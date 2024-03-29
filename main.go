package main

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// Globals:

var (
	gcpProjectID        string
	firebaseServiceFile string
	sendgridAPIKey      string
	sendgridTemplateID  string
)

const (
	appID = "canwork-api-chat-notification"
)

// Init function gets run automatically
func init() {
	if appengine.IsDevAppServer() {
		firebaseServiceFile = getEnv("CANWORK_FIREBASE_SERVICE_FILE", "")
		gcpProjectID = getEnv("GCP_PROJECT_ID", "")
	}
	sendgridAPIKey = getEnv("CANYA_SENDGRID_API_KEY", "")
	sendgridTemplateID = getEnv("CANYA_SENDGRID_TEMPLATE_ID", "d-bd1327c67b294710aa2b5dcdfe0da944")

	http.HandleFunc("/", handleRoot)
}

func main() {
	appengine.Main()
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	sendgrid.DefaultClient.HTTPClient = urlfetch.Client(ctx)

	client, err := getNewFirestoreClient(ctx)
	writeLogIfError(ctx, err)
	defer client.Close()

	psi := client.Collection("notifications").Where("chat", "==", true).Documents(ctx)
	log.Infof(ctx, "Fetched firebase notifications where chat == true")

	for {
		x, err := psi.Next()
		if err != nil {
			break
		}
		userID := x.Ref.ID
		log.Infof(ctx, "getting user from firestore: %s", userID)
		docsnap, err := client.Doc(fmt.Sprintf("users/%s", userID)).Get(ctx)
		if err != nil {
			log.Errorf(ctx, "failed to retrieve user %s: %s", userID, err.Error())
		} else {
			var user User
			if err := docsnap.DataTo(&user); err != nil {
				log.Errorf(ctx, "failed parsing user %s: %s", userID, err.Error())
			} else {
				log.Infof(ctx, "Got user to notify: %+s", user)
				sent := sendEmail(ctx, w, user.Name, user.Email)
				if sent {
					_, err := client.Doc(fmt.Sprintf("notifications/%s", userID)).Update(ctx, []firestore.Update{{Path: "chat", Value: false}})
					if err != nil {
						log.Errorf(ctx, "Error setting flag on notifications: %s", userID)
					}
				}
			}
		}
	}
	log.Infof(ctx, "Finished sending emails")
}

func sendEmail(ctx context.Context, w http.ResponseWriter, name string, email string) bool {
	m := mail.NewV3Mail()

	senderAddress := "support@canya.com"
	senderName := "CanYa support"
	e := mail.NewEmail(senderName, senderAddress)
	m.SetFrom(e)

	m.SetTemplateID(sendgridTemplateID)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(name, email),
	}
	p.AddTos(tos...)

	p.SetDynamicTemplateData("subject", "You have unread chat messages")
	p.SetDynamicTemplateData("title", "You have unread chat messages on CanWork")
	p.SetDynamicTemplateData("body", "People on CanWork are waiting for your response!")
	p.SetDynamicTemplateData("returnLinkText", "Visit CANWork")
	p.SetDynamicTemplateData("returnLinkUrl", "https://canwork.io/inbox/chat")

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(sendgridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		log.Errorf(ctx, "%s", err)
		return false
	} else {
		log.Infof(ctx, "API response status code: %d", response.StatusCode)
		log.Debugf(ctx, "API response body: %s", response.Body)
		log.Debugf(ctx, "API response headers: %+v", response.Headers)
		fmt.Fprintln(w, fmt.Sprintf("email sent with sendgrid response body: %s", response.Body))
		return true
	}
}

package main

import (
	"fmt"
	"net/http"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// Globals:

var (
	// firebaseServiceFile string
	sendgridAPIKey     string
	sendgridTemplateID string
)

const (
	gcpProjectID = "staging-can-work"
	appID        = "canwork-api-chat-notification"
)

// Init function gets run automatically
func init() {

	// firebaseServiceFile = getEnv("CANWORK_FIREBASE_SERVICE_FILE", "")
	// if firebaseServiceFile == "" {
	// 	panic(fmt.Sprintf("unable to find required environment variable: CANWORK_FIREBASE_SERVICE_FILE"))
	// }

	sendgridAPIKey = getEnv("CANYA_SENDGRID_API_KEY", "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")
	if sendgridAPIKey == "" {
		panic(fmt.Sprintf("unable to find required environment variable: CANYA_SENDGRID_API_KEY"))
	}

	sendgridTemplateID = getEnv("CANYA_SENDGRID_TEMPLATE_ID", "d-bd1327c67b294710aa2b5dcdfe0da944")

	http.HandleFunc("/", handleRoot)

}

func main() {
	appengine.Main()
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	sendgrid.DefaultClient.HTTPClient = urlfetch.Client(ctx)

	/*
	 * Send a test email with the templateID
	 * Taken from: https://github.com/sendgrid/sendgrid-go/blob/master/USE_CASES.md#transactional-templates
	 */

	m := mail.NewV3Mail()

	address := "alex@canya.com"
	name := "Big Boi"
	e := mail.NewEmail(name, address)
	m.SetFrom(e)

	m.SetTemplateID(sendgridTemplateID)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail("Cam", "cam@canya.com"),
		mail.NewEmail("Alex", "alex@canya.com"),
	}
	p.AddTos(tos...)

	p.SetDynamicTemplateData("subject", "Subject")
	p.SetDynamicTemplateData("title", "Test from GAE Standard Golang")
	p.SetDynamicTemplateData("body", "sdfsdfsd")
	p.SetDynamicTemplateData("returnLinkText", "Visit CANWork")
	p.SetDynamicTemplateData("returnLinkUrl", "https://canwork.io")

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(sendgridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		log.Errorf(ctx, "failed to hit sendgrid API: %s", err.Error())
		fmt.Fprintln(w, fmt.Sprintf("email not sent: %s", err.Error()))
	} else {
		log.Debugf(ctx, "API response status code: %d", response.StatusCode)
		log.Debugf(ctx, "API response body: %s", response.Body)
		log.Debugf(ctx, "API response headers: %+v", response.Headers)
		fmt.Fprintln(w, fmt.Sprintf("email sent with sendgrid response body: %s", response.Body))
	}

	// var err error

	// logger.("logggggg")
	// Creates a client.

	// client, err := getNewFirestoreClient(ctx)
	// writeLogIfError(ctx, err)
	// defer client.Close()

	// myID := "GW0A2f0pTOc559hfCT0sQqa1kgE3"

	// psi := client.Collection(fmt.Sprintf("who/%s/user", myID)).Documents(ctx)

	// for {
	// 	x, err := psi.Next()
	// 	if err != nil {
	// 		break
	// 	}
	// 	fmt.Fprintln(w, fmt.Sprintf("%+v", x.Data()))
	// }

}

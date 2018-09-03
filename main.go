package main

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

// Globals:

var (
	ctx                 context.Context
	firebaseServiceFile string
	sendgridAPIKey      string
)

const (
	firestoreProjectID = "staging-can-work"
)

// Init function gets run automatically
func init() {
	firebaseServiceFile = getEnv("CANWORK_FIREBASE_SERVICE_FILE", "")
	if firebaseServiceFile == "" {
		panic(fmt.Sprintf("unable to find required environment variable: CANWORK_FIREBASE_SERVICE_FILE"))
	}
	sendgridAPIKey = getEnv("CANYA_SENDGRID_API_KEY", "")
	if sendgridAPIKey == "" {
		panic(fmt.Sprintf("unable to find required environment variable: CANYA_SENDGRID_API_KEY"))
	}
}

func main() {
	http.HandleFunc("/", handleRoot)
	appengine.Main()
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := appengine.NewContext(r)

	client, err := getNewFirestoreClient(ctx)
	writeLogIfError(ctx, err)
	defer client.Close()

	myID := "GW0A2f0pTOc559hfCT0sQqa1kgE3"

	psi := client.Collection(fmt.Sprintf("who/%s/user", myID)).Documents(ctx)

	for {
		x, err := psi.Next()
		if err != nil {
			break
		}
		fmt.Fprintln(w, fmt.Sprintf("%+v", x.Data()))
	}

}

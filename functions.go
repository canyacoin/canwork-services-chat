package main

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/appengine/log"
)

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	"cloud.google.com/go/firestore"
// 	sendgrid "github.com/sendgrid/sendgrid-go"
// 	"github.com/sendgrid/sendgrid-go/helpers/mail"
// 	"google.golang.org/api/option"
// 	"google.golang.org/appengine/log"
// )

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getNewFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	return firestore.NewClient(ctx, gcpProjectID, option.WithServiceAccountFile(firebaseServiceFile))
}

func writeLogIfError(ctx context.Context, err error) {
	if err != nil {
		log.Errorf(ctx, "Err: %s", err.Error())
	}
}

// func sendEmail() {
// 	from := mail.NewEmail("Example User", "valthrex@gmail.com")
// 	fmt.Println(from)
// 	subject := "Sending with SendGrid is Fun"
// 	to := mail.NewEmail("Example User", "valthrex@gmail.com")
// 	plainTextContent := "and easy to do anywhere, even with Go"
// 	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
// 	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
// 	client := sendgrid.NewSendClient(sendgridAPIKey)
// 	response, err := client.Send(message)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println(response.StatusCode)
// 		fmt.Println(response.Body)
// 		fmt.Println(response.Headers)
// 	}
// }

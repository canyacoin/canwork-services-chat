package main

import (
	"context"
	"fmt"
	"os"

	"firebase.google.com/go"

	"cloud.google.com/go/firestore"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func getEnv(key, fallback string) string {
	returnVal := fallback
	if value, ok := os.LookupEnv(key); ok {
		returnVal = value
	}
	if returnVal == "" {
		panic(fmt.Sprintf("Unable to retrieve key: %s", key))
	}
	return returnVal
}

func getNewFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	if !appengine.IsDevAppServer() {
		credentials, error := google.FindDefaultCredentials(ctx, compute.ComputeScope)
		if error != nil {
			panic(fmt.Sprintf("Can't get default firebase credentials: %v", error))
		}
		conf := &firebase.Config{ProjectID: credentials.ProjectID}
		fmt.Printf(credentials.ProjectID)
		app, err := firebase.NewApp(ctx, conf)
		if err != nil {
			panic(fmt.Sprintf("Failed to load firebase: %v", error))
		}
		return app.Firestore(ctx)
	}

	return firestore.NewClient(ctx, gcpProjectID, option.WithServiceAccountFile(firebaseServiceFile))
}

func writeLogIfError(ctx context.Context, err error) {
	if err != nil {
		log.Errorf(ctx, "Err: %s", err.Error())
	}
}

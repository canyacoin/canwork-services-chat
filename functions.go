package main

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/appengine/log"
)

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

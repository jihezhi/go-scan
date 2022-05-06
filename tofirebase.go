package main

import (
	"context"
	"log"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func StoreToFirestore(logTransfer LogTransfer) {
	ctx := context.Background()
	config := &firebase.Config{ProjectID: "fantom-test-f3bf5"}
	opt := option.WithCredentialsFile("fantom-test-f3bf5-39bbe2b9a415.json")
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatal(err)
	}
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer firestoreClient.Close()

	document, err := firestoreClient.Collection("fantom").Doc(logTransfer.From.String()).Get(ctx)

	if len(document.Data()) > 0 {
		index := len(document.Data()) + 1
		value := map[string]interface{}{
			"From":         logTransfer.From.String(),
			"To":           logTransfer.To.String(),
			"TokenAddress": logTransfer.TokenAddress.String(),
			"BLockNumber":  logTransfer.BlockNumber,
			"Tokens":       logTransfer.Tokens.Int64(),
			"TimeStamp":    logTransfer.Time,
		}
		firestoreClient.Collection("fantom").Doc(logTransfer.From.String()).Update(ctx, []firestore.Update{
			{
				Path:  strconv.Itoa(index),
				Value: value,
			},
		})
	} else {
		firestoreClient.Collection("fantom").Doc(logTransfer.From.String()).Set(ctx, map[string]interface{}{
			"1": map[string]interface{}{
				"From":         logTransfer.From.String(),
				"To":           logTransfer.To.String(),
				"TokenAddress": logTransfer.TokenAddress.String(),
				"BLockNumber":  logTransfer.BlockNumber,
				"Tokens":       logTransfer.Tokens.Int64(),
				"TimeStamp":    logTransfer.Time,
			}})
	}

}

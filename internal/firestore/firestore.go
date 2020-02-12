package firestore

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Save is function for firestore.
func Save(collection string, id string, data map[string]string) {
	// Use the application default credentials
	ctx := context.Background()
	sa := option.WithCredentialsFile("configs/key/erc-checker-firebase-adminsdk-92q34-5687c47029.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	// add data to firestore
	_, err = client.Collection(collection).Doc(id).Set(ctx, data)
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
}

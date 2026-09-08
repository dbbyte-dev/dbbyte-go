package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dbbyte-dev/dbbyte-go"
)

func main() {
	dbb, err := dbbyte.New(dbbyte.ClientOptions{
		APIKey: os.Getenv("DBBYTE_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	uploaded, err := dbb.Objects.Upload(ctx, "my-bucket", dbbyte.UploadOptions{
		Key:         "images/logo.png",
		Body:        []byte("...file bytes..."),
		ContentType: "image/png",
	})
	if err != nil {
		var dbErr *dbbyte.Error
		if errors.As(err, &dbErr) {
			fmt.Printf("upload failed: %s (status %d)\n", dbErr.Message, dbErr.Status)
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("uploaded: %s (%d bytes)\n", uploaded.Key, uploaded.Size)

	signed, err := dbb.Objects.Sign(ctx, "my-bucket", "images/logo.png", dbbyte.SignOptions{
		ExpiresIn: 3600,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("signed URL:", signed.URL)
}

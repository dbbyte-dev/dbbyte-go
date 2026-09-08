// Package dbbyte is the official Go SDK for the DBByte storage/CDN API.
//
// Basic usage:
//
//	dbb, err := dbbyte.New(dbbyte.ClientOptions{APIKey: "dbk_..."})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	uploaded, err := dbb.Objects.Upload(ctx, "my-bucket", dbbyte.UploadOptions{
//	    Key:  "images/logo.png",
//	    Body: fileBytes,
//	})
package dbbyte

import "context"

// DBByte is the top-level client. Construct with New.
type DBByte struct {
	client  *client
	Objects *ObjectsResource
}

// New constructs a DBByte client. APIKey is required; all other options
// have sane defaults (see ClientOptions).
func New(opts ClientOptions) (*DBByte, error) {
	c, err := newClient(opts)
	if err != nil {
		return nil, err
	}
	return &DBByte{
		client:  c,
		Objects: &ObjectsResource{client: c},
	}, nil
}

// Health calls GET /health.
func (d *DBByte) Health(ctx context.Context) (*HealthResponse, error) {
	var out HealthResponse
	if err := d.client.requestJSON(ctx, "GET", "/health", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

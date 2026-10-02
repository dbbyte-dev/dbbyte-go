package dbbyte

import "context"

// CDNResource groups cache invalidation operations.
type CDNResource struct {
	client *client
}

// Purge invalidates a public file or all public files in a bucket.
func (c *CDNResource) Purge(ctx context.Context, opts PurgeOptions) (*PurgeResponse, error) {
	var out PurgeResponse
	if err := c.client.requestJSON(ctx, "POST", "/purge", nil, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

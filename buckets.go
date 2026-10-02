package dbbyte

import (
	"context"
	"net/http"
)

// BucketsResource groups account-level bucket management operations.
type BucketsResource struct {
	client *client
}

// List returns all buckets owned by the account.
func (b *BucketsResource) List(ctx context.Context) (*BucketListResponse, error) {
	var out BucketListResponse
	if err := b.client.requestJSON(ctx, http.MethodGet, "/buckets", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get returns one bucket owned by the account.
func (b *BucketsResource) Get(ctx context.Context, slug string) (*BucketResponse, error) {
	var out BucketResponse
	path := "/buckets/" + encodeSegment(slug)
	if err := b.client.requestJSON(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a bucket with the supplied settings.
func (b *BucketsResource) Create(ctx context.Context, opts CreateBucketOptions) (*BucketResponse, error) {
	var out BucketResponse
	if err := b.client.requestJSON(ctx, http.MethodPost, "/buckets", nil, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update changes supported boolean settings on a bucket.
func (b *BucketsResource) Update(ctx context.Context, slug string, opts UpdateBucketOptions) (*BucketResponse, error) {
	var out BucketResponse
	path := "/buckets/" + encodeSegment(slug)
	if err := b.client.requestJSON(ctx, http.MethodPatch, path, nil, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a non-default bucket and its stored contents.
func (b *BucketsResource) Delete(ctx context.Context, slug string) (*BucketDeleteResponse, error) {
	var out BucketDeleteResponse
	path := "/buckets/" + encodeSegment(slug)
	if err := b.client.requestJSON(ctx, http.MethodDelete, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

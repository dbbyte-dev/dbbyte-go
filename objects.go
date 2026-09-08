package dbbyte

import (
	"context"
	"strconv"
)

// ObjectsResource groups every object-level operation (upload, list,
// download, delete, sign, versions) under one bucket-scoped namespace,
// matching the Node and Python SDKs.
type ObjectsResource struct {
	client *client
}

// Upload uploads a single object to bucket.
func (o *ObjectsResource) Upload(ctx context.Context, bucket string, opts UploadOptions) (*UploadedObject, error) {
	query := queryParams{}
	if opts.Overwrite {
		query["overwrite"] = "true"
	}

	var out UploadedObject
	err := o.client.requestForm(ctx, "POST", "/storage/"+bucket+"/upload", query, opts.Key, opts.Body, opts.ContentType, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// List returns a page of objects in bucket.
func (o *ObjectsResource) List(ctx context.Context, bucket string, opts ListObjectsOptions) (*ListObjectsResponse, error) {
	query := queryParams{
		"prefix": opts.Prefix,
		"cursor": opts.Cursor,
	}
	if opts.Limit > 0 {
		query["limit"] = strconv.Itoa(opts.Limit)
	}

	var out ListObjectsResponse
	if err := o.client.requestJSON(ctx, "GET", "/storage/"+bucket+"/list", query, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Download fetches a single object's bytes. Pass version to fetch a
// specific version instead of the latest.
func (o *ObjectsResource) Download(ctx context.Context, bucket, key string, version string) (*DownloadResult, error) {
	query := queryParams{"version": version}
	return o.client.requestBinary(ctx, "GET", "/storage/"+bucket+"/download/"+encodeKey(key), query)
}

// DownloadAllVersions returns a zip of every version of key.
func (o *ObjectsResource) DownloadAllVersions(ctx context.Context, bucket, key string) (*DownloadResult, error) {
	query := queryParams{"version": "all"}
	return o.client.requestBinary(ctx, "GET", "/storage/"+bucket+"/download/"+encodeKey(key), query)
}

// Delete deletes an object. Pass version to delete a specific version
// instead of the latest.
func (o *ObjectsResource) Delete(ctx context.Context, bucket, key string, version string) (*DeleteResponse, error) {
	query := queryParams{"version": version}
	var out DeleteResponse
	if err := o.client.requestJSON(ctx, "DELETE", "/storage/"+bucket+"/delete/"+encodeKey(key), query, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAllVersions permanently deletes every version of key.
func (o *ObjectsResource) DeleteAllVersions(ctx context.Context, bucket, key string) (*DeleteResponse, error) {
	query := queryParams{"version": "all"}
	var out DeleteResponse
	if err := o.client.requestJSON(ctx, "DELETE", "/storage/"+bucket+"/delete/"+encodeKey(key), query, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Sign creates a signed, time-limited URL for key.
func (o *ObjectsResource) Sign(ctx context.Context, bucket, key string, opts SignOptions) (*SignedUrlResponse, error) {
	body := map[string]interface{}{}
	if opts.Version != "" {
		body["version"] = opts.Version
	}
	if opts.ExpiresIn > 0 {
		body["expiresIn"] = opts.ExpiresIn
	}

	var out SignedUrlResponse
	if err := o.client.requestJSON(ctx, "POST", "/storage/"+bucket+"/sign/"+encodeKey(key), nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Versions lists the version history of key.
func (o *ObjectsResource) Versions(ctx context.Context, bucket, key string, opts ListVersionsOptions) (*ListVersionsResponse, error) {
	query := queryParams{"cursor": opts.Cursor}
	if opts.Limit > 0 {
		query["limit"] = strconv.Itoa(opts.Limit)
	}

	var out ListVersionsResponse
	if err := o.client.requestJSON(ctx, "GET", "/storage/"+bucket+"/versions/"+encodeKey(key), query, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

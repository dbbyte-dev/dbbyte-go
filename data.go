package dbbyte

import (
	"context"
	"strconv"
	"strings"
)

// DataResource groups operations for published file and Google Docs APIs.
type DataResource struct {
	client *client
}

// Get reads a page from a published Data API. Filters are exact string matches.
func (d *DataResource) Get(ctx context.Context, slug string, opts DataQueryOptions) (*DataPage, error) {
	query := queryParams{"cursor": opts.Cursor}
	if len(opts.Fields) > 0 {
		query["fields"] = strings.Join(opts.Fields, ",")
	}
	if opts.Limit > 0 {
		query["limit"] = strconv.Itoa(opts.Limit)
	}
	for field, value := range opts.Filters {
		query[field] = value
	}

	var out DataPage
	path := "/v1/data/" + encodeSegment(slug)
	if err := d.client.requestJSON(ctx, "GET", path, query, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create inserts a JSON record into a write-enabled Data API.
func (d *DataResource) Create(ctx context.Context, slug string, record DataRecord) (*DataRecordResponse, error) {
	var out DataRecordResponse
	path := "/v1/data/" + encodeSegment(slug)
	if err := d.client.requestJSON(ctx, "POST", path, nil, record, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AppendText appends to a write-enabled Google Docs Data API.
func (d *DataResource) AppendText(ctx context.Context, slug, text string) (*AppendTextResponse, error) {
	var out AppendTextResponse
	path := "/v1/data/" + encodeSegment(slug)
	body := map[string]string{"appendText": text}
	if err := d.client.requestJSON(ctx, "POST", path, nil, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches writable fields on a JSON record.
func (d *DataResource) Update(ctx context.Context, slug, recordID string, patch DataRecord) (*DataRecordResponse, error) {
	var out DataRecordResponse
	path := "/v1/data/" + encodeSegment(slug) + "/" + encodeSegment(recordID)
	if err := d.client.requestJSON(ctx, "PATCH", path, nil, patch, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a record from a write-enabled JSON Data API.
func (d *DataResource) Delete(ctx context.Context, slug, recordID string) (*DataDeleteResponse, error) {
	var out DataDeleteResponse
	path := "/v1/data/" + encodeSegment(slug) + "/" + encodeSegment(recordID)
	if err := d.client.requestJSON(ctx, "DELETE", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

package dbbyte

// DataRecord is a record returned by a published Data API.
type DataRecord map[string]interface{}

// Bucket is the public representation returned by bucket endpoints.
type Bucket struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Slug              string `json:"slug"`
	Region            string `json:"region"`
	ReplicaCount      int    `json:"replicaCount"`
	IsPublic          bool   `json:"isPublic"`
	CDN               bool   `json:"cdn"`
	VersioningEnabled bool   `json:"versioningEnabled"`
	FileCount         int64  `json:"fileCount"`
	Size              int64  `json:"size"`
	IsDefault         bool   `json:"isDefault"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

// BucketListResponse contains the account's buckets.
type BucketListResponse struct {
	Buckets []Bucket `json:"buckets"`
}

// BucketResponse contains a bucket and optional create success flag.
type BucketResponse struct {
	Success bool   `json:"success,omitempty"`
	Bucket  Bucket `json:"bucket"`
}

// CreateBucketOptions configures Buckets.Create.
type CreateBucketOptions struct {
	Name         string `json:"name"`
	Region       string `json:"region"`
	ReplicaCount int    `json:"replicaCount,omitempty"`
	IsPublic     bool   `json:"isPublic,omitempty"`
	CDN          bool   `json:"cdn,omitempty"`
}

// UpdateBucketOptions configures Buckets.Update. Pointer fields allow explicit false values.
type UpdateBucketOptions struct {
	VersioningEnabled *bool `json:"versioningEnabled,omitempty"`
	IsPublic          *bool `json:"isPublic,omitempty"`
	CDN               *bool `json:"cdn,omitempty"`
}

// BucketDeleteResponse is returned when a bucket is removed.
type BucketDeleteResponse struct {
	Success bool `json:"success"`
}

// DataQueryOptions configures Data.Get. Filters use exact string matches.
type DataQueryOptions struct {
	Fields  []string
	Cursor  string
	Limit   int
	Filters map[string]string
}

// DataPage is one page from a published Data API.
type DataPage struct {
	Data       []DataRecord `json:"data"`
	NextCursor *string      `json:"next_cursor"`
}

// DataRecordResponse wraps a created or updated record.
type DataRecordResponse struct {
	Data DataRecord `json:"data"`
}

// AppendTextResponse is returned by Data.AppendText.
type AppendTextResponse struct {
	Success            bool `json:"success"`
	AppendedCharacters int  `json:"appendedCharacters"`
}

// DataDeleteResponse is returned by Data.Delete.
type DataDeleteResponse struct {
	Success bool `json:"success"`
}

// PurgeScope selects a file or whole-bucket CDN purge.
type PurgeScope string

const (
	PurgeFile   PurgeScope = "file"
	PurgeBucket PurgeScope = "bucket"
)

// PurgeOptions configures CDN.Purge. Path is required for file purges.
type PurgeOptions struct {
	Slug  string     `json:"slug"`
	Scope PurgeScope `json:"scope"`
	Path  string     `json:"path,omitempty"`
}

// PurgeResponse reports CDN purge results.
type PurgeResponse struct {
	Success bool     `json:"success"`
	Purged  int      `json:"purged"`
	Failed  []string `json:"failed"`
}

// HealthResponse mirrors GET /health.
type HealthResponse struct {
	Status         string       `json:"status"` // "ok" | "degraded"
	Timestamp      string       `json:"timestamp"`
	UptimeSeconds  int          `json:"uptimeSeconds"`
	Checks         HealthChecks `json:"checks"`
	ResponseTimeMs int          `json:"responseTimeMs"`
}

type HealthChecks struct {
	API      HealthCheckStatus    `json:"api"`
	Database HealthDatabaseStatus `json:"database"`
}

type HealthCheckStatus struct {
	Status string `json:"status"`
}

type HealthDatabaseStatus struct {
	Status    string `json:"status"`
	LatencyMs *int   `json:"latencyMs"`
}

// UploadedObject is returned on a successful upload. Matches the storage
// server's real response shape exactly - not inferred.
type UploadedObject struct {
	Success  bool   `json:"success"`
	Key      string `json:"key"`  // server-generated storage key, not the name you passed
	Name     string `json:"name"` // the name you uploaded it as
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
	FileID   string `json:"fileId"`
}

// UploadOptions configures Objects.Upload.
type UploadOptions struct {
	Key         string
	Body        []byte
	ContentType string
	Overwrite   bool
}

type FileVersionPreview struct {
	Key       string `json:"key"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mimeType"`
	IsLatest  bool   `json:"isLatest"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type FileListItem struct {
	Key             string               `json:"key"`
	Name            string               `json:"name"`
	Size            int64                `json:"size"`
	IsPublic        bool                 `json:"isPublic"`
	CreatedAt       string               `json:"createdAt"`
	UpdatedAt       string               `json:"updatedAt"`
	VersionCount    int                  `json:"versionCount"`
	Versions        []FileVersionPreview `json:"versions"`
	HasMoreVersions bool                 `json:"hasMoreVersions"`
}

// ListObjectsOptions configures Objects.List.
//
// Prefix filters by the uploaded, human-readable file name.
type ListObjectsOptions struct {
	Prefix string
	Cursor string
	Limit  int // 0 means "not set"
}

type ListObjectsResponse struct {
	Files      []FileListItem `json:"files"`
	NextCursor *string        `json:"nextCursor"`
	HasMore    bool           `json:"hasMore"`
}

// SignOptions configures Objects.Sign.
type SignOptions struct {
	Version   string
	ExpiresIn int // seconds; 0 means "not set"
}

type SignedUrlResponse struct {
	URL       string `json:"url"`
	ExpiresIn int    `json:"expiresIn"`
}

type VersionEntry struct {
	VersionKey string `json:"versionKey"`
	Size       int64  `json:"size"`
	MimeType   string `json:"mimeType"`
	IsLatest   bool   `json:"isLatest"`
	IsDeleted  bool   `json:"isDeleted"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// ListVersionsOptions configures Objects.Versions.
type ListVersionsOptions struct {
	Limit  int // 0 means "not set"
	Cursor string
}

type ListVersionsResponse struct {
	Versions   []VersionEntry `json:"versions"`
	NextCursor *string        `json:"nextCursor"`
	HasMore    bool           `json:"hasMore"`
}

type DeleteResponse struct {
	Success         bool `json:"success"`
	DeletedVersions *int `json:"deletedVersions,omitempty"` // only present when deleting every version at once
}

// DownloadResult holds a downloaded object's bytes and metadata. Filename
// is nil when the server didn't send a Content-Disposition filename.
type DownloadResult struct {
	Data        []byte
	ContentType string
	Filename    *string
}

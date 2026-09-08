package dbbyte

// HealthResponse mirrors GET /health.
type HealthResponse struct {
	Status        string            `json:"status"` // "ok" | "degraded"
	Timestamp     string            `json:"timestamp"`
	UptimeSeconds int               `json:"uptimeSeconds"`
	Checks        HealthChecks      `json:"checks"`
	ResponseTimeMs int              `json:"responseTimeMs"`
}

type HealthChecks struct {
	API      HealthCheckStatus     `json:"api"`
	Database HealthDatabaseStatus  `json:"database"`
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
// NOTE: as of the current API, Prefix filters on the storage-generated
// `key` field (a timestamp-prefixed, slash-flattened server ID), not the
// human-readable `name` you uploaded with. It will rarely match what you
// expect - confirmed via round-trip testing against a real server. Treat
// as unreliable until the backend filters on `name` instead.
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

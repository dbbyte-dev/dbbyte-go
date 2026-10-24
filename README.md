# dbbyte-go

Official Go SDK for the DBByte storage and Data APIs.

The client exposes `Buckets`, `Objects`, `Data`, and `CDN` resources. Bucket
management requires an account-level API key. See the [REST API reference](https://dbbyte.com/api-reference)
and [OpenAPI specification](https://dbbyte.com/openapi.yaml) for operations,
permissions, and response shapes.

## Install

```bash
go get github.com/dbbyte-dev/dbbyte-go
```

## Usage

```go
dbb, err := dbbyte.New(dbbyte.ClientOptions{
    APIKey: "dbk_...",
})
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()

uploaded, err := dbb.Objects.Upload(ctx, "my-bucket", dbbyte.UploadOptions{
    Key:  "images/logo.png",
    Body: fileBytes,
})
```

See `examples/basic` for a fuller walkthrough.

### Direct browser uploads

`Objects.UploadURL` mints a short-lived link so a browser can upload straight
to DBByte without your write API key ever reaching the frontend. Call it from
your own server:

```go
link, err := dbb.Objects.UploadURL(ctx, "my-bucket", dbbyte.UploadURLOptions{
    Name:        "song.mp3",
    ContentType: "audio/mpeg",
    Size:        5 * 1024 * 1024, // max bytes the link accepts (up to 50 MB)
    ExpiresIn:   900,             // seconds; the server caps this at 3600
})
if err != nil {
    log.Fatal(err)
}
// Hand link.URL, link.Method ("PUT") and link.Headers to the browser, which
// sends the file bytes as the request body.
```

## Errors

Every failed request returns a `*dbbyte.Error` with `Message`, `Status`,
`Hint`, and `RetryAfter`. The API has no `code` field — don't rely on one.

A `429` can mean two different things, and only one is retryable:

- **Burst rate limit** — `RetryAfter` is set (seconds). Transient, safe to
  retry after that delay.
- **Plan API request quota exceeded** — `RetryAfter` is `0`. Retrying
  immediately won't help until the billing period resets or the plan is
  upgraded.

Check `RetryAfter`, not the message text, to tell them apart.

```go
var dbErr *dbbyte.Error
if errors.As(err, &dbErr) {
    if dbErr.RetryAfter > 0 {
        time.Sleep(time.Duration(dbErr.RetryAfter) * time.Second)
        // retry
    }
}
```

## Base URL / gateway note

By default the client talks to `https://api.dbbyte.com`. In production the
Cloudflare Worker gateway in front of that host inserts `/api/v1` itself,
so paths used internally by this SDK do **not** include `/v1`. If you point
`BaseURL` at a Next.js dev server directly (bypassing the Worker), include
`/api/v1` in `BaseURL` instead.

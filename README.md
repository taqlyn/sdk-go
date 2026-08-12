# Taqlyn Go SDK

Taqlyn **server** SDK for creating and managing short links. Requests are
authenticated with an Ed25519 private key and the `X-Taqlyn-*` headers.

This package is server-only. Do not embed its private key in mobile, browser, or
desktop apps. It intentionally does not implement Match, resolve, or deferred
mobile handoff flows.

## Install

```bash
go get github.com/taqlyn/sdk-go
```

## Quickstart

Use the same environment variables as the Node server SDK:

```bash
export TAQLYN_BASE_URL=http://localhost:8080
export TAQLYN_CLIENT_ID=app_test_abc
export TAQLYN_PRIVATE_KEY='-----BEGIN PRIVATE KEY-----
...
-----END PRIVATE KEY-----'
```

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	taqlyn "github.com/taqlyn/sdk-go"
)

func main() {
	client, err := taqlyn.NewClient(taqlyn.Config{
		BaseURL:       os.Getenv("TAQLYN_BASE_URL"),
		ClientID:      os.Getenv("TAQLYN_CLIENT_ID"),
		PrivateKeyPEM: os.Getenv("TAQLYN_PRIVATE_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	link, err := client.CreateShortLink(context.Background(), taqlyn.CreateShortLinkRequest{
		DestinationWeb: "https://example.com/offer",
		Mode:           taqlyn.ModeWebOnly,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(link.ShortURL)
}
```

`TAQLYN_PRIVATE_KEY` must be a PKCS#8 Ed25519 PEM, not an `sk_*` handle.
Literal `\n` sequences in the environment value are accepted. Never commit the
private key.

## API

- `CreateShortLink`
- `GetShortLink`
- `PatchShortLink`
- `DeleteShortLink`

All methods accept a `context.Context`. A non-success HTTP response is returned
as `*taqlyn.APIError`.

## Signing

The SDK signs this canonical message, with no trailing newline:

```text
taqlyn-v1
{METHOD}
{PATH}
{unixTimestamp}
{clientId}
{hex(sha256(body))}
```

The signature is standard base64-encoded Ed25519 and is sent with
`X-Taqlyn-Client-Id`, `X-Taqlyn-Timestamp`, and `X-Taqlyn-Signature`.

## Test

```bash
go test ./...
```

## License

MIT — see [LICENSE](./LICENSE).

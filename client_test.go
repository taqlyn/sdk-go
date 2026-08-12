package taqlyn

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateShortLink(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", request.Method)
		}
		if request.URL.Path != "/v1/short-links" {
			t.Errorf("path = %s, want /v1/short-links", request.URL.Path)
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if got := string(body); got != goldenBody {
			t.Errorf("body = %s, want %s", got, goldenBody)
		}
		if got := request.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("content type = %q", got)
		}
		if got := request.Header.Get("X-Taqlyn-Client-Id"); got != "app_test_abc" {
			t.Errorf("client ID header = %q", got)
		}
		if got := request.Header.Get("X-Taqlyn-Timestamp"); got != "1700000000" {
			t.Errorf("timestamp header = %q", got)
		}
		if got := request.Header.Get("X-Taqlyn-Signature"); got != goldenSignature {
			t.Errorf("signature = %q, want golden signature", got)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{
			"id":"sl_test",
			"code":"Ab12Cd",
			"shortUrl":"https://go.example/Ab12Cd",
			"host":"go.example",
			"mode":"web_only",
			"destinationWeb":"https://example.com/offer",
			"env":"sandbox",
			"orgId":"org_test",
			"appId":"app_test"
		}`))
	}))
	defer server.Close()

	client, err := NewClient(Config{
		BaseURL:       server.URL + "/",
		ClientID:      "app_test_abc",
		PrivateKeyPEM: goldenPrivateKeyPEM(),
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.now = func() time.Time { return time.Unix(1700000000, 0) }

	link, err := client.CreateShortLink(context.Background(), CreateShortLinkRequest{
		DestinationWeb: "https://example.com/offer",
		Mode:           ModeWebOnly,
	})
	if err != nil {
		t.Fatalf("CreateShortLink: %v", err)
	}
	if link.ID != "sl_test" || link.ShortURL != "https://go.example/Ab12Cd" {
		t.Errorf("unexpected link: %#v", link)
	}
}

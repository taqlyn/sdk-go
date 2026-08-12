package taqlyn

import (
	"encoding/base64"
	"encoding/pem"
	"net/http"
	"testing"
)

const (
	goldenKeyDERBase64 = "MC4CAQAwBQYDK2VwBCIEIAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8g"
	goldenBody         = `{"destinationWeb":"https://example.com/offer","mode":"web_only"}`
	goldenSignature    = "zTe0VimeWAe6dzpPxAIn+DDR46E58G63ypiSTkXd1jT1o3oxYJ4jzAof05lf3s/8sbZ7l46VjDh8ohtB+NISAA=="
)

func goldenPrivateKeyPEM() string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: mustDecodeBase64(goldenKeyDERBase64),
	}))
}

func mustDecodeBase64(value string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		panic(err)
	}
	return decoded
}

func TestSignedHeadersGoldenVector(t *testing.T) {
	privateKey, err := parsePrivateKey(goldenPrivateKeyPEM())
	if err != nil {
		t.Fatalf("parse private key: %v", err)
	}

	headers := signedHeaders(
		privateKey,
		http.MethodPost,
		"/v1/short-links",
		"app_test_abc",
		1700000000,
		[]byte(goldenBody),
	)

	if got := headers.Get("X-Taqlyn-Client-Id"); got != "app_test_abc" {
		t.Errorf("client ID header = %q", got)
	}
	if got := headers.Get("X-Taqlyn-Timestamp"); got != "1700000000" {
		t.Errorf("timestamp header = %q", got)
	}
	if got := headers.Get("X-Taqlyn-Signature"); got != goldenSignature {
		t.Errorf("signature = %q, want %q", got, goldenSignature)
	}
}

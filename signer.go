package taqlyn

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const signingVersion = "taqlyn-v1"

func parsePrivateKey(privateKeyPEM string) (ed25519.PrivateKey, error) {
	normalized := strings.ReplaceAll(privateKeyPEM, `\n`, "\n")
	block, _ := pem.Decode([]byte(normalized))
	if block == nil {
		return nil, fmt.Errorf("taqlyn: private key is not valid PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("taqlyn: parse PKCS#8 private key: %w", err)
	}

	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("taqlyn: private key is not Ed25519")
	}
	return privateKey, nil
}

func canonicalMessage(method, path, timestamp, clientID string, body []byte) string {
	bodyHash := sha256.Sum256(body)
	return strings.Join([]string{
		signingVersion,
		strings.ToUpper(method),
		path,
		timestamp,
		clientID,
		hex.EncodeToString(bodyHash[:]),
	}, "\n")
}

func signedHeaders(
	privateKey ed25519.PrivateKey,
	method string,
	path string,
	clientID string,
	timestamp int64,
	body []byte,
) http.Header {
	timestampString := strconv.FormatInt(timestamp, 10)
	message := canonicalMessage(method, path, timestampString, clientID, body)
	signature := ed25519.Sign(privateKey, []byte(message))

	headers := make(http.Header)
	headers.Set("X-Taqlyn-Client-Id", clientID)
	headers.Set("X-Taqlyn-Timestamp", timestampString)
	headers.Set("X-Taqlyn-Signature", base64.StdEncoding.EncodeToString(signature))
	return headers
}

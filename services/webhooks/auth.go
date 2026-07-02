package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// generateSignature function generates the signature for webhook.
// it returns the value in the format t=2435....23.s=eajn...jda.
// t represents the timestamp of generation and s is the signature itself.
func (s *Service) generateSignature(text []byte, job Job) (string, error) {
	// get the timestamp and generate payload
	timestamp := time.Now().Unix()
	payload := fmt.Sprintf("%d.%s", timestamp, string(text))

	// get the secret
	secret, err := s.webhookRepository.GetWebhookSecretByID(job.webhook.ID)
	if err != nil {
		return "", err
	}

	// generate the sha256 hash
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(payload))
	signature := hex.EncodeToString(hash.Sum(nil))

	return fmt.Sprintf("t=%d.s=%s", timestamp, signature), nil
}

// An example function to verify the webhook signature in golang.
func VerifySignature(payload string, header string, secret string) (bool, error) {
	// extract the data
	var timestampStr, signatureStr string

	parts := strings.Split(header, ".")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "t":
			timestampStr = kv[1]
		case "s":
			signatureStr = kv[1]
		}
	}

	if timestampStr == "" || signatureStr == "" {
		return false, errors.New("invalid or missing signature")
	}

	// prevent replay attacks by checking the timestamp
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false, errors.New("invalid timestamp")
	}

	currentTime := time.Now().Unix()

	if math.Abs(float64(currentTime-timestamp)) > 300 {
		return false, errors.New("maybe a replay attack")
	}

	// recreate the signature
	stringToSign := fmt.Sprintf("%s.%s", timestampStr, payload)

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(stringToSign))
	genSignature := hex.EncodeToString(hash.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(signatureStr), []byte(genSignature)) != 1 {
		return false, errors.New("webhook signature mismatch")
	}

	return true, nil
}

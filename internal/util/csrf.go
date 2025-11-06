package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GenerateCSRFToken(secret string) string {
	tokenID := uuid.New().String()

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(tokenID))
	signature := h.Sum(nil)

	token := fmt.Sprintf("%s.%s", tokenID, base64.URLEncoding.EncodeToString(signature))
	return token
}

func ValidateCSRFToken(token, secret string) bool {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return false
	}

	tokenID := parts[0]
	signatureStr := parts[1]

	signature, err := base64.URLEncoding.DecodeString(signatureStr)
	if err != nil {
		return false
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(tokenID))
	expectedSignature := h.Sum(nil)

	return hmac.Equal(signature, expectedSignature)
}

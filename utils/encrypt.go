package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func GenerateSignature(orderID, statusCode, grossAmount, serverKey string) string {
	signatureStr := orderID + statusCode + grossAmount + serverKey
	hash := sha512.Sum512([]byte(signatureStr))
	return hex.EncodeToString(hash[:])
}

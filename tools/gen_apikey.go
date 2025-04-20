package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func main() {
	serviceName := "user-service"
	signatureKay := "IphDM6yqXe0n0o2CyMpV"
	requestAt := fmt.Sprintf("%d", time.Now().Unix())

	raw := fmt.Sprintf("%s:%s:%s", serviceName, signatureKay, requestAt)

	hash := sha256.New()
	hash.Write([]byte(raw))
	apiKey := hex.EncodeToString(hash.Sum(nil))

	fmt.Println("x-serice-name :", serviceName)
	fmt.Println("x-request-at :", requestAt)
	fmt.Println("x-api-key :", apiKey)
}

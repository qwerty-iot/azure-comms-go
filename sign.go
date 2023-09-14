package azurecomms

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func parseConnectionString(cs string, require ...string) (map[string]string, error) {
	m := map[string]string{}
	for _, s := range strings.Split(cs, ";") {
		if s == "" {
			continue
		}
		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			return nil, errors.New("malformed connection string")
		}
		m[kv[0]] = kv[1]
	}
	for _, k := range require {
		if s := m[k]; s == "" {
			return nil, fmt.Errorf("%s is required", k)
		}
	}
	return m, nil
}

func signRequest(key string, req *http.Request, body []byte) error {

	now := time.Now().UTC()
	nowString := now.Format(time.RFC1123)
	if strings.HasSuffix(nowString, "UTC") {
		nowString = nowString[:len(nowString)-3] + "GMT"
	}

	hash := sha256.Sum256(body)
	hashB64 := base64.StdEncoding.EncodeToString(hash[:])

	stringToSign := req.Method + "\n" +
		req.URL.Path + "?" + req.URL.Query().Encode() + "\n" +
		nowString + ";" + req.Host + ";" + hashB64

	rawKey, _ := base64.StdEncoding.DecodeString(key)

	hm := hmac.New(sha256.New, rawKey)
	hm.Write([]byte(stringToSign))
	//hm.Write(rawKey)
	signature := hm.Sum(nil)

	//req.Header.Set("x-ms-date", nowString)
	//req.Header.Set("x-ms-content-sha256", hashB64)
	req.Header["x-ms-date"] = []string{nowString}
	req.Header["x-ms-content-sha256"] = []string{hashB64}
	//req.Header["host"] = []string{req.Host}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "HMAC-SHA256 SignedHeaders=x-ms-date;host;x-ms-content-sha256&Signature="+base64.StdEncoding.EncodeToString(signature))

	return nil
}

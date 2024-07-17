package web

import (
	"encoding/json"
	"io"
)

type ctxUserIdKey string

const UserIdKey ctxUserIdKey = "userId"

func DecodeJsonRequest(rc io.ReadCloser, data interface{}) error {
	defer rc.Close()
	decoder := json.NewDecoder(rc)
	return decoder.Decode(data)
}

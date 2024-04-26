package testutils

import (
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
)

func GenEncryptedPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func ReqStructToHTTPBody(reqStruct any) *bytes.Reader {
	reqBodyBytes, err := json.Marshal(reqStruct)
	if err != nil {
		panic(err)
	}
	reqBody := bytes.NewReader(reqBodyBytes)
	return reqBody
}

func StructToString(inputStruct any) string {
	inputStructBytes, err := json.Marshal(inputStruct)
	if err != nil {
		panic(err)
	}
	return string(inputStructBytes)
}

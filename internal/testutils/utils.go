package testutils

import (
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func GenEncryptedPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func GenHTTPJSONReq(method, url string, reqStruct any) *http.Request {
	req, err := http.NewRequest(method, url, structToHTTPBody(reqStruct))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func structToHTTPBody(reqStruct any) *bytes.Reader {
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

func StringToStruct(inputString string, outputStruct any) {
	err := json.Unmarshal([]byte(inputString), outputStruct)
	if err != nil {
		panic(err)
	}
}

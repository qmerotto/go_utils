package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct{}

type ClientHelper interface {
	Post(url string, body interface{}, headers http.Header, auth map[string]string) (*http.Response, error)
}

func (c *Client) Post(url string, body interface{}, headers http.Header, auth map[string]string) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	setAuthorizationOnRequest(request, auth)
	if err != nil {
		return nil, err
	}
	request.Header = mergeHeaders(request.Header, headers)
	client := &http.Client{}
	return client.Do(request)
}

func setAuthorizationOnRequest(request *http.Request, auth map[string]string) {
	if auth == nil {
		return
	}
	username, isUsernamePresent := auth["username"]
	password, isPasswordPresent := auth["password"]
	jwt, jwtExists := auth["jwt"]
	if isUsernamePresent && isPasswordPresent {
		request.SetBasicAuth(username, password)
		return
	}
	if jwtExists {
		request.Header = mergeHeaders(request.Header, http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", jwt)}})
	}
}

func mergeHeaders(existingHeaders http.Header, insertedHeaders http.Header) http.Header {
	if existingHeaders == nil {
		return insertedHeaders
	}
	for key, value := range insertedHeaders {
		_, headerExists := existingHeaders[key]
		if headerExists {
			continue
		}
		existingHeaders[key] = value
	}
	return existingHeaders
}

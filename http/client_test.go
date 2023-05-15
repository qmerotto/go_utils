package http

import (
	"net/http"
	"reflect"
	"testing"
)

func TestPost(t *testing.T) {
	client := Client{}
	body := map[string]string{"oui": "non"}
	headers := http.Header{"Content-Type": []string{"application/json"}}

	_, err := client.Post("https://mock-tikee-api.enlaps.fr", body, headers, nil)
	if err == nil {
		t.Fatalf("https://mock-tikee-api.enlaps.fr endpoint does not exist.")
	}

	_, err = client.Post("https://mock-tikee-api.enlaps.fr", "non_json_value", headers, nil)
	if err == nil {
		t.Fatalf("Post should take a json as body.")
	}
}

func TestMergeHeaders(t *testing.T) {
	existingHeaders := http.Header{"mockheader": []string{"mockvalue"}}
	insertedHeaders := http.Header{"mockheader": []string{"mockvalue2"}, "mynewmockheader": []string{"mockvalue"}}
	expectedHeaders := http.Header{"mockheader": []string{"mockvalue"}, "mynewmockheader": []string{"mockvalue"}}

	mergeHeaders(existingHeaders, insertedHeaders)

	if reflect.DeepEqual(existingHeaders, expectedHeaders) == false {
		t.Fatalf("Those headers should be equal")
	}
}

func TestSetAuthorizationOnRequest(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "https://mock-tikee-api.enlaps.fr", nil)
	if err != nil {
		t.Fatalf("Unable to initialize http request")
	}
	setAuthorizationOnRequest(request, nil)
	_, authorizationExist := request.Header["Authorization"]
	if authorizationExist == true {
		t.Fatalf("Authorization should not be initialized yet")
	}
	setAuthorizationOnRequest(request, map[string]string{"username": "mock_username", "password": "mock_password"})
	authorizationHeader, authorizationExist := request.Header["Authorization"]
	if authorizationExist == false {
		t.Fatalf("Authorization headers should have been initialized")
	}
	setAuthorizationOnRequest(request, map[string]string{"jwt": "mock_jwt"})
	jwtAuthorizationHeader, authorizationExist := request.Header["Authorization"]
	if authorizationExist == false {
		t.Fatalf("Authorization headers should have been initialized")
	}
	if authorizationHeader[0] != jwtAuthorizationHeader[0] {
		t.Fatalf("Authorization headers should not have been modified")
	}
}

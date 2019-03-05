package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// http recorder act like browser that accept the result of http request
	recorder := httptest.NewRecorder()

	// create http handler from our handlerFunc, "handler" is handler defined in main.go
	hf := http.HandlerFunc(handler)

	// serve http request to our recorder
	hf.ServeHTTP(recorder, req)

	// check status
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong http status code: got %v want %v",
			status, http.StatusOK)
	}

	// check response body
	expected := "Hello World!"
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}

}

func TestRouter(t *testing.T) {
	r := muxRouter()

	// Create a new server using the "httptest" libraries `NewServer` method
	// docs: https://golang.org/pkg/net/http/httptest/#NewServer
	mockServer := httptest.NewServer(r)

	res, err := http.Get(mockServer.URL + "/hello")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Status mut be ok, got %v", res.StatusCode)
	}

	defer res.Body.Close()

	// read the body
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	// bytes to string
	resString := string(b)
	expected := "Hello World!"

	if resString != expected {
		t.Errorf("Response should be %v got: %v", expected, resString)
	}
}

func TestRouterForNonExistentRoute(t *testing.T) {
	r := muxRouter()
	mockServer := httptest.NewServer(r)

	resp, err := http.Post(mockServer.URL+"/hello", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// status expected to be wrong
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status should be 405, got: %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)

	}

	respString := string(b)
	expected := ""
	if respString != expected {
		t.Errorf("Response should be %s got %s", expected, respString)
	}
}

func TestStaticFileServer(t *testing.T) {
	r := muxRouter()
	mockServer := httptest.NewServer(r)

	// hit the /assets endpoint
	resp, err := http.Get(mockServer.URL + "/assets")
	if err != nil {
		t.Fatal(err)
	}

	// expect to be ok
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got: %v", resp.StatusCode)
	}

	// test content type header is text/html
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Content type expected %s got %s", expectedContentType, contentType)
	}
}

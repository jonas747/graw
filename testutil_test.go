package graw

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestNewServerFromResponse(t *testing.T) {
	expected := []byte("10101010101___")
	server := newServerFromResponse(200, expected)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("request to server failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("request returned bad status: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		t.Fatalf("no body in response")
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if !reflect.DeepEqual(respBytes, expected) {
		t.Errorf("response incorrect; expected %s, got %s", respBytes, expected)
	}
}

func TestResponseIs(t *testing.T) {
	expected := []byte("ksjdnksbf")
	if responseIs(&http.Response{StatusCode: 200, Body: nil}, 200, expected) {
		t.Error("failed to identify nil body")
	}

	if !responseIs(&http.Response{StatusCode: 200, Body: nil}, 200, nil) {
		t.Error("failed to accept nil body with nil expectation")
	}

	if responseIs(&http.Response{
		StatusCode: 201,
		Body:       bytesCloser{bytes.NewBuffer(expected), nil},
	}, 200, expected) {
		t.Error("failed to identify status code difference")
	}

	if responseIs(&http.Response{
		StatusCode: 200,
		Body: bytesCloser{
			bytes.NewBuffer(expected),
			fmt.Errorf("AN ERROR"),
		},
	}, 200, expected) {
		t.Error("faulty read of response body did not become a diff")
	}

	if responseIs(&http.Response{
		StatusCode: 200,
		Body:       bytesCloser{bytes.NewBuffer(expected), nil},
	}, 200, []byte("sdfsdj")) {
		t.Error("body comparison failed; should have returned false")
	}

	if !responseIs(&http.Response{
		StatusCode: 200,
		Body:       bytesCloser{bytes.NewBuffer(expected), nil},
	}, 200, expected) {
		t.Error("body comparison failed; should have returned true")
	}
}
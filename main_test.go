package main

import (
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"
	"encoding/json"
	"fmt"
)

func TestGetRoot(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    w := httptest.NewRecorder()
    getRoot(w, req)
    res := w.Result()
    defer res.Body.Close()
    data, err := ioutil.ReadAll(res.Body)

    if err != nil {
        t.Errorf("expected error to be nil got %v", err)
    }

    if string(data) != "This is my website!\n" {
        t.Errorf("Wrong response!")
    }
}

func TestGetGithub(t *testing.T) {
    expectedUsername := "octocat"
    expectedName := "The Octocat"

    req := httptest.NewRequest(http.MethodGet, "/github/" + expectedUsername + "/repositories", nil)
    w := httptest.NewRecorder()
    getGithub(w, req)
    res := w.Result()
    defer res.Body.Close()
    data, err := ioutil.ReadAll(res.Body)

	if err != nil {
        fmt.Println("Error while reading the data", err.Error())
    }

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)

	if err != nil {
        fmt.Println("Error while decoding the data", err.Error())
    }
	
	if result["Username"].(string) != expectedUsername {
        t.Errorf("expected " + expectedUsername + ", got " + result["Username"].(string))
    }

	if result["Name"].(string) != expectedName {
        t.Errorf("expected " + expectedName + ", got " + result["Name"].(string))
    }
}
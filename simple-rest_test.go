package main

import (
	"testing"
	"net/http"
	"bytes"
	"./entities"
	"./fileStore"
	"encoding/json"
)

const url = "http://localhost:8080"
const test_prefix = "test_"

func TestCreateSite (t *testing.T) {
	defer RemoveTestData(t)
	var emptyAP = []entities.AccessPoint{}
	example_site := entities.Site{test_prefix + "1", "role1", "uri1", emptyAP}
	site_json, _ := example_site.ToJson()
	resp, err := http.Post(url + "/sites/example1", "application/json", bytes.NewBuffer(site_json))
	if err != nil {
		t.Error("Error running test: " + err.Error())
		return
	}

	var returned_site entities.Site
	err = json.NewDecoder(resp.Body).Decode(&returned_site)
	if err != nil {
		t.Error("Error running test: " + err.Error())
		return
	}

	if !returned_site.EqualTo(&example_site) {
		t.Error("Error creating site")
		return
	}
}

func RemoveTestData(t *testing.T) {
	fs := fileStore.FileStore{}
	fs.SetPrefix("./data/")
	err := fs.RemoveTestFiles()
	if err != nil {
		t.Error(err.Error())
	}
}

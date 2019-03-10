package main

import (
	"testing"
	"net/http"
	"bytes"
	"./entities"
	"./fileStore"
	"encoding/json"
	"fmt"
)

const url = "http://localhost:8080"
// Prefix allows us to select data to remove at end of testing.
const test_prefix = "test"

// Test:
// 	that a site can be created
//	that two sites with the same name can not be created
func TestCreateSite(t *testing.T) {
	fmt.Println("RUNNING: Test Create Site")
	// Always remove whatever testing data we created.
	defer RemoveTestData(t)
	var emptyAP = []entities.AccessPoint{}
	example_site := entities.Site{test_prefix + "one", test_prefix + "role1", test_prefix + "uri1", emptyAP}

	fmt.Println("\tCreating Site:", example_site)
	createTestSite(t, example_site, 200)

	fmt.Println("\tTrying create duplicate site. (Failure expected)")
	createTestSite(t, example_site, 400)
}

// Test:
//	that a site can be created and then edited
//	test that a site that does not exist can not be edited
func TestEditSite(t *testing.T) {
	fmt.Println("RUNNING: Test Update Site")
	// Always remove whatever testing data we created.
	defer RemoveTestData(t)
	var emptyAP = []entities.AccessPoint{}
	example_site := entities.Site{test_prefix + "two", test_prefix + "role1", test_prefix + "uri1", emptyAP}

	fmt.Println("\tCreating Site:", example_site)
	createTestSite(t, example_site, 200)

	example_site_update := entities.Site{test_prefix + "two", test_prefix + "role_update", test_prefix + "uri_update", emptyAP}
	fmt.Println("\tUpdating Site:", example_site_update)
	editTestSite(t, example_site_update, 200)

	fmt.Println("\tTrying to update nonexistant site. (Failure expected)")
	example_site_update_fake := entities.Site{test_prefix + "fake", test_prefix + "role_update", test_prefix + "uri_update", emptyAP}
	editTestSite(t, example_site_update_fake, 400)

	// TODO check access points arent updated.

}


// Test:
// 	that a site can be created and then deleted
//	that a site that was already deleted can not be deleted again
//	that a site that does not exist can not be deleted
func TestDeleteSite (t *testing.T) {
	fmt.Println("RUNNING: Test Delete Site")
	defer RemoveTestData(t)
	var emptyAP = []entities.AccessPoint{}
	example_site := entities.Site{test_prefix + "three", test_prefix + "role1", test_prefix + "uri1", emptyAP}

	fmt.Println("\tCreating Site:", example_site)
	createTestSite(t, example_site, 200)

	fmt.Println("\tDeleting Site:", example_site)
	deleteTestSite(t, example_site.Name, 200)

	fmt.Println("\tDeleting Site that has been deleted (Failure Expected)")
	deleteTestSite(t, example_site.Name, 400)

	fmt.Println("\tDeleting Site that never existed (Failure Expected)")
	deleteTestSite(t, "cats_r_cool", 400)
}

// =============== Helper functions ================= //
func createTestSite(t *testing.T, site entities.Site, expected_response_code int) {
	site_json, _ := site.ToJson()
	resp, err := http.Post(url + "/sites", "application/json", bytes.NewBuffer(site_json))
	if err != nil {
		t.Error("Error running test: " + err.Error())
		return
	}

	// Ensure we got the expected response code.
	if resp.StatusCode != expected_response_code {
		t.Error("Returned repsonse code:", resp.StatusCode, " does not match expected: ", expected_response_code)
		return
	}

	// We expect a successful response.
	if expected_response_code == 200 {
		var returned_site entities.Site
		err = json.NewDecoder(resp.Body).Decode(&returned_site)
		if err != nil {
			t.Error("Error running test: " + err.Error())
			return
		}

		// Ensure we were returned the data we sent
		if !returned_site.EqualTo(&site) {
			t.Error("Error creating site")
			return
		}

		// Check that the site exists in the filestore.
		getTestSite(t, site.Name, 200, site)
		fmt.Println("\t\tSite successfully created: ", site)
	} else {
		fmt.Println("\t\tExpected site creation failure: ", site)
	}
}

func editTestSite(t *testing.T, site entities.Site, expected_response_code int) {
	site_json, err := site.ToJson()
	if err != nil {
		t.Error("Error running test: " + err.Error())
		return
	}
	body := bytes.NewBuffer(site_json)
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", url + "/sites", body)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Error running test: " + err.Error())
		return
	}
	defer resp.Body.Close()

	// Ensure we got the expected response code.
	if resp.StatusCode != expected_response_code {
		t.Error("Returned repsonse code:", resp.StatusCode, " does not match expected: ", expected_response_code)
		return
	}

	// We expect a successful response.
	if expected_response_code == 200 {
		var returned_site entities.Site
		err = json.NewDecoder(resp.Body).Decode(&returned_site)
		if err != nil {
			t.Error("Error running test: " + err.Error())
			return
		}

		// Ensure we were returned the data we sent
		if !returned_site.EqualTo(&site) {
			t.Error("Error creating site")
			return
		}

		// Check that the site exists in the filestore and matches our edit.
		getTestSite(t, site.Name, 200, site)
		fmt.Println("\t\tSite successfully updated: ", site)
	} else {
		fmt.Println("\t\tExpected site update failure: ", site)
	}

}

func getTestSite(t *testing.T, site_name string, expected_response_code int, expected_response_site entities.Site) {
	resp, err := http.Get(url + "/sites/" + site_name)
	if err != nil {
		t.Error("Error running test: ", site_name)
		return
	}

	if resp.StatusCode != expected_response_code {
		t.Error("Returned repsonse code:", resp.StatusCode, " does not match expected: ", expected_response_code)
		return
	}

	if expected_response_code == 200 {
		// Load our site.
		var returned_site entities.Site
		err = json.NewDecoder(resp.Body).Decode(&returned_site)
		if err != nil {
			t.Error("Error running test: " + err.Error())
			return
		}

		// Ensure the correct site was returned.
		if !returned_site.EqualTo(&expected_response_site) {
			t.Error("Returned site: ", returned_site, " does not match expected: ", expected_response_site)
			return
		}
		fmt.Println("\t\tSite successfully got: ", site_name)
	} else {
		fmt.Println("\t\tExpected site get failure: ", site_name)
	}
}

func deleteTestSite(t *testing.T, site_name string, expected_response_code int) {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url + "/sites/" + site_name, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Error running test: " + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != expected_response_code {
		t.Error("Returned repsonse code:", resp.StatusCode, " does not match expected: ", expected_response_code)
		return
	}

	// We are expecting a successful delete.
	if expected_response_code == 200 {
		// Try to get the site to ensure it was removed from the file store.
		expected_response_site := entities.Site{}
		getTestSite(t, site_name, 400, expected_response_site)
		fmt.Println("\t\tSite successfully deleted: ", site_name)
	} else {
		fmt.Println("\t\tExpected site delete failure: ", site_name)
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

package main

import (
	"encoding/json"
	"net/http"
	//"log"
	//"fmt"
	//"errors"
	"github.com/gorilla/mux"
	"./fileStore"
	"./entities"
)

const FileStorePrefix = "./data/"

func main() {
	fs := fileStore.FileStore{}
	fs.SetPrefix(FileStorePrefix)
	router := mux.NewRouter()
	router.HandleFunc("/sites/create", CreateSite).Methods("POST")
	router.HandleFunc("/sites", GetSites).Methods("GET")
	router.HandleFunc("/sites/{name}", GetSite).Methods("GET")
	http.ListenAndServe(":8080", router)
}

func CreateSite(w http.ResponseWriter, r *http.Request) {
	var site entities.Site
	_ = json.NewDecoder(r.Body).Decode(&site)

	// Check if site exists in File Store.
	fs := fileStore.FileStore{}
	fs.SetPrefix(FileStorePrefix)
	exists := fs.Exists(site.Name)
	if exists {
		sendError(w, "A site already exists with this name")
	} else {
		// Doesn't exist, so we can create.
		site_json, err := site.ToJson()
		if err != nil {
			sendError(w, err.Error())
			return
		}
		// Write created object to our File Store.
		err = fs.Write(site.Name, site_json)
		if err != nil {
			sendError(w, err.Error())
			return
		} else {
			// Set the proper response code and return the created item.
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(site)
		}
	}
}

func GetSites(w http.ResponseWriter, r *http.Request) {
	// Get all Site names in the FileStore
	fs := fileStore.FileStore{}
	fs.SetPrefix(FileStorePrefix)
	site_names, err := fs.GetFiles()
	if err != nil {
		sendError(w, err.Error())
	} else {
		var sites []entities.Site
		// Load all site objects
		for _, site_name := range site_names {
			// Get File data.
			file_data, err := fs.Load(site_name)
			if err != nil {
				sendError(w, err.Error())
				return
			}
			// Build site object from file data.
			site, err := entities.SiteFromJson(file_data)
			if err != nil {
				sendError(w, err.Error())
				return
			}
			sites = append(sites, site)
		}
		json.NewEncoder(w).Encode(sites)
	}
}

func GetSite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Check if site exists in the File Store.
	fs := fileStore.FileStore{}
	fs.SetPrefix(FileStorePrefix)
	exists := fs.Exists(params["name"])

	if exists {
		// Get File data.
		file_data, err := fs.Load(params["name"])
		if err != nil {
			sendError(w, err.Error())
			return
		}
		// Build site object from file data.
		site, err := entities.SiteFromJson(file_data)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		json.NewEncoder(w).Encode(site)
	} else {
		sendError(w, "Site does not exist")
	}
}


func sendError(w http.ResponseWriter, msg string) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(entities.ErrorResponse{msg})
}

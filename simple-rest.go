package main

import (
	"encoding/json"
	"net/http"
	//"log"
	//"fmt"
	"errors"
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
	router.HandleFunc("/sites/{name}/accesspoints", GetAPs).Methods("GET")
	router.HandleFunc("/sites/{name}/accesspoints/{label}", APHandler).Methods("GET", "POST")
	http.ListenAndServe(":8080", router)
}

func WriteSiteToStore(site entities.Site) (error) {
	fs := fileStore.FileStore{}
	fs.SetPrefix(FileStorePrefix)

	site_json, err := site.ToJson()
	if err != nil {
		return err
	}
	// Write created object to our File Store.
	err = fs.Write(site.Name, site_json)
	if err != nil {
		return err
	}

	return nil
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
		err := WriteSiteToStore(site)
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

func GetSiteFromStore(w http.ResponseWriter, r *http.Request) (entities.Site, error) {
	params := mux.Vars(r)
	var site entities.Site
	var err error = nil

	// Check if site exists in the File Store.
	fs := fileStore.FileStore{}
	fs.SetPrefix(FileStorePrefix)
	exists := fs.Exists(params["name"])

	if exists {
		// Get File data.
		file_data, err := fs.Load(params["name"])
		if err != nil {
			return site, err
		}
		// Build site object from file data.
		site, err = entities.SiteFromJson(file_data)
		if err != nil {
			return site, err
		}
	} else {
		err = errors.New("Site does not exist")
	}

	return site, err
}

func GetSite(w http.ResponseWriter, r *http.Request) {
	site, err := GetSiteFromStore(w, r)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(site)
}

func GetAPs(w http.ResponseWriter, r *http.Request) {
	site, err := GetSiteFromStore(w, r)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(site.Access_points)
}

func GetAP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	site, err := GetSiteFromStore(w, r)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	// Find the accesspoint
	var ap entities.AccessPoint
	for _, site_ap := range site.Access_points {
		if site_ap.Label == params["label"] {
			ap = site_ap
			break
		}
	}

	// ap doesn't exist
	if (entities.AccessPoint{}) == ap {
		sendError(w, "Access point does not exist")
		return
	}

	json.NewEncoder(w).Encode(ap)
}

func CreateAP(w http.ResponseWriter, r *http.Request) {
	// Get the site
	site, err := GetSiteFromStore(w, r)
	if err != nil {
		sendError(w, err.Error())
	}

	// Parse the access point
	var ap entities.AccessPoint
	_ = json.NewDecoder(r.Body).Decode(&ap)

	// Check for accesspoint label
	found := 0
	for i := 0; i < len(site.Access_points); i++ {
		// If label already exists, update Url
		if site.Access_points[i].Label == ap.Label {
			found = 1
			site.Access_points[i].Url = ap.Url
			break
		}
	}

	// If not found, then append to accesspoint list
	if found == 0 {
		site.Access_points = append(site.Access_points, ap)
	}

	// Rewrite entire site to file - I think this is easier than piece-wise update
	err = WriteSiteToStore(site)
	if err != nil {
		sendError(w, err.Error())
		return
	} else {
		// Set the proper response code and return the created item.
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(ap)
	}
}

func APHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || len(r.Method) == 0 {
		GetAP(w, r)
	} else if r.Method == "POST" {
		CreateAP(w, r)
	} else {
		return
	}
}

func sendError(w http.ResponseWriter, msg string) {
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(entities.ErrorResponse{msg})
}
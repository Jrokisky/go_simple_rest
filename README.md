# CMSC 621 Project 1

## Team
	* Frank Serna
	* Justin Rokisky

## Project Details

### Setup Instructions
	* Install Gorilla Mux: go get github.com/gorilla/mux
	* Ensure a directory exists with the name of the File Prefix in 'simple-rest.go'

### Architecture
	* We store data using our file store package
		* each site is stored in its own file under the File Prefix
		* this allows for easy existance checks as well as limits the size of data that needs to be written to disk on updates.
		* this also allows less conflict if multiple operations are done concurrently
	* We use the /sites endpoint for:
		* GET => gets all sites
		* PUT => updates a given site
		* POST => creates a new site
	* We use the /sites/$SITENAME endpoint for:
		* GET => get site identified by $SITENAME
		* DELETE => delete the site identified by $SITENAME
	* We use the /sites/$SITENAME/accesspoints for:
		* GET => gets all access points of the given site
		* PUT => updates a given access point
		* POST => creates an access point on the given site
	* We use the /sites/$SITENAME/accesspoints/$APLABEL
		* GET => gets the accesspoint identified by $APLABLE from the site identified by $SITENAME
		* DELETE => deletes the accesspoint identified by $APLABLE from the site identified by $SITENAME

### Assumptions Made
We assumed that the goal of this project was to focus primarily on build a REST API that used JSON for messaging.
Due to this, the following reasonable assumptions were made:
	* Site names would be limited to a single string of lower case letters
	* 

### Testing Instructions
We added a test suite to help test our application.
To test:
	* run the application
	* in a separate terminal run: go test

# CMSC 621 Project 1

## Team: Frank Serna & Justin Rokisky

## About
This software provides a RESTful API for interacting with sites and access points.  Sites have the following properties:

```go
type Site struct {
	Name string
	Role string
	Uri string
	Access_points []AccessPoint
}
```
and access points have the following properties:

```go
type AccessPoint struct {
	Label string
	Url string
}
```

Data is persisted by storing each site in json form in a file named site.Name. This allows for easy existance checks, limits the size of data that needs to be written to disk on updates and also allows less conflict if multiple operations are done concurrently


Interaction is provided by GET, POST, PUT, DELETE commands explained below.

### Running the server
In a terminal, type
```bash
go run simple-rest.go
```

### Running the test suite
* Run the application using the instructions above.
* In a separate terminal, type
```bash
go test
```

### RESTful API commands
#### GET requests
GET requests provide viewable JSON properties of the object fetched.  These can be sites or access points.  Examples using a browser follow:
* View all sites: 
```bash
http://localhost:8080/sites
```
* View a specific site named $SITE_NAME:
```bash
http://localhost:8080/sites/$SITE_NAME
```
* View all accesspoints to a specific $SITE_NAME:
```bash
http://localhost:8080/sites/$SITE_NAME/accesspoints
```
* View a specific access point $AP_NAME at site $SITE_NAME:
```bash
http://localhost:8080/sites/$SITE_NAME/accesspoints/$AP_NAME
```
#### POST requests
POST requests create or update a site or access point object.  These are submitted via JSON.  Because name and label designate the site and accesspoint id, POST creates only the specified JSON object if it does not exist, otherwise it updates the object with the same resource id.  Examples using curl are as follows:
* POST a new site foo with empty access points:
```bash
curl -d '{"Name":"foo","Role":"cat","Uri":"karate","Access_points":null}' -H "Content-Type: application/json" http://localhost:8080/sites
```
* POST a new site foo with access points:
```bash
curl -d '{"Name":"foo","Role":"cat","Uri":"karate","Access_points":[{"Label":"foo","Url":"bar"},{"Label":"baz","Url":"beep"}]}' -H "Content-Type: application/json" http://localhost:8080/sites
```
* POST a new accesspoint dog to site foo:
```bash
curl -d '{"Label":"dog","Url":"cat"}' -H "Content-Type: application/json" http://localhost:8080/sites/foo/accesspoints
```
* POST update existing accesspoint dog to site foo with a new url tiger:
```bash
curl -d '{"Label":"dog","Url":"tiger"}' -H "Content-Type: application/json" http://localhost:8080/sites/foo/accesspoints
```
#### PUT requests
For this API, PUT and POST requests are synonymous as the resource id (name for site and label for accesspoint) is assumed to be provided by the user.  Please see the POST requests section.
#### DELETE requests
DELETE requests remove site or accesspoint objects.  Since each site stores access points, a site deletion will delete all access points associated with the site.  Examples follow:
* DELETE a site named test:
```bash
curl -X "DELETE" http://localhost:8080/sites/test
```
* DELETE an accesspoint dog on site test:
```bash
curl -X "DELETE" http://localhost:8080/sites/test/accesspoints/dog
```

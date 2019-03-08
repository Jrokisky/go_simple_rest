package entities

import (
	"encoding/json"
)

type Site struct {
	Name string
	Role string
	Uri string
	Access_points []AccessPoint
}

type AccessPoint struct {
	Label string
	Url string
}

type ErrorResponse struct {
	Error string
}

type SuccessResponse struct {
	Success string
}

// TODO: how does this handle access points?
func (s *Site) ToJson() ([]byte, error) {
	json, err := json.Marshal(s)
	return json, err
}

func (ap *AccessPoint) ToJson() ([]byte, error) {
	json, err := json.Marshal(ap)
	return json, err
}

func SiteFromJson(json_data []byte) (Site, error) {
	var site Site
	err := json.Unmarshal(json_data, &site)
	return site, err
}

func AccessPointFromJson(json_data []byte) (AccessPoint, error) {
	var accessPoint AccessPoint
	err := json.Unmarshal(json_data, &accessPoint)
	return accessPoint, err
}

// TODO: validation to ensure Sites cant have duplicate names
// TODO: validation to ensure that Sites cant have accesspoints with the same labels

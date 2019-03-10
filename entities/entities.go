package entities

import (
	"encoding/json"
	"regexp"
	"errors"
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

func (s *Site) EqualTo(s2 *Site, ignore_access_points bool) (bool) {
	if ignore_access_points {
		var emptyAP = []AccessPoint{}
		s.Access_points = emptyAP
		s2.Access_points = emptyAP
	}
	s_json, err := s.ToJson()
	if err != nil {
		return false
	}

	s2_json, err2 := s2.ToJson()
	if err2 != nil {
		return false
	}

	return string(s_json) == string(s2_json)
}

func (ap *AccessPoint) EqualTo(ap2 *AccessPoint) (bool) {
	ap_json, err := ap.ToJson()
	if err != nil {
		return false
	}

	ap2_json, err2 := ap2.ToJson()
	if err2 != nil {
		return false
	}

	return string(ap_json) == string(ap2_json)
}

func (s *Site) Validate() (error) {
	isAlpha := regexp.MustCompile(`^[a-z]+$`).MatchString
	if !isAlpha(s.Name) {
		return errors.New("Site name can only contain lowercase letters")
	}

	apLabels := make(map[string]int)
	for _, ap := range s.Access_points {
		if apLabels[ap.Label] == 1 {
			return errors.New("Access Point labels must be unique")
		} else {
			apLabels[ap.Label] = 1
		}
	}
	return nil
}

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

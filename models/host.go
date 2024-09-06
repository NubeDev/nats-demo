package models

type Host struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

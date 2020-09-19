package models

type ToDo struct {
	Id          string `json:"Id" bson:"_id,omitempty"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Completed   bool   `json:"Completed"`
}

package models

type Orders struct {
	Status string `json:"status" bson:"status"`
	Data   string `json:"data" bson:"data"`
}

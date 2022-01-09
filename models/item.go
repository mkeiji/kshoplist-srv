package models

type Item struct {
	Id      int    `json:"id"`
	StoreId int    `json:"storeId"`
	Name    string `json:"name"`
}

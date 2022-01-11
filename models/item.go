package models

import "time"

type Item struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	StoreId   int       `json:"storeId"`
	Name      string    `json:"name"`
}

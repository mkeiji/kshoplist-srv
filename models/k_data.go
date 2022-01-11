package models

import "kshoplistSrv/enums"

type Kdata struct {
	Type   enums.KmsgType `json:"type"`
	Action enums.Actions  `json:"action"`
	Items  []Item         `json:"items"`
}

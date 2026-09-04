package model

import (
	"gorm.io/datatypes"
)

type Setting struct {
	Key   string         `gorm:"primaryKey" json:"key"`
	Value datatypes.JSON `json:"value"`
}

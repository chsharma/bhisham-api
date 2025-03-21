package models

import "time"

type Bhisham struct {
	ID           int        `json:"id" db:"id"`
	SerialNo     string     `json:"serial_no" db:"serial_no"`
	BhishamName  *string    `json:"bhisham_name,omitempty" db:"bhisham_name"`
	CreatedBy    *string    `json:"created_by,omitempty" db:"created_by"`
	CreatedAt    *time.Time `json:"created_at,omitempty" db:"created_at"`
	IsComplete   *int       `json:"is_complete,omitempty" db:"is_complete"`
	CompleteBy   *string    `json:"complete_by,omitempty" db:"complete_by"`
	CompleteTime *string    `json:"complete_time,omitempty" db:"complete_time"`
}

type BhishamMapping struct {
	BhishamID  int    `json:"bhisham_id"`
	MCNo       int    `json:"mc_no"`
	MCName     string `json:"mc_name"`
	CubeNumber int    `json:"cube_number"`
	CCNo       string `json:"cc_no"`
	CCName     string `json:"cc_name"`
}

type BhishamKit struct {
	BhishamID int    `json:"bhisham_id"`
	MCNo      int    `json:"mc_no"`
	CCNo      string `json:"cc_no"`
	KitCode   string `json:"kitcode"`
	KitName   string `json:"kitname"`
	NoOfKit   int    `json:"no_of_kit"`
}

type GetChildCube struct {
	BhishamID int `json:"bhisham_id"`
	MCNo      int `json:"mc_no"`
}

type GetKits struct {
	BhishamID  int    `json:"bhisham_id"`
	MCNo       int    `json:"mc_no"`
	CubeNumber string `json:"cc_no"`
}

package models

import "time"

type Bhisham struct {
	ID           int        `json:"id"`
	SerialNo     string     `json:"serial_no"`
	BhishamName  *string    `json:"bhisham_name"`
	CreatedBy    *string    `json:"created_by"`
	CreatedAt    *time.Time `json:"created_at"`
	IsComplete   *int       `json:"is_complete"`
	CompleteBy   *string    `json:"complete_by"`
	CompleteTime *string    `json:"complete_time"`
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

type BhishamMappingData struct {
	BhishamID      int
	SerialNo       string
	MCNo           int
	CCNo           string
	CCName         string
	KitCode        string
	KitName        string
	NoOfKit        int
	SKUCode        string
	ItemName       string
	BatchNoSrNo    string
	MFD            string
	EXP            string
	ManufacturedBy string
	TotalQty       int
	CCEPC          string
	MCEPC          string
	MCName         string
	NoOfItem       int
	IsCube         int
	CubeNumber     int
}

type GetBhishamID struct {
	BhishamID int `json:"bhisham_id"`
}

type UpdateBhishamData struct {
	BhishamID  int    `json:"bhisham_id"`
	MCNo       int    `json:"mc_no"`
	CubeNumber int    `json:"cube_number"`
	KitName    string `json:"kit_name"`
	BatchCode  string `json:"batch_code"`
	MFD        string `json:"mfd"`
	EXP        string `json:"exp"`
	ID         int    `json:"id"`
	UpdateType int    `json:"update_typeid"`
}

type KitItems struct {
	ID             int    `json:"id"`
	MCNo           int    `json:"mc_no"`
	CubeNumber     int    `json:"cube_number"`
	KitName        string `json:"kit_name"`
	KitNo          int    `json:"kit_no"`
	BatchNoSrNo    string `json:"batch_no_sr_no"`
	SKUName        string `json:"sku_name"`
	Mfd            string `json:"mfd"`
	Exp            string `json:"exp"`
	ManufacturedBy string `json:"manufactured_by"`
	SKUQty         int    `json:"sku_qty"`
}

type BhishamData struct {
	ID             int    `json:"id"`
	BhishamID      int    `json:"bhisham_id"`
	MCNo           int    `json:"mc_no"`
	MCName         string `json:"mc_name"`
	MCEPC          string `json:"mc_epc"`
	CCNo           string `json:"cc_no"`
	CCName         string `json:"cc_name"`
	CCEPC          string `json:"cc_epc"`
	KitCode        string `json:"kitcode"`
	KitNo          int    `json:"kit_no"`
	KitEPC         string `json:"kit_epc"`
	KitBatchNo     string `json:"kit_batch_no"`
	KitExpiry      string `json:"kit_expiry"`
	KitQty         int    `json:"kit_qty"`
	SKUCode        string `json:"sku_code"`
	SKUName        string `json:"sku_name"`
	BatchNoSrNo    string `json:"batch_no_sr_no"`
	Mfd            string `json:"mfd"`
	Exp            string `json:"exp"`
	ManufacturedBy string `json:"manufactured_by"`
	SKUQty         int    `json:"sku_qty"`
	CubeNumber     int    `json:"cube_number"`
	KitName        string `json:"kitname"`
	NoOfKit        int    `json:"no_of_kit"`
}

type UpdateType struct {
	UpdateTypeID int    `json:"update_typeid"`
	Name         string `json:"name"`
}

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
	IsHHSynch    *int       `json:"is_hh_synch"`
	HHSerial     *string    `json:"hh_serial"`
	HHSynchTime  *string    `json:"hh_synch_time"`
	HHSynchCount *int       `json:"hh_synch_count"`
}

type BhishamMapping struct {
	BhishamID   int    `json:"bhisham_id"`
	MCNo        int    `json:"mc_no"`
	MCName      string `json:"mc_name"`
	CubeNumber  int    `json:"cube_number"`
	CCNo        string `json:"cc_no"`
	CCName      string `json:"cc_name"`
	Total       int    `json:"total_item"`
	TotalUpdate int    `json:"total_update_item"`
}

type BhishamKit struct {
	BhishamID   int    `json:"bhisham_id"`
	MCNo        int    `json:"mc_no"`
	CCNo        string `json:"cc_no"`
	KitCode     string `json:"kitcode"`
	KitName     string `json:"kitname"`
	KitSlug     string `json:"kit_slug"`
	NoOfKit     int    `json:"no_of_kit"`
	KitExpiry   string `json:"kit_expiry"`
	Total       int    `json:"total_item"`
	TotalUpdate int    `json:"total_update_item"`
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
	KitSlug        string
	SKUSlug        string
}

type GetBhishamID struct {
	BhishamID int `json:"bhisham_id"`
}

type UpdateBhishamData struct {
	BhishamID  int    `json:"bhisham_id"`
	MCNo       int    `json:"mc_no"`
	CubeNumber int    `json:"cube_number"`
	KitCode    string `json:"kit_code"`
	KitSlug    string `json:"kit_slug"`
	SkuCode    string `json:"sku_code"`
	SkuSlug    string `json:"sku_slug"`
	BatchCode  string `json:"batch_code"`
	MFD        string `json:"mfd"`
	EXP        string `json:"exp"`
	ID         int    `json:"id"`
	UpdateType int    `json:"update_typeid"`
}

type KitItems struct {
	ID             int     `json:"id"`
	MCNo           int     `json:"mc_no"`
	CubeNumber     int     `json:"cube_number"`
	KitCode        string  `json:"kit_code"`
	KitName        string  `json:"kit_name"`
	KitSlug        string  `json:"kit_slug"`
	KitNo          int     `json:"kit_no"`
	BatchNoSrNo    string  `json:"batch_no_sr_no"`
	SKUName        string  `json:"sku_name"`
	SKUCode        string  `json:"sku_code"`
	SKUSlug        string  `json:"sku_slug"`
	Mfd            string  `json:"mfd"`
	Exp            string  `json:"exp"`
	ManufacturedBy string  `json:"manufactured_by"`
	SKUQty         int     `json:"sku_qty"`
	IsUpdate       *string `json:"is_update"`
	UpdateTime     *string `json:"update_time"`
	UpdateBy       *string `json:"updated_by"`
}

type BhishamData struct {
	ID             int     `json:"id"`
	BhishamID      int     `json:"bhisham_id"`
	MCNo           int     `json:"mc_no"`
	MCName         string  `json:"mc_name"`
	MCEPC          string  `json:"mc_epc"`
	CCNo           string  `json:"cc_no"`
	CCName         string  `json:"cc_name"`
	CCEPC          string  `json:"cc_epc"`
	KitCode        string  `json:"kitcode"`
	KitNo          int     `json:"kit_no"`
	KitEPC         string  `json:"kit_epc"`
	KitBatchNo     string  `json:"kit_batch_no"`
	KitExpiry      string  `json:"kit_expiry"`
	KitQty         int     `json:"kit_qty"`
	SKUCode        string  `json:"sku_code"`
	SKUName        string  `json:"sku_name"`
	BatchNoSrNo    string  `json:"batch_no_sr_no"`
	Mfd            string  `json:"mfd"`
	Exp            string  `json:"exp"`
	ManufacturedBy string  `json:"manufactured_by"`
	SKUQty         int     `json:"sku_qty"`
	CubeNumber     int     `json:"cube_number"`
	KitName        string  `json:"kitname"`
	NoOfKit        int     `json:"no_of_kit"`
	IsUpdate       *string `json:"is_update"`
	UpdateTime     *string `json:"update_time"`
	UpdateBy       *string `json:"updated_by"`
}

type UpdateType struct {
	UpdateTypeID int    `json:"update_typeid"`
	Name         string `json:"name"`
}

type SerialNo struct {
	SerialNo string `json:"loginid"`
}

type ReportBhishamMapping struct {
	McName         string `json:"mc_name"`
	CcNo           string `json:"cc_no"`
	CcName         string `json:"cc_name"`
	KitName        string `json:"kitname"`
	NoOfKit        int    `json:"no_of_kit"`
	SkuCode        string `json:"sku_code"`
	SkuName        string `json:"sku_name"`
	BatchNoSrNo    string `json:"batch_no_sr_no"`
	Mfd            string `json:"mfd"`
	Exp            string `json:"exp"`
	ManufacturedBy string `json:"manufactured_by"`
	NoOfItem       int    `json:"no_of_item"`
}

package repositories

import (
	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/models"
	"database/sql"
	"errors"
)

type BhishamRepository struct {
	DB *sql.DB
}

func (r *BhishamRepository) CreateBhisham(bhisham models.Bhisham) (map[string]interface{}, error) {
	// Validate SerialNo
	if bhisham.SerialNo == "" {
		return helper.CreateDynamicResponse("Serial No is required", false, nil, 400, nil), errors.New("serial_no cannot be empty")
	}

	// Begin Transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return helper.CreateDynamicResponse("Failed to start transaction", false, nil, 500, nil), err
	}

	// Rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if SerialNo already exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM bhisham WHERE serial_no = $1)", bhisham.SerialNo).Scan(&exists)
	if err != nil {
		return helper.CreateDynamicResponse("Error checking existing Serial No", false, nil, 500, nil), err
	}
	if exists {
		return helper.CreateDynamicResponse("Serial No already exists", false, nil, 400, nil), errors.New("duplicate serial_no")
	}

	// Insert New Record and Return ID
	var newID int
	insertBhishamQuery := `INSERT INTO bhisham (serial_no, bhisham_name, created_by, created_at)
	                       VALUES ($1, $2, $3, NOW()) RETURNING id`

	err = tx.QueryRow(insertBhishamQuery, bhisham.SerialNo, bhisham.BhishamName, bhisham.CreatedBy).Scan(&newID)
	if err != nil {
		return helper.CreateDynamicResponse("Error Creating Bhisham", false, nil, 400, nil), err
	}

	// Insert into bhisham_mapping from default_bhisham
	insertMappingQuery := `INSERT INTO bhisham_mapping (
		bhisham_id, serial_no, mc_no, cc_no, cc_name, kitcode, kitname, no_of_kit, sku_code, 
		item_name, batch_no_sr_no, mfd, exp, manufactured_by, total_qty, cc_epc, mc_epc, 
		mc_name, no_of_item, is_cube, cube_number
	) SELECT $1, $2, mc_no, cc_no, cc_name, kitcode, kitname, no_of_kit, sku_code, 
		item_name, batch_no_sr_no, mfd, exp, manufactured_by, total_qty, cc_epc, mc_epc, 
		mc_name, no_of_item, is_cube, cube_number 
	FROM public.default_bhisham ORDER BY mc_no, cube_number;`

	_, err = tx.Exec(insertMappingQuery, newID, bhisham.SerialNo)
	if err != nil {
		return helper.CreateDynamicResponse("Error inserting into Bhisham Mapping", false, nil, 400, nil), err
	}

	// Commit Transaction
	err = tx.Commit()
	if err != nil {
		return helper.CreateDynamicResponse("Transaction commit failed", false, nil, 500, nil), err
	}

	// Success Response with Inserted ID
	return helper.CreateDynamicResponse("Bhisham Created Successfully", true, map[string]interface{}{"id": newID}, 200, nil), nil
}

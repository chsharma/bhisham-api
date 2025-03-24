package repositories

import (
	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/models"
	"bhisham-api/internal/app/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type BhishamRepository struct {
	DB *sql.DB
}

func (r *BhishamRepository) oldCreateBhisham(bhisham models.Bhisham) (map[string]interface{}, error) {
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

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte("00000"), bcrypt.DefaultCost)
	if err != nil {
		return helper.CreateDynamicResponse("Error hashing password", false, nil, 500, nil), err
	}

	var user_id = utils.GenerateId()

	query := `INSERT INTO user_login (user_id, name, login_id, pwd, role_id, created_at) 
              VALUES ($1, $2, $3, $4, $5, NOW())`

	_, err = tx.Exec(query, user_id, bhisham.BhishamName, bhisham.SerialNo, string(hashedPwd), 3)
	if err != nil {
		return helper.CreateDynamicResponse("Error creating user", false, nil, 500, nil), err
	}

	// Commit Transaction
	err = tx.Commit()
	if err != nil {
		return helper.CreateDynamicResponse("Transaction commit failed", false, nil, 500, nil), err
	}

	// Success Response with Inserted ID
	return helper.CreateDynamicResponse("Bhisham Created Successfully", true, map[string]interface{}{"id": newID}, 200, nil), nil
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

	// Generate 8-digit bhisham_id with leading zeros
	bhishamIdPadded := fmt.Sprintf("%08d", newID)

	// Insert into bhisham_mapping from default_bhisham with formatted EPCs
	insertMappingQuery := `INSERT INTO bhisham_mapping (
        bhisham_id, serial_no, mc_no, cc_no, cc_name, kitcode, kitname, no_of_kit, sku_code, 
        item_name, batch_no_sr_no, mfd, exp, manufactured_by, total_qty, cc_epc, mc_epc, 
        mc_name, no_of_item, is_cube, cube_number
    ) SELECT $1, $2, mc_no, cc_no, cc_name, kitcode, kitname, no_of_kit, sku_code, 
        item_name, batch_no_sr_no, mfd, exp, manufactured_by, total_qty, 
        'CA' || mc_no || $3 || '00000000000' || LPAD(cube_number::text, 2, '0'),  -- CC EPC
        'A0' || mc_no || $3 || '0000000000000',  -- MC EPC
        mc_name, no_of_item, is_cube, cube_number 
    FROM public.default_bhisham ORDER BY mc_no, cube_number;`

	_, err = tx.Exec(insertMappingQuery, newID, bhisham.SerialNo, bhishamIdPadded)
	if err != nil {
		return helper.CreateDynamicResponse("Error inserting into Bhisham Mapping", false, nil, 400, nil), err
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte("00000"), bcrypt.DefaultCost)
	if err != nil {
		return helper.CreateDynamicResponse("Error hashing password", false, nil, 500, nil), err
	}

	Updatequery := `UPDATE public.bhisham_mapping 
	SET kit_expiry = temp.min_expiry
	FROM (
		SELECT cc_no, kitname, MIN(CAST(exp AS DATE)) AS min_expiry 
		FROM public.bhisham_mapping 
		WHERE LENGTH(exp) > 5 AND bhisham_id = $1
		GROUP BY cc_no, kitname
	) temp  
	WHERE public.bhisham_mapping.kitname = temp.kitname 
	AND public.bhisham_mapping.cc_no = temp.cc_no
	AND public.bhisham_mapping.bhisham_id = $1`

	_, err = tx.Exec(Updatequery, newID)
	if err != nil {
		return helper.CreateDynamicResponse("Error updating expiry >"+err.Error(), false, nil, 500, nil), err
	}

	var user_id = utils.GenerateId()

	query := `INSERT INTO user_login (user_id, name, login_id, pwd, role_id, created_at) 
              VALUES ($1, $2, $3, $4, $5, NOW())`

	_, err = tx.Exec(query, user_id, bhisham.BhishamName, bhisham.SerialNo, string(hashedPwd), 3)
	if err != nil {
		return helper.CreateDynamicResponse("Error creating user >"+err.Error(), false, nil, 500, nil), err
	}

	// Commit Transaction
	err = tx.Commit()
	if err != nil {
		return helper.CreateDynamicResponse("Transaction commit failed", false, nil, 500, nil), err
	}

	// Success Response with Inserted ID
	return helper.CreateDynamicResponse("Bhisham Created Successfully", true, map[string]interface{}{"id": newID}, 200, nil), nil
}

func (r *BhishamRepository) CreateBhishamData(BhishamID int, UserID string) (map[string]interface{}, error) {

	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM bhisham_data WHERE id = $1)", BhishamID).Scan(&exists)
	if err != nil {
		return helper.CreateDynamicResponse("Error checking existing Serial No", false, nil, 500, nil), err
	}
	if exists {
		return helper.CreateDynamicResponse("Data already exists", false, nil, 500, nil), errors.New("duplicate serial_no")
	}

	query := `SELECT bhisham_id, serial_no, mc_no, cc_no, cc_name, kitcode, kitname, no_of_kit, 
	sku_code, item_name, batch_no_sr_no, mfd, exp, manufactured_by, total_qty, 
	cc_epc, mc_epc, mc_name, no_of_item, is_cube, cube_number FROM bhisham_mapping WHERE bhisham_id=$1`

	rows, err := r.DB.Query(query, BhishamID)
	if err != nil {
		return helper.CreateDynamicResponse("Error  executing query > "+err.Error(), false, nil, 200, nil), nil
	}
	defer rows.Close()

	var mappings []models.BhishamMappingData
	uniqueCCNos := make(map[string]bool)
	uniqueCCKit := make(map[string]bool)
	var Cubes []string

	// Read all records
	for rows.Next() {
		var data models.BhishamMappingData
		if err := rows.Scan(&data.BhishamID, &data.SerialNo, &data.MCNo, &data.CCNo, &data.CCName, &data.KitCode, &data.KitName, &data.NoOfKit,
			&data.SKUCode, &data.ItemName, &data.BatchNoSrNo, &data.MFD, &data.EXP, &data.ManufacturedBy, &data.TotalQty,
			&data.CCEPC, &data.MCEPC, &data.MCName, &data.NoOfItem, &data.IsCube, &data.CubeNumber); err != nil {

			return helper.CreateDynamicResponse("Error  scanning row > "+err.Error(), false, nil, 200, nil), nil
		}
		mappings = append(mappings, data)

		// Unique CCNo values
		if _, exists := uniqueCCNos[data.CCNo]; !exists {
			uniqueCCNos[data.CCNo] = true
			Cubes = append(Cubes, data.CCNo)
		}
	}

	insertQuery := `INSERT INTO bhisham_data 
	(bhisham_id, mc_no, mc_name, mc_epc, cc_no, cc_name, cc_epc, kitcode, kit_no, kit_epc, kit_batch_no, 
	kit_expiry, kit_qty, sku_code, sku_name, batch_no_sr_no, mfd, exp, manufactured_by, sku_qty, cube_number,kitname,no_of_kit) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)`

	// Use a transaction for batch insertion
	tx, err := r.DB.Begin()
	if err != nil {
		return helper.CreateDynamicResponse("Error starting transaction > "+err.Error(), false, nil, 200, nil), nil
	}
	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		tx.Rollback()
		return helper.CreateDynamicResponse("Error preparing insert statement > "+err.Error(), false, nil, 200, nil), nil
	}
	defer stmt.Close()

	for _, ccNo := range Cubes {
		kitNo := 1 // Start Kit Numbering

		for _, mapping := range mappings {
			if mapping.CCNo == ccNo {
				key := fmt.Sprintf("%s-%s", mapping.CCNo, mapping.KitName)
				if _, exists := uniqueCCKit[key]; !exists {
					uniqueCCKit[key] = true
					KitItems := GetProductsByKit(mappings, mapping.KitName, mapping.CCNo)
					minExpiry := FindMinExpiry(KitItems)

					for k := 0; k < int(mapping.NoOfKit); k++ {
						for _, item := range KitItems {
							_, err := stmt.Exec(
								item.BhishamID, item.MCNo, item.MCName, item.MCEPC,
								item.CCNo, item.CCName, item.CCEPC, item.KitCode, kitNo, GenerateKitEPC(BhishamID, item.MCNo, item.CubeNumber, kitNo),
								item.BatchNoSrNo, minExpiry, item.NoOfKit, item.SKUCode, item.ItemName,
								item.BatchNoSrNo, item.MFD, item.EXP, item.ManufacturedBy, item.TotalQty, item.CubeNumber, item.KitName, item.NoOfKit,
							)
							if err != nil {
								tx.Rollback()
								return helper.CreateDynamicResponse("Error inserting data > "+err.Error(), false, nil, 200, nil), nil
							}
						}
						kitNo++
					}
				}
			}
		}
	}
	var ID int
	updateBhishamQuery := `UPDATE public.bhisham SET is_complete=1, complete_time=NOW(), complete_by=$1 WHERE id=$2 RETURNING id`

	err = tx.QueryRow(updateBhishamQuery, UserID, BhishamID).Scan(&ID)
	if err != nil {
		return helper.CreateDynamicResponse("Error Creating Bhisham >"+err.Error(), false, nil, 400, nil), err
	}
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return helper.CreateDynamicResponse("Error committing transaction > "+err.Error(), false, nil, 200, nil), nil
	}

	return helper.CreateDynamicResponse("Bhisham Created Successfully", true, nil, 200, nil), nil
}

func (r *BhishamRepository) UpdateBhishamData(obj models.UpdateBhishamData, UserID string) (map[string]interface{}, error) {
	var updateBhishamQuery, updateBhishamMappingQuery string
	var queryParams []interface{}
	var ParamsMapping []interface{}

	// Determine the update queries based on UpdateType
	switch obj.UpdateType {
	case 1:
		updateBhishamQuery = `UPDATE public.bhisham_data 
						  SET mfd=$1, exp=$2, batch_no_sr_no=$3, is_update=1, update_time=NOW(), updated_by=$7
						  WHERE bhisham_id=$4 AND sku_name=$5 AND mc_no=$6`
		updateBhishamMappingQuery = `UPDATE public.bhisham_mapping 
						  SET mfd=$1, exp=$2, batch_no_sr_no=$3, is_update=1, update_time=NOW(), updated_by=$7
						  WHERE bhisham_id=$4 AND item_name=$5 AND mc_no=$6`
		queryParams = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.BhishamID, obj.KitName, obj.MCNo, UserID}

	case 2:
		updateBhishamQuery = `UPDATE public.bhisham_data 
						  SET mfd=$1, exp=$2, batch_no_sr_no=$3, is_update=1, update_time=NOW(), updated_by=$8
						  WHERE bhisham_id=$4 AND sku_name=$5 AND mc_no=$6 AND cube_number=$7`
		updateBhishamMappingQuery = `UPDATE public.bhisham_mapping 
						  SET mfd=$1, exp=$2, batch_no_sr_no=$3, is_update=1, update_time=NOW(), updated_by=$8
						  WHERE bhisham_id=$4 AND item_name=$5 AND mc_no=$6 AND cube_number=$7`
		queryParams = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.BhishamID, obj.KitName, obj.MCNo, obj.CubeNumber, UserID}

	case 3:
		updateBhishamQuery = `UPDATE public.bhisham_data 
						  SET mfd=$1, exp=$2, batch_no_sr_no=$3, is_update=1, update_time=NOW(), updated_by=$5
						  WHERE id=$4`
		updateBhishamMappingQuery = `UPDATE public.bhisham_mapping 
						  SET mfd=$1, exp=$2, batch_no_sr_no=$3, is_update=1, update_time=NOW(), updated_by=$8
						  WHERE bhisham_id=$4 AND item_name=$5 AND mc_no=$6 AND cube_number=$7`
		queryParams = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.ID, UserID}
		ParamsMapping = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.BhishamID, obj.KitName, obj.MCNo, obj.CubeNumber, UserID}

	default:
		return helper.CreateDynamicResponse("Invalid Update Type", false, nil, 400, nil), fmt.Errorf("invalid update type: %d", obj.UpdateType)
	}

	// Start a transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return helper.CreateDynamicResponse("Transaction start error: "+err.Error(), false, nil, 500, nil), err
	}

	// Execute bhisham_data update
	if err := executeUpdateQuery(tx, updateBhishamQuery, queryParams); err != nil {
		tx.Rollback()
		return helper.CreateDynamicResponse("Error updating bhisham_data: "+err.Error(), false, nil, 400, nil), err
	}

	// Execute bhisham_mapping update (use correct params for type 3)
	if obj.UpdateType == 3 {
		if err := executeUpdateQuery(tx, updateBhishamMappingQuery, ParamsMapping); err != nil {
			tx.Rollback()
			return helper.CreateDynamicResponse("Error updating bhisham_mapping: "+err.Error(), false, nil, 400, nil), err
		}
	} else {
		if err := executeUpdateQuery(tx, updateBhishamMappingQuery, queryParams); err != nil {
			tx.Rollback()
			return helper.CreateDynamicResponse("Error updating bhisham_mapping: "+err.Error(), false, nil, 400, nil), err
		}
	}

	Updatequery := `UPDATE public.bhisham_mapping 
	SET kit_expiry = temp.min_expiry
	FROM (
		SELECT cc_no, kitname, MIN(CAST(exp AS DATE)) AS min_expiry 
		FROM public.bhisham_mapping 
		WHERE LENGTH(exp) > 5 AND bhisham_id = $1
		GROUP BY cc_no, kitname
	) temp  
	WHERE public.bhisham_mapping.kitname = temp.kitname 
	AND public.bhisham_mapping.cc_no = temp.cc_no
	AND public.bhisham_mapping.bhisham_id = $1`

	_, err = tx.Exec(Updatequery, obj.BhishamID)
	if err != nil {
		return helper.CreateDynamicResponse("Error updating expiry >"+err.Error(), false, nil, 500, nil), err
	}

	Updatequery = `UPDATE public.bhisham_data
	SET kit_expiry = temp.min_expiry
	FROM (
		SELECT kit_epc, MIN(CAST(exp AS DATE)) AS min_expiry 
		FROM public.bhisham_data 
		WHERE LENGTH(exp) > 5 AND bhisham_id = $1
		GROUP BY kit_epc
	) temp  
	WHERE public.bhisham_data.kit_epc = temp.kit_epc 
	AND public.bhisham_data.bhisham_id = $1`

	_, err = tx.Exec(Updatequery, obj.BhishamID)
	if err != nil {
		return helper.CreateDynamicResponse("Error updating expiry >"+err.Error(), false, nil, 500, nil), err
	}

	// Log the update in `update_bhisham_data` table
	logQuery := `INSERT INTO public.update_bhisham_data 
				 (bhisham_id, mc_no, cube_number, kit_name, batch_code, mfd, exp, update_type_id, created_by) 
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	logParams := []interface{}{obj.BhishamID, obj.MCNo, obj.CubeNumber, obj.KitName, obj.BatchCode, obj.MFD, obj.EXP, obj.UpdateType, UserID}

	if _, err := tx.Exec(logQuery, logParams...); err != nil {
		tx.Rollback()
		return helper.CreateDynamicResponse("Error logging update: "+err.Error(), false, nil, 400, nil), err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return helper.CreateDynamicResponse("Transaction commit error: "+err.Error(), false, nil, 500, nil), err
	}

	return helper.CreateDynamicResponse("Bhisham Updated Successfully", true, nil, 200, nil), nil
}

func (r *BhishamRepository) oldUpdateBhishamMapping(obj models.UpdateBhishamData, UserID string) (map[string]interface{}, error) {
	var updateBhishamQuery string
	var queryParams []interface{}

	// Determine the update query based on UpdateType
	switch obj.UpdateType {
	case 1:
		updateBhishamQuery = `UPDATE public.bhisham_mapping 
							  SET mfd=$1, exp=$2, batch_no_sr_no=$3 
							  WHERE bhisham_id=$4 AND sku_name=$5 AND mc_no=$6`
		queryParams = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.BhishamID, obj.KitName, obj.MCNo}

	case 2:
		updateBhishamQuery = `UPDATE public.bhisham_mapping 
							  SET mfd=$1, exp=$2, batch_no_sr_no=$3 
							  WHERE bhisham_id=$4 AND sku_name=$5 AND mc_no=$6 AND cube_number=$7`
		queryParams = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.BhishamID, obj.KitName, obj.MCNo, obj.CubeNumber}

	case 3:
		updateBhishamQuery = `UPDATE public.bhisham_mapping 
							  SET mfd=$1, exp=$2, batch_no_sr_no=$3 
							  WHERE id=$4`
		queryParams = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.ID}

	default:
		return helper.CreateDynamicResponse("Invalid Update Type", false, nil, 400, nil), fmt.Errorf("invalid update type: %d", obj.UpdateType)
	}

	// Execute Update Query (Exec instead of QueryRow)
	res, err := r.DB.Exec(updateBhishamQuery, queryParams...)
	if err != nil {
		return helper.CreateDynamicResponse("Error updating Bhisham data: "+err.Error(), false, nil, 400, nil), err
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching update count: "+err.Error(), false, nil, 400, nil), err
	}
	if rowsAffected == 0 {
		return helper.CreateDynamicResponse("No records updated", false, nil, 200, nil), nil
	}

	// Log update in `update_bhisham_data` table
	logQuery := `INSERT INTO public.update_bhisham_data 
				 (bhisham_id, mc_no, cube_number, kit_name, batch_code, mfd, exp, update_type_id, created_by) 
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	logParams := []interface{}{obj.BhishamID, obj.MCNo, obj.CubeNumber, obj.KitName, obj.BatchCode, obj.MFD, obj.EXP, obj.UpdateType, UserID}

	_, err = r.DB.Exec(logQuery, logParams...)
	if err != nil {
		return helper.CreateDynamicResponse("Error logging update: "+err.Error(), false, nil, 400, nil), err
	}

	return helper.CreateDynamicResponse(fmt.Sprintf("Bhisham Updated Successfully (%d rows affected)", rowsAffected), true, nil, 200, nil), nil
}

func (r *BhishamRepository) UpdateBhishamMapping(obj models.UpdateBhishamData, UserID string) (map[string]interface{}, error) {
	var updateBhishamQuery string
	var queryParams []interface{}

	// Start transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return helper.CreateDynamicResponse("Failed to start transaction: "+err.Error(), false, nil, 500, nil), err
	}

	// Determine the update queries based on UpdateType
	switch obj.UpdateType {
	case 1:
		updateBhishamQuery = `UPDATE public.bhisham_mapping 
					  SET mfd=$1, exp=$2, batch_no_sr_no=$3, is_update=1, update_time=NOW(), updated_by=$7
					  WHERE bhisham_id=$4 AND item_name=$5 AND mc_no=$6`
		queryParams = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.BhishamID, obj.KitName, obj.MCNo, UserID}

	case 2:
		updateBhishamQuery = `UPDATE public.bhisham_mapping 
					  SET mfd=$1, exp=$2, batch_no_sr_no=$3, is_update=1, update_time=NOW(), updated_by=$8
					  WHERE bhisham_id=$4 AND item_name=$5 AND mc_no=$6 AND cube_number=$7`
		queryParams = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.BhishamID, obj.KitName, obj.MCNo, obj.CubeNumber, UserID}

	case 3:
		updateBhishamQuery = `UPDATE public.bhisham_mapping 
					  SET mfd=$1, exp=$2, batch_no_sr_no=$3, is_update=1, update_time=NOW(), updated_by=$5
					  WHERE id=$4`
		queryParams = []interface{}{obj.MFD, obj.EXP, obj.BatchCode, obj.ID, UserID}

	default:
		tx.Rollback() // Rollback transaction on invalid update type
		return helper.CreateDynamicResponse("Invalid Update Type", false, nil, 400, nil), fmt.Errorf("invalid update type: %d", obj.UpdateType)
	}

	// Execute bhisham_mapping update within the transaction
	res, err := tx.Exec(updateBhishamQuery, queryParams...)
	if err != nil {
		tx.Rollback()
		return helper.CreateDynamicResponse("Error updating bhisham_mapping: "+err.Error(), false, nil, 400, nil), err
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return helper.CreateDynamicResponse("Error fetching update count: "+err.Error(), false, nil, 400, nil), err
	}
	if rowsAffected == 0 {
		tx.Rollback()
		return helper.CreateDynamicResponse("No records updated", false, nil, 200, nil), nil
	}

	// Update expiry in bhisham_mapping
	Updatequery := `UPDATE public.bhisham_mapping 
	SET kit_expiry = temp.min_expiry
	FROM (
		SELECT cc_no, kitname, MIN(CAST(exp AS DATE)) AS min_expiry 
		FROM public.bhisham_mapping 
		WHERE LENGTH(exp) > 5 AND bhisham_id = $1
		GROUP BY cc_no, kitname
	) temp  
	WHERE public.bhisham_mapping.kitname = temp.kitname 
	AND public.bhisham_mapping.cc_no = temp.cc_no
	AND public.bhisham_mapping.bhisham_id = $1`

	_, err = tx.Exec(Updatequery, obj.BhishamID)
	if err != nil {
		tx.Rollback()
		return helper.CreateDynamicResponse("Error updating expiry >"+err.Error(), false, nil, 500, nil), err
	}

	// Log the update in `update_bhisham_data` table
	logQuery := `INSERT INTO public.update_bhisham_data 
				 (bhisham_id, mc_no, cube_number, kit_name, batch_code, mfd, exp, update_type_id, created_by) 
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	logParams := []interface{}{obj.BhishamID, obj.MCNo, obj.CubeNumber, obj.KitName, obj.BatchCode, obj.MFD, obj.EXP, obj.UpdateType, UserID}

	_, err = tx.Exec(logQuery, logParams...)
	if err != nil {
		tx.Rollback()
		return helper.CreateDynamicResponse("Error logging update: "+err.Error(), false, nil, 400, nil), err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return helper.CreateDynamicResponse("Transaction commit failed: "+err.Error(), false, nil, 500, nil), err
	}

	return helper.CreateDynamicResponse(fmt.Sprintf("Bhisham Updated Successfully (%d rows affected)", rowsAffected), true, nil, 200, nil), nil
}

func FindMinExpiry(items []models.BhishamMappingData) string {
	var minExpiry string
	var minTime time.Time
	hasValidExpiry := false

	for _, item := range items {
		exp := strings.TrimSpace(item.EXP) // Trim spaces

		// Ignore empty or "NA" values
		if exp == "" || strings.ToUpper(exp) == "NA" {
			continue
		}

		// Parse expiry date
		expTime, err := time.Parse("2006-01-02", exp) // Adjust format if needed
		if err != nil {
			continue
		}

		// Initialize minTime or update if new expiry is earlier
		if !hasValidExpiry || expTime.Before(minTime) {
			minTime = expTime
			minExpiry = exp
			hasValidExpiry = true
		}
	}

	// If no valid expiry found, return "NA"
	if !hasValidExpiry {
		return "NA"
	}
	return minExpiry
}

func GenerateKitEPC(bhishamID, mcNo, boxNumber, packNo int) string {
	// Format components with strict padding:
	// - bhishamID: 8 digits (padded with leading zeros)
	// - mcNo: used as-is (must be 3 digits for consistent length)
	// - boxNumber: 2 digits (padded with leading zeros)
	// - packNo: 2 digits (padded with leading zeros) as last two characters

	// Validate mcNo is exactly 3 digits
	mcNoStr := fmt.Sprintf("%03d", mcNo)
	if len(mcNoStr) != 3 {
		mcNoStr = fmt.Sprintf("%03d", mcNo%1000) // Force to 3 digits
	}

	// Pad all other components
	bhishamIDPadded := fmt.Sprintf("%08d", bhishamID)
	boxNumberPadded := fmt.Sprintf("%02d", boxNumber)
	packNoPadded := fmt.Sprintf("%02d", packNo)

	// Construct fixed-format EPC (24 characters total)
	// Format: "BA0" + mcNo (3) + "CA" + mcNo (3) + "0" + boxNumber (2) + bhishamID (8) + "00" + packNo (2)
	epc := fmt.Sprintf("BA0%sCA%s0%s%s00%s",
		mcNoStr,
		mcNoStr,
		boxNumberPadded,
		bhishamIDPadded,
		packNoPadded)

	// Final validation (should always be 24 chars with this construction)
	if len(epc) != 24 {
		// Fallback to ensure length if assumptions fail
		if len(epc) > 24 {
			epc = epc[:24]
		} else {
			epc += strings.Repeat("0", 24-len(epc))
		}
	}

	return epc
}

func GetProductsByKit(mappings []models.BhishamMappingData, kitName string, ccNo string) []models.BhishamMappingData {
	var result []models.BhishamMappingData
	for _, item := range mappings {
		if strings.EqualFold(item.KitName, kitName) && strings.EqualFold(item.CCNo, ccNo) {
			result = append(result, item) // Append full struct
		}
	}
	return result
}

func GetBox(mappings []models.BhishamMappingData) []string {
	uniqueCCNo := make(map[string]bool) // To store unique values
	var result []string

	for _, item := range mappings {
		if _, exists := uniqueCCNo[item.CCNo]; !exists {
			uniqueCCNo[item.CCNo] = true
			result = append(result, item.CCNo)
		}
	}
	return result
}

// Helper function to execute an update query and check affected rows
func executeUpdateQuery(tx *sql.Tx, query string, params []interface{}) error {
	res, err := tx.Exec(query, params...)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no records updated")
	}
	return nil
}

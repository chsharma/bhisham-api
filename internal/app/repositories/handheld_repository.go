package repositories

import (
	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/models"
	"database/sql"
)

type HandheldRepository struct {
	DB *sql.DB
}

func (r *HandheldRepository) GetBhishamID(SerialNo string) (map[string]interface{}, error) {
	var BhishamID int

	// Use parameterized query to avoid SQL injection
	query := `SELECT id FROM public.bhisham WHERE serial_no = $1 LIMIT 1`
	rows, err := r.DB.Query(query, SerialNo)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching Bhisham ID", false, nil, 500, nil), err
	}
	defer rows.Close() // Ensure rows are closed after use

	// Check if there's a result
	if rows.Next() {
		if err := rows.Scan(&BhishamID); err != nil {
			return helper.CreateDynamicResponse("Error scanning Bhisham ID", false, nil, 500, nil), err
		}
	} else {
		return helper.CreateDynamicResponse("No record found", false, nil, 404, nil), nil
	}

	return helper.CreateDynamicResponse("Bhisham ID fetched successfully", true, BhishamID, 200, nil), nil
}

func (r *HandheldRepository) GetAllBhishamData(BhishamID int) (map[string]interface{}, error) {
	var bhishamID int

	// Fetch bhisham_id
	query := `SELECT id FROM public.bhisham WHERE id = $1 AND is_complete = 1 LIMIT 1`
	err := r.DB.QueryRow(query, BhishamID).Scan(&bhishamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.CreateDynamicResponse("Bhisham data is not completed.", false, nil, 404, nil), nil
		}
		return helper.CreateDynamicResponse("Error fetching Bhisham ID", false, nil, 500, nil), err
	}

	// Query for Bhisham data (Streaming for Large Data)
	query = `SELECT id, bhisham_id, mc_no, mc_name, mc_epc, cc_no, cc_name, cc_epc, kitcode, kit_no, kit_epc, 
	                 kit_batch_no, kit_expiry, kit_qty, sku_code, sku_name, batch_no_sr_no, mfd, exp, 
	                 manufactured_by, sku_qty, cube_number, kitname, no_of_kit ,is_update,update_time,updated_by
	          FROM public.bhisham_data 
	          WHERE bhisham_id = $1 
	          ORDER BY bhisham_id, mc_no, cube_number`

	rows, err := r.DB.Query(query, bhishamID)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching data", false, nil, 500, nil), err
	}
	defer rows.Close()

	// Use a buffered streaming approach
	bhishamDataList := make([]models.BhishamData, 0, 100) // Allocate only 100 at a time

	for rows.Next() {
		var data models.BhishamData
		err := rows.Scan(
			&data.ID, &data.BhishamID, &data.MCNo, &data.MCName, &data.MCEPC, &data.CCNo, &data.CCName, &data.CCEPC,
			&data.KitCode, &data.KitNo, &data.KitEPC, &data.KitBatchNo, &data.KitExpiry, &data.KitQty,
			&data.SKUCode, &data.SKUName, &data.BatchNoSrNo, &data.Mfd, &data.Exp,
			&data.ManufacturedBy, &data.SKUQty, &data.CubeNumber, &data.KitName, &data.NoOfKit, &data.IsUpdate, &data.UpdateTime, &data.UpdateBy,
		)
		if err != nil {
			return helper.CreateDynamicResponse("Error scanning data", false, nil, 500, nil), err
		}
		bhishamDataList = append(bhishamDataList, data)

		// If memory usage is a concern, consider processing data here instead of appending
	}

	// Return empty array if no records found
	if len(bhishamDataList) == 0 {
		return helper.CreateDynamicResponse("No records found", true, []models.BhishamData{}, 200, nil), nil
	}

	// Update `bhisham` table
	var newID int
	UpdateQuery := `UPDATE bhisham 
	                SET is_hh_synch = 1, hh_synch_time = NOW(), hh_synch_count = hh_synch_count + 1 
	                WHERE id = $1 RETURNING id`

	err = r.DB.QueryRow(UpdateQuery, BhishamID).Scan(&newID)
	if err != nil {
		return helper.CreateDynamicResponse("Error updating bhisham", false, nil, 500, nil), err
	}

	// Response count
	Count := map[string]interface{}{
		"total_rows": len(bhishamDataList),
	}
	return helper.CreateDynamicResponse("Data fetched successfully", true, bhishamDataList, 200, Count), nil
}

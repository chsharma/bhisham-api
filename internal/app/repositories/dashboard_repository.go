package repositories

import (
	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/models"
	"database/sql"
)

type DashboardRepository struct {
	DB *sql.DB
}

func (r *DashboardRepository) GetDashboardStats() (map[string]interface{}, error) {

	var totalBhisham int
	err := r.DB.QueryRow("SELECT COUNT(1) FROM public.bhisham;").Scan(&totalBhisham)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching total count", false, nil, 500, nil), err
	}

	var totalIMotherCube int
	err = r.DB.QueryRow("SELECT COUNT(*) FROM(SELECT DISTINCT bhisham_id,mc_no  FROM public.bhisham_mapping );").Scan(&totalIMotherCube)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching total count", false, nil, 500, nil), err
	}

	var totalICube int
	err = r.DB.QueryRow("SELECT COUNT(*) FROM(SELECT DISTINCT bhisham_id,mc_no,cc_no  FROM public.bhisham_mapping where  is_cube=1);").Scan(&totalICube)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching total count", false, nil, 500, nil), err
	}

	var totalIKits int
	err = r.DB.QueryRow("SELECT COUNT(*) FROM(SELECT DISTINCT bhisham_id,mc_no,cc_no,kitcode  FROM public.bhisham_mapping where  is_cube=1);").Scan(&totalIKits)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching total count", false, nil, 500, nil), err
	}

	response := map[string]interface{}{
		"bhisham": totalBhisham,
		"mc":      totalIMotherCube,
		"cc":      totalICube,
		"kits":    totalIKits,
	}

	return helper.CreateDynamicResponse("All users fetched successfully", true, response, 200, nil), nil

}

func (r *DashboardRepository) GetBhisham() (map[string]interface{}, error) {

	rows, err := r.DB.Query(`SELECT id, serial_no, bhisham_name, created_by, created_at, is_complete, complete_by, complete_time, is_hh_synch,hh_serial,hh_synch_time,hh_synch_count FROM public.bhisham order by id desc`)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching bhishams", false, nil, 500, nil), err
	}
	defer rows.Close()

	var bhs []models.Bhisham
	for rows.Next() {
		var bh models.Bhisham
		if err := rows.Scan(
			&bh.ID, &bh.SerialNo, &bh.BhishamName, &bh.CreatedBy, &bh.CreatedAt, &bh.IsComplete,
			&bh.CompleteBy, &bh.CompleteTime, &bh.IsHHSynch, &bh.HHSerial, &bh.HHSynchTime, &bh.HHSynchCount,
		); err != nil {
			return helper.CreateDynamicResponse("Error scanning bhisham data", false, nil, 500, nil), err
		}
		bhs = append(bhs, bh)
	}
	return helper.CreateDynamicResponse("All bhishams fetched successfully", true, bhs, 200, nil), nil

}

func (r *DashboardRepository) GetChildCube(BhishamID, MotherCubeID int) (map[string]interface{}, error) {
	rows, err := r.DB.Query(`SELECT bhisham_id, mc_no, mc_name, cube_number, cc_no, cc_name, count(bhisham_id) as total, SUM(is_update) as total_update
                            FROM public.bhisham_mapping 
                            WHERE is_cube = 1 AND bhisham_id = $1 AND mc_no = $2  group by bhisham_id, mc_no, mc_name, cube_number, cc_no, cc_name
                            ORDER BY cube_number`, BhishamID, MotherCubeID)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching cubes", false, nil, 500, nil), err
	}
	defer rows.Close()

	var bhs []models.BhishamMapping
	for rows.Next() {
		var bh models.BhishamMapping
		if err := rows.Scan(
			&bh.BhishamID, &bh.MCNo, &bh.MCName, &bh.CubeNumber, &bh.CCNo, &bh.CCName, &bh.Total, &bh.TotalUpdate,
		); err != nil {
			return helper.CreateDynamicResponse("Error scanning cubes data", false, nil, 500, nil), err
		}
		bhs = append(bhs, bh)
	}

	// Check if no records were found
	if len(bhs) == 0 {
		return helper.CreateDynamicResponse("No cubes found", true, []models.BhishamMapping{}, 200, nil), nil
	}

	return helper.CreateDynamicResponse("All cubes fetched successfully", true, bhs, 200, nil), nil
}

func (r *DashboardRepository) GetChildKits(BhishamID, MotherCubeID, CCNo int) (map[string]interface{}, error) {
	rows, err := r.DB.Query(`SELECT bhisham_id, mc_no, cc_no, kitcode, kitname, kit_slug, no_of_kit ,COALESCE(kit_expiry, 'NA') AS kit_expiry, count(bhisham_id) as total, SUM(is_update) as total_update
                            FROM public.bhisham_data
                            WHERE is_cube = 1 AND bhisham_id = $1 AND mc_no = $2 AND cube_number = $3 group by bhisham_id, mc_no, cc_no, kitcode, kitname, kit_slug, no_of_kit, kit_expiry`,
		BhishamID, MotherCubeID, CCNo)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching kits "+err.Error(), false, nil, 500, nil), err
	}
	defer rows.Close()

	var kits []models.BhishamKit
	for rows.Next() {
		var kit models.BhishamKit
		if err := rows.Scan(
			&kit.BhishamID, &kit.MCNo, &kit.CCNo, &kit.KitCode, &kit.KitName, &kit.KitSlug, &kit.NoOfKit, &kit.KitExpiry, &kit.Total, &kit.TotalUpdate,
		); err != nil {
			return helper.CreateDynamicResponse("Error scanning kits data "+err.Error(), false, nil, 500, nil), err
		}
		kits = append(kits, kit)
	}

	// Handle case where no records are found
	if len(kits) == 0 {
		return helper.CreateDynamicResponse("No kits found", true, []models.BhishamKit{}, 200, nil), nil
	}

	return helper.CreateDynamicResponse("All kits fetched successfully", true, kits, 200, nil), nil
}

// remove kit no
func (r *DashboardRepository) GetKitItems(BhishamID, MotherCubeID, CCNo int, KitName string) (map[string]interface{}, error) {
	query := `SELECT id, mc_no, cube_number, kitname, kit_no, batch_no_sr_no, 
                     sku_name, mfd, exp, manufactured_by, sku_qty, 
                     is_update, update_time, updated_by, sku_code, kitcode, kit_slug, sku_slug
              FROM public.bhisham_data 
              WHERE bhisham_id = $1 
                AND mc_no = $2 
                AND cube_number = $3 
                AND kit_slug = $4 order by mc_no, cube_number,kit_no`

	rows, err := r.DB.Query(query, BhishamID, MotherCubeID, CCNo, KitName)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching items", false, nil, 500, nil), err
	}
	defer rows.Close()

	var kits []models.KitItems
	for rows.Next() {
		var kit models.KitItems
		err := rows.Scan(
			&kit.ID,
			&kit.MCNo,
			&kit.CubeNumber,
			&kit.KitName,
			&kit.KitNo,
			&kit.BatchNoSrNo,
			&kit.SKUName,
			&kit.Mfd,
			&kit.Exp,
			&kit.ManufacturedBy,
			&kit.SKUQty,
			&kit.IsUpdate,
			&kit.UpdateTime,
			&kit.UpdateBy,
			&kit.SKUCode,
			&kit.KitCode,
			&kit.KitSlug,
			&kit.SKUSlug,
		)
		if err != nil {
			return helper.CreateDynamicResponse("Error scanning items data "+err.Error(), false, nil, 500, nil), err
		}
		kits = append(kits, kit)
	}

	// Check for errors that may have occurred during iteration
	if err = rows.Err(); err != nil {
		return helper.CreateDynamicResponse("Error after row iteration", false, nil, 500, nil), err
	}

	// Handle case where no records are found
	if len(kits) == 0 {
		return helper.CreateDynamicResponse("No items found", true, []models.KitItems{}, 200, nil), nil
	}

	return helper.CreateDynamicResponse("Items fetched successfully", true, kits, 200, nil), nil
}

func (r *DashboardRepository) GetMappingKitItems(BhishamID, MotherCubeID, CCNo int, KitName string) (map[string]interface{}, error) {
	query := `SELECT id, mc_no, cube_number, kitname, 1 as kit_no, batch_no_sr_no, item_name as sku_name, mfd, exp, manufactured_by, no_of_item as sku_qty,is_update,update_time,updated_by,sku_code, kitcode, kit_slug, sku_slug FROM public.bhisham_mapping 
			  WHERE bhisham_id=$1 AND mc_no=$2 AND cube_number=$3 and kit_slug=$4`

	rows, err := r.DB.Query(query, BhishamID, MotherCubeID, CCNo, KitName)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching items "+err.Error(), false, nil, 500, nil), err
	}
	defer rows.Close()

	var kits []models.KitItems
	for rows.Next() {
		var kit models.KitItems
		if err := rows.Scan(
			&kit.ID, &kit.MCNo, &kit.CubeNumber, &kit.KitName, &kit.KitNo,
			&kit.BatchNoSrNo, &kit.SKUName, &kit.Mfd, &kit.Exp, &kit.ManufacturedBy, &kit.SKUQty, &kit.IsUpdate, &kit.UpdateTime, &kit.UpdateBy, &kit.SKUCode, &kit.KitCode,
			&kit.KitSlug,
			&kit.SKUSlug,
		); err != nil {
			return helper.CreateDynamicResponse("Error scanning items data", false, nil, 500, nil), err
		}
		kits = append(kits, kit)
	}

	// Handle case where no records are found
	if len(kits) == 0 {
		return helper.CreateDynamicResponse("No items found", true, []models.KitItems{}, 200, nil), nil
	}

	return helper.CreateDynamicResponse("All items fetched successfully", true, kits, 200, nil), nil
}

func (r *DashboardRepository) GetUpdateType() (map[string]interface{}, error) {

	rows, err := r.DB.Query(`select update_typeid,name from public.data_update_type order by order_by`)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching data_update_type", false, nil, 500, nil), err
	}
	defer rows.Close()

	var bhs []models.UpdateType
	for rows.Next() {
		var bh models.UpdateType
		if err := rows.Scan(
			&bh.UpdateTypeID, &bh.Name,
		); err != nil {
			return helper.CreateDynamicResponse("Error scanning data_update_type ? "+err.Error(), false, nil, 500, nil), err
		}
		bhs = append(bhs, bh)
	}
	return helper.CreateDynamicResponse("All bhidata_update_type fetched successfully", true, bhs, 200, nil), nil

}

func (r *DashboardRepository) GetAllMappingBhishamData(BhishamID int) (map[string]interface{}, error) {
	// Query for Bhisham data (Streaming for Large Data)
	query := `SELECT mc_name, cc_no, cc_name, kitname, no_of_kit, sku_code, item_name as sku_name, batch_no_sr_no, mfd, exp, 
	                 manufactured_by, no_of_item 
	          FROM public.bhisham_mapping
	          WHERE bhisham_id = $1
	          ORDER BY bhisham_id, mc_no, cube_number`

	rows, err := r.DB.Query(query, BhishamID)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching data", false, nil, 500, nil), err
	}
	defer rows.Close()

	// Use a buffered streaming approach
	bhishamDataList := make([]models.ReportBhishamMapping, 0, 100) // Allocate only 100 at a time

	for rows.Next() {
		var data models.ReportBhishamMapping
		err := rows.Scan(
			&data.McName, &data.CcNo, &data.CcName, &data.KitName, &data.NoOfKit,
			&data.SkuCode, &data.SkuName, &data.BatchNoSrNo, &data.Mfd, &data.Exp,
			&data.ManufacturedBy, &data.NoOfItem,
		)
		if err != nil {
			return helper.CreateDynamicResponse("Error scanning data", false, nil, 500, nil), err
		}
		bhishamDataList = append(bhishamDataList, data)
	}

	// Return empty array if no records found
	if len(bhishamDataList) == 0 {
		return helper.CreateDynamicResponse("No records found", true, []models.ReportBhishamMapping{}, 200, nil), nil
	}

	// Response count
	count := map[string]interface{}{
		"total_rows": len(bhishamDataList),
	}
	return helper.CreateDynamicResponse("Data fetched successfully", true, bhishamDataList, 200, count), nil
}

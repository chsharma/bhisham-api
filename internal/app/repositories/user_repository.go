package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/models"
	"bhisham-api/internal/app/utils"
	"bhisham-api/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

// CreateUser - Register new user with hashed password
func (r *UserRepository) CreateUser(user models.User) (map[string]interface{}, error) {
	// Hash Password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return helper.CreateDynamicResponse("Error hashing password", false, nil, 500, nil), err
	}

	var user_id = utils.GenerateId()
	user.UserID = user_id

	query := `INSERT INTO user_login (user_id, name, login_id, pwd, role_id, created_at) 
              VALUES ($1, $2, $3, $4, $5, NOW())`

	_, err = r.DB.Exec(query, user.UserID, user.Name, user.LoginID, string(hashedPwd), user.RoleID)
	if err != nil {
		return helper.CreateDynamicResponse("Error creating user", false, nil, 500, nil), err
	}

	return helper.CreateDynamicResponse("User Created Successfully", true, user, 200, nil), nil
}

// UpdateUser updates user details while checking for duplicate login_id
func (r *UserRepository) UpdateUser(user models.User) (map[string]interface{}, error) {
	// Check if the login_id already exists for another user
	var existingUserID string
	checkQuery := `SELECT user_id FROM user_login WHERE login_id=$1 AND user_id<>$2`
	err := r.DB.QueryRow(checkQuery, user.LoginID, user.UserID).Scan(&existingUserID)
	if err == nil {
		return helper.CreateDynamicResponse("Login ID already in use", false, nil, 400, nil), fmt.Errorf("duplicate login_id")
	} else if err != sql.ErrNoRows {
		return helper.CreateDynamicResponse("Error checking login ID", false, nil, 500, nil), err
	}

	// Proceed with updating the user
	query := `UPDATE user_login SET name=$1, login_id=$2, role_id=$3 WHERE user_id=$4`
	res, err := r.DB.Exec(query, user.Name, user.LoginID, user.RoleID, user.UserID)
	if err != nil {
		return helper.CreateDynamicResponse("Error updating user", false, nil, 500, nil), err
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching update count", false, nil, 500, nil), err
	}
	if rowsAffected == 0 {
		return helper.CreateDynamicResponse("No user found to update", false, nil, 404, nil), nil
	}

	return helper.CreateDynamicResponse("User updated successfully", true, nil, 200, nil), nil
}

// GetUsers retrieves all users from the user_login table
func (r *UserRepository) GetUsers() (map[string]interface{}, error) {
	query := `SELECT user_id, name, login_id, role_id, created_at, active FROM user_login`
	rows, err := r.DB.Query(query)
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching users: "+err.Error(), false, nil, 500, nil), err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.LoginID, &user.RoleID, &user.CreatedAt, &user.Active); err != nil {
			return helper.CreateDynamicResponse("Error scanning user: "+err.Error(), false, nil, 500, nil), err
		}

		users = append(users, user)
	}

	return helper.CreateDynamicResponse("Users fetched successfully", true, users, 200, nil), nil
}

// UpdatePassword updates the password for a given user
func (r *UserRepository) UpdatePassword(userID, newPassword string) (map[string]interface{}, error) {
	// Hash the new password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return helper.CreateDynamicResponse("Error hashing new password", false, nil, 500, nil), err
	}

	query := `UPDATE user_login SET pwd=$1 WHERE user_id=$2`
	res, err := r.DB.Exec(query, string(hashedPwd), userID)
	if err != nil {
		return helper.CreateDynamicResponse("Error updating password", false, nil, 500, nil), err
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return helper.CreateDynamicResponse("Error fetching update count: "+err.Error(), false, nil, 500, nil), err
	}
	if rowsAffected == 0 {
		return helper.CreateDynamicResponse("No user found to update password", false, nil, 404, nil), nil
	}

	return helper.CreateDynamicResponse("Password updated successfully", true, nil, 200, nil), nil
}

// AuthenticateUser - Verify login credentials
func (r *UserRepository) AuthenticateUser(loginID, password string) (map[string]interface{}, error) {
	var user models.User
	var hashedPwd string

	// Corrected SQL Query
	query := `SELECT user_id, name, login_id, pwd, active, role_id, created_at FROM user_login WHERE login_id = $1`
	err := r.DB.QueryRow(query, loginID).Scan(&user.UserID, &user.Name, &user.LoginID, &hashedPwd, &user.Active, &user.RoleID, &user.CreatedAt)
	if err != nil {
		return helper.CreateDynamicResponse("Invalid login credentials", false, nil, 500, nil), err
	}

	// Check if account is active
	if !user.Active {
		return helper.CreateDynamicResponse("Account is inactive", false, nil, 500, nil), errors.New("account is inactive")
	}

	// Verify Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
	if err != nil {
		return helper.CreateDynamicResponse("Invalid login credentials", false, nil, 500, nil), err
	}

	// Generate Token
	newToken, _ := middleware.GenerateNewToken(user.LoginID) // Use user.UserID instead of separate user_id
	token := map[string]interface{}{
		"token": newToken,
	}

	// Return User Data & Token
	return helper.CreateDynamicResponse("Login successful", true, user, 200, token), nil
}

package repositories

import (
	"database/sql"
	"errors"

	"bhisham-api/internal/app/helper"
	"bhisham-api/internal/app/middleware"
	"bhisham-api/internal/app/models"
	"bhisham-api/internal/app/utils"

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
	newToken, _ := middleware.GenerateNewToken(user.UserID) // Use user.UserID instead of separate user_id
	token := map[string]interface{}{
		"token": newToken,
	}

	// Return User Data & Token
	return helper.CreateDynamicResponse("Login successful", true, user, 200, token), nil
}

package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID              uint      `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	PrivacySettings string    `json:"privacy_settings"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UserService handles user-related operations
type UserService struct {
	DB *sql.DB
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

// Create creates a new user
func (s *UserService) Create(username, email, password string, privacySettings map[string]interface{}) (*User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Convert privacy settings to JSON
	privacyJSON, err := json.Marshal(privacySettings)
	if err != nil {
		return nil, err
	}

	// Insert the user into the database
	query := `
		INSERT INTO users (username, email, password, privacy_settings)
		VALUES (?, ?, ?, ?)
		RETURNING id, username, email, privacy_settings, created_at, updated_at
	`

	var user User
	err = s.DB.QueryRow(
		query,
		username,
		email,
		string(hashedPassword),
		string(privacyJSON),
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PrivacySettings,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id uint) (*User, error) {
	query := `
		SELECT id, username, email, password, privacy_settings, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user User
	err := s.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.PrivacySettings,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, username, email, password, privacy_settings, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var user User
	err := s.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.PrivacySettings,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// UpdatePrivacySettings updates a user's privacy settings
func (s *UserService) UpdatePrivacySettings(userID uint, settings map[string]interface{}) error {
	// Convert settings to JSON
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	// Update the user's privacy settings
	query := `
		UPDATE users
		SET privacy_settings = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err = s.DB.Exec(query, string(settingsJSON), userID)
	return err
}

// Authenticate verifies a user's credentials
func (s *UserService) Authenticate(email, password string) (*User, error) {
	// Get the user by email
	user, err := s.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

package repository

import (
	"go-ticket/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	*Repository[models.User]
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		Repository: NewRepository[models.User](db, "users"),
		db:         db,
	}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT * FROM users 
		WHERE email = $1 
		AND deleted_at IS NULL
	`
	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	query := `
		SELECT * FROM users 
		WHERE phone = $1 
		AND deleted_at IS NULL
	`
	err := r.db.Get(&user, query, phone)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (
			id, fullname, email, phone, password,
			created_at, updated_at
		) VALUES (
			:id, :fullname, :email, :phone, :password,
			:created_at, :updated_at
		)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         user.ID,
		"fullname":   user.Fullname,
		"email":      user.Email,
		"phone":      user.Phone,
		"password":   "-",
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
	return err
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users SET
			fullname = :fullname,
			email = :email,
			phone = :phone,
			password = :password,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         user.ID,
		"fullname":   user.Fullname,
		"email":      user.Email,
		"phone":      user.Phone,
		"password":   user.Password,
		"updated_at": user.UpdatedAt,
	})
	return err
}

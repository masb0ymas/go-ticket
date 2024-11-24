package repository

import (
	"go-ticket/models"

	"github.com/jmoiron/sqlx"
)

type LocationRepository struct {
	*Repository[models.Location]
}

func NewLocationRepository(db *sqlx.DB) *LocationRepository {
	return &LocationRepository{
		Repository: NewRepository[models.Location](db, "locations"),
	}
}

// Custom methods for LocationRepository
func (r *LocationRepository) FindByCity(city string) ([]models.Location, error) {
	query := `
		SELECT * FROM locations 
		WHERE LOWER(city) LIKE LOWER($1) 
		AND deleted_at IS NULL
	`

	var locations []models.Location
	err := r.db.Select(&locations, query, "%"+city+"%")
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *LocationRepository) FindByCountry(country string) ([]models.Location, error) {
	query := `
		SELECT * FROM locations 
		WHERE LOWER(country) = LOWER($1) 
		AND deleted_at IS NULL
	`

	var locations []models.Location
	err := r.db.Select(&locations, query, country)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *LocationRepository) Create(location *models.Location) error {
	query := `
		INSERT INTO locations (
			id, name, address, city, state, country, postal_code,
			created_at, updated_at
		) VALUES (
			:id, :name, :address, :city, :state, :country, :postal_code,
			:created_at, :updated_at
		)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":          location.ID,
		"name":        location.Name,
		"address":     location.Address,
		"city":        location.City,
		"state":       location.State,
		"country":     location.Country,
		"postal_code": location.PostalCode,
		"created_at":  location.CreatedAt,
		"updated_at":  location.UpdatedAt,
	})
	return err
}

func (r *LocationRepository) Update(location *models.Location) error {
	query := `
		UPDATE locations SET
			name = :name,
			address = :address,
			city = :city,
			state = :state,
			country = :country,
			postal_code = :postal_code,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":          location.ID,
		"name":        location.Name,
		"address":     location.Address,
		"city":        location.City,
		"state":       location.State,
		"country":     location.Country,
		"postal_code": location.PostalCode,
		"updated_at":  location.UpdatedAt,
	})
	return err
}

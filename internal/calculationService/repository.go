package calculationService

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("calculation not found")

type Repository interface {
	Create(calc Calculation) (Calculation, error)
	GetAll() ([]Calculation, error)
	GetByID(id string) (Calculation, error)
	Update(calc Calculation) (Calculation, error)
	Delete(id string) error
}

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(calc Calculation) (Calculation, error) {
	if err := r.db.Create(&calc).Error; err != nil {
		return Calculation{}, err
	}
	return calc, nil
}

func (r *gormRepository) GetAll() ([]Calculation, error) {
	var list []Calculation
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gormRepository) GetByID(id string) (Calculation, error) {
	var calc Calculation
	if err := r.db.First(&calc, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Calculation{}, ErrNotFound
		}
		return Calculation{}, err
	}
	return calc, nil
}

func (r *gormRepository) Update(calc Calculation) (Calculation, error) {
	// ensure exists
	_, err := r.GetByID(calc.ID)
	if err != nil {
		return Calculation{}, err
	}
	if err := r.db.Save(&calc).Error; err != nil {
		return Calculation{}, err
	}
	return calc, nil
}

func (r *gormRepository) Delete(id string) error {
	if err := r.db.Delete(&Calculation{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

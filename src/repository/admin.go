package repository

import (
	model "api/src/model"
	"log"

	"gorm.io/gorm"
)

type IAdminRepository interface {
	// 企業登録
	InsertCompany(tx *gorm.DB, m *model.Company) error
}

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) IAdminRepository {
	return &AdminRepository{db}
}

// 企業登録
func (a *AdminRepository) InsertCompany(tx *gorm.DB, m *model.Company) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

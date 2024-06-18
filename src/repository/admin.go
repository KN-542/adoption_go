package repository

import (
	"api/src/model/ddl"
	"log"

	"gorm.io/gorm"
)

type IAdminRepository interface {
	// 企業登録
	Insert(tx *gorm.DB, m *ddl.Company) error
}

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) IAdminRepository {
	return &AdminRepository{db}
}

// 企業登録
func (a *AdminRepository) Insert(tx *gorm.DB, m *ddl.Company) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

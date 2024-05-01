package repository

import (
	model "api/src/model"
	"log"

	"gorm.io/gorm"
)

type IRoleRepository interface {
	// 登録
	Insert(tx *gorm.DB, m *model.CustomRole) error
	// 付与ロール登録
	InsertAssociation(tx *gorm.DB, m *model.RoleAssociation) error
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{db}
}

// 登録
func (r *RoleRepository) Insert(tx *gorm.DB, m *model.CustomRole) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 付与ロール登録
func (r *RoleRepository) InsertAssociation(tx *gorm.DB, m *model.RoleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

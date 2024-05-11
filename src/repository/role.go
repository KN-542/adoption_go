package repository

import (
	"api/src/model/ddl"
	"api/src/model/entity"
	"log"

	"gorm.io/gorm"
)

type IRoleRepository interface {
	// 登録
	Insert(tx *gorm.DB, m *ddl.CustomRole) (*entity.CustomRole, error)
	// 付与ロール登録
	InsertAssociation(tx *gorm.DB, m *ddl.RoleAssociation) error
	// 該当ロールのマスタID取得
	GetRoleIDs(m *ddl.CustomRole) ([]entity.RoleAssociation, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{db}
}

// 登録
func (r *RoleRepository) Insert(tx *gorm.DB, m *ddl.CustomRole) (*entity.CustomRole, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &entity.CustomRole{
		CustomRole: *m,
	}, nil
}

// 付与ロール登録
func (r *RoleRepository) InsertAssociation(tx *gorm.DB, m *ddl.RoleAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 該当ロールのマスタID取得
func (r *RoleRepository) GetRoleIDs(m *ddl.CustomRole) ([]entity.RoleAssociation, error) {
	var res []entity.RoleAssociation

	query := r.db.Model(&entity.RoleAssociation{}).
		Select(`master_role_id`).
		Where("role_id = ?", m.ID)

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

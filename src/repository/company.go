package repository

import (
	"api/src/model/ddl"
	"api/src/model/entity"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type ICompanyRepository interface {
	// 登録
	Insert(tx *gorm.DB, m *ddl.Company) (*entity.Company, error)
	// 検索
	Search(m *ddl.Company) ([]entity.Company, error)
	// 企業名重複確認
	IsDuplName(m *ddl.Company) error
}

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) ICompanyRepository {
	return &CompanyRepository{db}
}

// 登録
func (r *CompanyRepository) Insert(tx *gorm.DB, m *ddl.Company) (*entity.Company, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &entity.Company{
		Company: *m,
	}, nil
}

// 検索
func (r *CompanyRepository) Search(m *ddl.Company) ([]entity.Company, error) {
	var res []entity.Company

	query := r.db.Table("t_company").
		Select(`
			t_company.hash_key,
			t_company.name
		`)

	if m.Name != "" {
		query = query.Where("t_company.name LIKE ?", "%"+m.Name+"%")
	}

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 企業名重複確認
func (r *CompanyRepository) IsDuplName(m *ddl.Company) error {
	var count int64
	if err := r.db.Model(&ddl.Company{}).Where(&ddl.Company{
		Name: m.Name,
	}).Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("duplicate company name")
	}

	return nil
}

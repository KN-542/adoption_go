package repository

import (
	"log"

	"gorm.io/gorm"
)

type IDBRepository interface {
	// トランザクション開始
	TxStart() (*gorm.DB, error)
	// トランザクションコミット
	TxCommit(tx *gorm.DB) error
	// トランザクションロールバック
	TxRollback(tx *gorm.DB) error
}

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) IDBRepository {
	return &DBRepository{db}
}

// トランザクション開始
func (d *DBRepository) TxStart() (*gorm.DB, error) {
	tx := d.db.Begin()
	if err := tx.Error; err != nil {
		log.Printf("%v", err)
		return nil, tx.Error
	}

	return tx, nil
}

// トランザクションコミット
func (d *DBRepository) TxCommit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// トランザクションロールバック
func (d *DBRepository) TxRollback(tx *gorm.DB) error {
	if err := tx.Rollback().Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

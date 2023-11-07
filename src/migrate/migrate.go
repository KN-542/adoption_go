package main

import (
	"api/src/infra/db"
	"api/src/model"
	"api/src/model/enum"
	"api/src/repository"
	"fmt"

	"gorm.io/gorm"
)

func main() {
	dbConn := db.NewDB()

	dbConn.AutoMigrate(
		&model.Role{},
		&model.Site{},
		&model.User{},
		&model.Applicant{},
	)

	// 初期マスタデータ
	CreateData(dbConn)

	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
}

func CreateData(db *gorm.DB) {
	r := repository.NewMasterRepository(db)

	// m_site
	r.InsertSite(
		&model.Site{
			SiteID:     int(enum.RECRUIT),
			SiteNameJa: "リクナビNEXT",
		},
	)
	r.InsertSite(
		&model.Site{
			SiteID:     int(enum.MYNAVI),
			SiteNameJa: "マイナビ",
		},
	)
	r.InsertSite(
		&model.Site{
			SiteID:     int(enum.DODA),
			SiteNameJa: "DODA",
		},
	)

	// m_role
	r.InsertRole(
		&model.Role{
			ID:     1,
			NameJa: "最高責任者",
		},
	)
	r.InsertRole(
		&model.Role{
			ID:     2,
			NameJa: "面接官",
		},
	)
}

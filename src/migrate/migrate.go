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
			SiteID:   int(enum.RECRUIT),
			SiteName: "リクナビNEXT",
		},
	)
	r.InsertSite(
		&model.Site{
			SiteID:   int(enum.MYNAVI),
			SiteName: "マイナビ",
		},
	)
	r.InsertSite(
		&model.Site{
			SiteID:   int(enum.DODA),
			SiteName: "DODA",
		},
	)

	// m_role
	r.InsertRole(
		&model.Role{
			ID:   1,
			Name: "Admin",
		},
	)
	r.InsertRole(
		&model.Role{
			ID:   2,
			Name: "Interviewer",
		},
	)
}

package main

import (
	"api/src/infra"
	"api/src/model"
	"api/src/model/enum"
	"api/src/repository"
	"fmt"

	"gorm.io/gorm"
)

func main() {
	dbConn := infra.NewDB()

	dbConn.AutoMigrate(
		&model.Role{},
		&model.Site{},
		&model.User{},
		&model.UserGroup{},
		&model.ApplicantStatus{},
		&model.Applicant{},
	)

	// 初期マスタデータ
	CreateData(dbConn)

	defer fmt.Println("Successfully Migrated")
	defer infra.CloseDB(dbConn)
}

func CreateData(db *gorm.DB) {
	r := repository.NewMasterRepository(db)

	// m_site
	r.InsertSite(
		&model.Site{
			ID:         uint(enum.RECRUIT),
			SiteNameJa: "リクナビNEXT",
		},
	)
	r.InsertSite(
		&model.Site{
			ID:         uint(enum.MYNAVI),
			SiteNameJa: "マイナビ",
		},
	)
	r.InsertSite(
		&model.Site{
			ID:         uint(enum.DODA),
			SiteNameJa: "DODA",
		},
	)
	r.InsertSite(
		&model.Site{
			ID:         uint(enum.OTHER),
			SiteNameJa: "その他",
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

	// m_applicant_status
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.SCHEDULE_UNANSWERED),
			StatusNameJa: "日程未回答",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.BOOK_CATEGORY_NOT_PRESENTED),
			StatusNameJa: "書類未提出",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_1),
			StatusNameJa: "1次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_2),
			StatusNameJa: "2次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_3),
			StatusNameJa: "3次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_4),
			StatusNameJa: "4次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_5),
			StatusNameJa: "5次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_6),
			StatusNameJa: "6次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_7),
			StatusNameJa: "7次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_8),
			StatusNameJa: "8次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_9),
			StatusNameJa: "9次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.INTERVIEW_10),
			StatusNameJa: "10次面接",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_1),
			StatusNameJa: "1次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_2),
			StatusNameJa: "2次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_3),
			StatusNameJa: "3次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_4),
			StatusNameJa: "4次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_5),
			StatusNameJa: "5次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_6),
			StatusNameJa: "6次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_7),
			StatusNameJa: "7次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_8),
			StatusNameJa: "8次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_9),
			StatusNameJa: "9次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.TASK_AFTER_INTERVIEW_10),
			StatusNameJa: "10次面接後課題",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_1),
			StatusNameJa: "1次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_2),
			StatusNameJa: "2次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_3),
			StatusNameJa: "3次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_4),
			StatusNameJa: "4次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_5),
			StatusNameJa: "5次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_6),
			StatusNameJa: "6次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_7),
			StatusNameJa: "7次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_8),
			StatusNameJa: "8次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_9),
			StatusNameJa: "9次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_INTERVIEW_10),
			StatusNameJa: "10次面接落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.OFFER),
			StatusNameJa: "内定",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.OFFER_COMMITMENT),
			StatusNameJa: "内定承諾",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.Failing_TO_PASS_DOCUMENTS),
			StatusNameJa: "書類選考落ち",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.WITHDRAWAL),
			StatusNameJa: "選考辞退",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.OFFER_DISMISSAL),
			StatusNameJa: "内定辞退",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			ID:           uint(enum.OFFER_COMMITMENT_DISMISSAL),
			StatusNameJa: "内定承諾後辞退",
		},
	)
}

package main

import (
	"api/src/infra"
	"api/src/model"
	"api/src/model/enum"
	"api/src/repository"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func main() {
	dbConn := infra.NewDB()

	dbConn.AutoMigrate(
		&model.Role{},
		&model.Site{},
		&model.ApplicantStatus{},
		&model.CalendarFreqStatus{},
		&model.User{},
		&model.UserGroup{},
		&model.UserSchedule{},
		&model.Applicant{},
	)

	/*
		論理名追加
	*/

	// m_role
	if err := AddTableComment(dbConn, "m_role", "ロールマスタ"); err != nil {
		log.Println(err)
	}
	mRole := map[string]string{
		"id":      "ID",
		"name_ja": "ロール名_日本語",
	}
	if err := AddColumnComments(dbConn, "m_role", mRole); err != nil {
		log.Println(err)
	}

	// m_site
	if err := AddTableComment(dbConn, "m_site", "媒体マスタ"); err != nil {
		log.Println(err)
	}
	mSite := map[string]string{
		"id":           "ID",
		"site_name_ja": "媒体名_日本語",
	}
	if err := AddColumnComments(dbConn, "m_site", mSite); err != nil {
		log.Println(err)
	}

	// m_applicant_status
	if err := AddTableComment(dbConn, "m_applicant_status", "選考状況マスタ"); err != nil {
		log.Println(err)
	}
	mApplicantStatus := map[string]string{
		"id":             "ID",
		"status_name_ja": "ステータス名_日本語",
	}
	if err := AddColumnComments(dbConn, "m_applicant_status", mApplicantStatus); err != nil {
		log.Println(err)
	}

	// m_calendar_freq_status
	if err := AddTableComment(dbConn, "m_calendar_freq_status", "予定頻度マスタ"); err != nil {
		log.Println(err)
	}
	mCalendarFreqStatus := map[string]string{
		"id":      "ID",
		"freq":    "頻度",
		"name_ja": "名前_日本語",
	}
	if err := AddColumnComments(dbConn, "m_calendar_freq_status", mCalendarFreqStatus); err != nil {
		log.Println(err)
	}

	// t_user
	if err := AddTableComment(dbConn, "t_user", "ユーザー"); err != nil {
		log.Println(err)
	}
	user := map[string]string{
		"id":            "ID",
		"hash_key":      "ハッシュキー",
		"name":          "氏名",
		"email":         "メールアドレス",
		"password":      "パスワード(ハッシュ化)",
		"init_password": "初回パスワード(ハッシュ化)",
		"role_id":       "ロールID",
		"created_at":    "登録日時",
		"updated_at":    "更新日時",
	}
	if err := AddColumnComments(dbConn, "t_user", user); err != nil {
		log.Println(err)
	}

	// t_user_group
	if err := AddTableComment(dbConn, "t_user_group", "ユーザーグループ"); err != nil {
		log.Println(err)
	}
	userGroup := map[string]string{
		"id":         "ID",
		"hash_key":   "ハッシュキー",
		"name":       "グループ名",
		"users":      "所属ユーザー",
		"created_at": "登録日時",
		"updated_at": "更新日時",
	}
	if err := AddColumnComments(dbConn, "t_user_group", userGroup); err != nil {
		log.Println(err)
	}

	// t_user_schedule
	if err := AddTableComment(dbConn, "t_user_schedule", "ユーザー予定"); err != nil {
		log.Println(err)
	}
	userSchedule := map[string]string{
		"id":             "ID",
		"hash_key":       "ハッシュキー",
		"user_hash_keys": "ハッシュキー(ユーザー)",
		"title":          "タイトル",
		"freq_id":        "頻度ID",
		"start":          "開始時刻",
		"end":            "終了時刻",
		"created_at":     "登録日時",
		"updated_at":     "更新日時",
	}
	if err := AddColumnComments(dbConn, "t_user_schedule", userSchedule); err != nil {
		log.Println(err)
	}

	// t_applicant
	if err := AddTableComment(dbConn, "t_applicant", "応募者"); err != nil {
		log.Println(err)
	}
	applicant := map[string]string{
		"id":               "ID",
		"hash_key":         "ハッシュキー",
		"site_id":          "サイトID",
		"status":           "ステータス",
		"name":             "氏名",
		"email":            "メールアドレス",
		"tel":              "TEL",
		"age":              "年齢",
		"resume":           "履歴書",
		"curriculum_vitae": "職務経歴書",
		"google_meet_url":  "Google Meet URL",
		"desired_at":       "希望面接日時",
		"created_at":       "登録日時",
		"updated_at":       "更新日時",
	}
	if err := AddColumnComments(dbConn, "t_applicant", applicant); err != nil {
		log.Println(err)
	}

	// 初期マスタデータ
	CreateData(dbConn)

	defer fmt.Println("Successfully Migrated")
	defer infra.CloseDB(dbConn)
}

// 論理名追加
func AddTableComment(db *gorm.DB, tableName string, comment string) error {
	sql := fmt.Sprintf("COMMENT ON TABLE %s IS '%s';", tableName, comment)
	return db.Exec(sql).Error
}
func AddColumnComments(db *gorm.DB, tableName string, comments map[string]string) error {
	for column, comment := range comments {
		sql := fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", tableName, column, comment)
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

// 初期マスタデータ作成
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

	// m_calendar_freq_status
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			ID:     uint(enum.FREQ_NONE),
			Freq:   "",
			NameJa: "なし",
		},
	)
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			ID:     uint(enum.FREQ_DAILY),
			Freq:   "daily",
			NameJa: "毎日",
		},
	)
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			ID:     uint(enum.FREQ_WEEKLY),
			Freq:   "weekly",
			NameJa: "毎週",
		},
	)
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			ID:     uint(enum.FREQ_MONTHLY),
			Freq:   "monthly",
			NameJa: "毎月",
		},
	)
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			ID:     uint(enum.FREQ_YEARLY),
			Freq:   "yearly",
			NameJa: "毎年",
		},
	)
}

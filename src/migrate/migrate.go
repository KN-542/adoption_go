package main

import (
	"api/src/infra"
	"api/src/model"
	"api/src/model/enum"
	"api/src/repository"
	"flag"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func main() {
	dbConn := infra.NewDB()

	migrate := flag.Bool("migrate", false, "Set to true to run migrations")
	drop := flag.Bool("drop", false, "Set to true to drop all tables")
	flag.Parse()

	if *migrate {
		dbConn.AutoMigrate(
			// m
			&model.LoginType{},
			&model.Site{},
			&model.Role{},
			&model.ApplicantStatus{},
			&model.CalendarFreqStatus{},
			&model.ApplyVariable{},
			&model.ViewRoleOperation{},
			&model.OperationLogEvent{},
			&model.NoticeType{},
			&model.AnalysisTerm{},
			&model.HashKeyPre{},
			&model.S3NamePre{},
			// t
			&model.Company{},
			&model.CustomRole{},
			&model.User{},
			&model.UserGroup{},
			&model.UserGroupAssociation{},
			&model.UserSchedule{},
			&model.Applicant{},
			&model.MailTemplate{},
			&model.Variable{},
			&model.MailPreview{},
			&model.Notice{},
			&model.OperationLog{},
			&model.HistoryOfUploadApplicant{},
		)

		/*
			論理名追加
		*/

		// m_login_type
		if err := AddTableComment(dbConn, "m_login_type", "ログイン種別マスタ"); err != nil {
			log.Println(err)
		}
		mLoginType := map[string]string{
			"id":   "ID",
			"type": "ログイン種別",
			"path": "遷移パス",
		}
		if err := AddColumnComments(dbConn, "m_login_type", mLoginType); err != nil {
			log.Println(err)
		}

		// m_role
		if err := AddTableComment(dbConn, "m_role", "ロールマスタ"); err != nil {
			log.Println(err)
		}
		mRole := map[string]string{
			"id":      "ID",
			"name_ja": "ロール名_日本語",
			"name_en": "ロール名_英語",
		}
		if err := AddColumnComments(dbConn, "m_role", mRole); err != nil {
			log.Println(err)
		}

		// m_site
		if err := AddTableComment(dbConn, "m_site", "媒体マスタ"); err != nil {
			log.Println(err)
		}
		mSite := map[string]string{
			"id":        "ID",
			"site_name": "媒体名",
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
			"status_name_en": "ステータス名_英語",
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
			"name_en": "名前_英語",
		}
		if err := AddColumnComments(dbConn, "m_calendar_freq_status", mCalendarFreqStatus); err != nil {
			log.Println(err)
		}

		// m_apply_variable
		if err := AddTableComment(dbConn, "m_apply_variable", "適用変数種別マスタ"); err != nil {
			log.Println(err)
		}
		mApplyVariable := map[string]string{
			"id":   "ID",
			"name": "種別名",
		}
		if err := AddColumnComments(dbConn, "m_apply_variable", mApplyVariable); err != nil {
			log.Println(err)
		}

		// m_view_role_operation
		if err := AddTableComment(dbConn, "m_view_role_operation", "画面操作ロールマスタ"); err != nil {
			log.Println(err)
		}
		mViewRoleOperation := map[string]string{
			"id":    "ID",
			"name":  "操作名",
			"roles": "操作ロール",
		}
		if err := AddColumnComments(dbConn, "m_view_role_operation", mViewRoleOperation); err != nil {
			log.Println(err)
		}

		// m_operation_log_event
		if err := AddTableComment(dbConn, "m_operation_log_event", "操作ログイベントマスタ"); err != nil {
			log.Println(err)
		}
		mOperationLogEvent := map[string]string{
			"id":    "ID",
			"event": "通知内容",
		}
		if err := AddColumnComments(dbConn, "m_operation_log_event", mOperationLogEvent); err != nil {
			log.Println(err)
		}

		// m_notice
		if err := AddTableComment(dbConn, "m_notice", "通知マスタ"); err != nil {
			log.Println(err)
		}
		mNotice := map[string]string{
			"id":     "ID",
			"notice": "通知内容",
		}
		if err := AddColumnComments(dbConn, "m_notice", mNotice); err != nil {
			log.Println(err)
		}

		// m_analysis_term
		if err := AddTableComment(dbConn, "m_analysis_term", "分析項目マスタ"); err != nil {
			log.Println(err)
		}
		mAnalysisTerm := map[string]string{
			"id":      "ID",
			"term_ja": "項目_日本語",
			"term_en": "項目_英語",
		}
		if err := AddColumnComments(dbConn, "m_analysis_term", mAnalysisTerm); err != nil {
			log.Println(err)
		}

		// m_hash_key_pre
		if err := AddTableComment(dbConn, "m_hash_key_pre", "ハッシュキープレビューマスタ"); err != nil {
			log.Println(err)
		}
		mHashKeyPre := map[string]string{
			"id":  "ID",
			"pre": "プレビュー",
		}
		if err := AddColumnComments(dbConn, "m_hash_key_pre", mHashKeyPre); err != nil {
			log.Println(err)
		}

		// m_s3_name_pre
		if err := AddTableComment(dbConn, "m_s3_name_pre", "S3ファイル名プレビューマスタ"); err != nil {
			log.Println(err)
		}
		mS3NamePre := map[string]string{
			"id":  "ID",
			"pre": "プレビュー",
		}
		if err := AddColumnComments(dbConn, "m_s3_name_pre", mS3NamePre); err != nil {
			log.Println(err)
		}

		// t_company
		if err := AddTableComment(dbConn, "t_company", "企業"); err != nil {
			log.Println(err)
		}
		company := map[string]string{
			"id":         "ID",
			"hash_key":   "ハッシュキー",
			"name":       "企業名",
			"logo":       "ロゴファイル名",
			"delete_flg": "削除フラグ",
			"created_at": "登録日時",
			"updated_at": "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_company", company); err != nil {
			log.Println(err)
		}

		// t_role
		if err := AddTableComment(dbConn, "t_role", "ロール"); err != nil {
			log.Println(err)
		}
		role := map[string]string{
			"id":         "ID",
			"hash_key":   "ハッシュキー",
			"roles":      "付与ロール",
			"company_id": "企業ID",
			"edit_flg":   "編集可能フラグ",
			"delete_flg": "削除可能フラグ",
			"created_at": "登録日時",
			"updated_at": "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_role", role); err != nil {
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
			"refresh_token": "リフレッシュトークン",
			"company_id":    "企業ID",
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
			"company_id": "企業ID",
			"created_at": "登録日時",
			"updated_at": "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_user_group", userGroup); err != nil {
			log.Println(err)
		}

		// t_user_group_association
		if err := AddTableComment(dbConn, "t_user_group_association", "ユーザーグループ紐づけ"); err != nil {
			log.Println(err)
		}
		userGroupAssociation := map[string]string{
			"user_group_id": "ユーザーグループ",
			"user_id":       "ユーザーID",
		}
		if err := AddColumnComments(dbConn, "t_user_group_association", userGroupAssociation); err != nil {
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
			"interview_flg":  "面接フラグ",
			"company_id":     "企業ID",
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
			"outer_id":         "媒体側ID",
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
			"users":            "面接官",
			"calendar_id":      "カレンダーID",
			"company_id":       "企業ID",
			"created_at":       "登録日時",
			"updated_at":       "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_applicant", applicant); err != nil {
			log.Println(err)
		}

		// t_mail_template
		if err := AddTableComment(dbConn, "t_mail_template", "メールテンプレート"); err != nil {
			log.Println(err)
		}
		mailTemplate := map[string]string{
			"id":         "ID",
			"hash_key":   "ハッシュキー",
			"title":      "メールテンプレート名",
			"subject":    "件名",
			"template":   "テンプレート",
			"desc":       "説明",
			"company_id": "企業ID",
			"edit_flg":   "編集可能フラグ",
			"delete_flg": "削除可能フラグ",
			"created_at": "登録日時",
			"updated_at": "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_mail_template", mailTemplate); err != nil {
			log.Println(err)
		}

		// t_variable
		if err := AddTableComment(dbConn, "t_variable", "変数"); err != nil {
			log.Println(err)
		}
		variable := map[string]string{
			"id":         "ID",
			"hash_key":   "ハッシュキー",
			"title":      "変数タイトル",
			"json_name":  "変数格納Json名",
			"company_id": "企業ID",
			"edit_flg":   "編集可能フラグ",
			"delete_flg": "削除可能フラグ",
			"created_at": "登録日時",
			"updated_at": "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_variable", variable); err != nil {
			log.Println(err)
		}

		// t_mail_preview
		if err := AddTableComment(dbConn, "t_mail_preview", "メールプレビュー"); err != nil {
			log.Println(err)
		}
		mailPreview := map[string]string{
			"template_id": "テンプレートID",
			"variable_id": "変数ID",
			"hash_key":    "ハッシュキー",
			"title":       "メールプレビュー名",
			"desc":        "説明",
			"company_id":  "企業ID",
			"edit_flg":    "編集可能フラグ",
			"delete_flg":  "削除可能フラグ",
			"created_at":  "登録日時",
			"updated_at":  "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_mail_preview", mailPreview); err != nil {
			log.Println(err)
		}

		// t_notice
		if err := AddTableComment(dbConn, "t_notice", "通知"); err != nil {
			log.Println(err)
		}
		notice := map[string]string{
			"id":           "ID",
			"hash_key":     "ハッシュキー",
			"type":         "種別",
			"from_user_id": "通知元ユーザーID",
			"to_user_id":   "通知先ユーザーID",
			"company_id":   "企業ID",
			"created_at":   "登録日時",
			"updated_at":   "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_notice", notice); err != nil {
			log.Println(err)
		}

		// t_operation_log
		if err := AddTableComment(dbConn, "t_operation_log", "操作ログ"); err != nil {
			log.Println(err)
		}
		operationLog := map[string]string{
			"id":         "ID",
			"hash_key":   "ハッシュキー",
			"event_id":   "イベントID",
			"user_id":    "対象ユーザーID",
			"log":        "ログ",
			"company_id": "企業ID",
			"created_at": "登録日時",
			"updated_at": "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_operation_log", operationLog); err != nil {
			log.Println(err)
		}

		// t_history_of_upload_applicant
		if err := AddTableComment(dbConn, "t_history_of_upload_applicant", "応募者アップロード履歴"); err != nil {
			log.Println(err)
		}
		historyOfUploadApplicant := map[string]string{
			"history_id": "履歴ID",
			"csv":        "アップロードcsv",
		}
		if err := AddColumnComments(dbConn, "t_history_of_upload_applicant", historyOfUploadApplicant); err != nil {
			log.Println(err)
		}

		// 初期マスタデータ
		CreateData(dbConn)

		defer fmt.Println("Successfully Migrated")
		defer infra.CloseDB(dbConn)
	} else if *drop {
		dbConn.Migrator().DropTable(
			// m
			&model.LoginType{},
			&model.Site{},
			&model.Role{},
			&model.ApplicantStatus{},
			&model.CalendarFreqStatus{},
			&model.ApplyVariable{},
			&model.ViewRoleOperation{},
			&model.OperationLogEvent{},
			&model.NoticeType{},
			&model.AnalysisTerm{},
			&model.HashKeyPre{},
			&model.S3NamePre{},
			// t
			&model.Company{},
			&model.CustomRole{},
			&model.User{},
			&model.UserGroup{},
			&model.UserGroupAssociation{},
			&model.UserSchedule{},
			&model.Applicant{},
			&model.MailTemplate{},
			&model.Variable{},
			&model.MailPreview{},
			&model.Notice{},
			&model.OperationLog{},
			&model.HistoryOfUploadApplicant{},
		)
	}
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
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.RECRUIT),
			},
			SiteName: "リクナビNEXT",
		},
	)
	r.InsertSite(
		&model.Site{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.MYNAVI),
			},
			SiteName: "マイナビ",
		},
	)
	r.InsertSite(
		&model.Site{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.DODA),
			},
			SiteName: "DODA",
		},
	)
	r.InsertSite(
		&model.Site{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.OTHER),
			},
			SiteName: "その他",
		},
	)

	// m_role
	r.InsertRole(
		&model.Role{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: 1,
			},
			NameJa: "最高責任者",
			NameEn: "Admin",
		},
	)
	r.InsertRole(
		&model.Role{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: 2,
			},
			NameJa: "面接官",
			NameEn: "Interviewer",
		},
	)

	// m_applicant_status
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.SCHEDULE_UNANSWERED),
			},
			StatusNameJa: "日程未回答",
			StatusNameEn: "Schedule Unanswered",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.BOOK_CATEGORY_NOT_PRESENTED),
			},
			StatusNameJa: "書類未提出",
			StatusNameEn: "Document Not Submitted",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_1),
			},
			StatusNameJa: "1次面接",
			StatusNameEn: "First Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_2),
			},
			StatusNameJa: "2次面接",
			StatusNameEn: "Second Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_3),
			},
			StatusNameJa: "3次面接",
			StatusNameEn: "Third Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_4),
			},
			StatusNameJa: "4次面接",
			StatusNameEn: "Fourth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_5),
			},
			StatusNameJa: "5次面接",
			StatusNameEn: "Fifth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_6),
			},
			StatusNameJa: "6次面接",
			StatusNameEn: "Sixth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_7),
			},
			StatusNameJa: "7次面接",
			StatusNameEn: "Seventh Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_8),
			},
			StatusNameJa: "8次面接",
			StatusNameEn: "Eighth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_9),
			},
			StatusNameJa: "9次面接",
			StatusNameEn: "Ninth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.INTERVIEW_10),
			},
			StatusNameJa: "10次面接",
			StatusNameEn: "Tenth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_1),
			},
			StatusNameJa: "1次面接後課題",
			StatusNameEn: "Task After First Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_2),
			},
			StatusNameJa: "2次面接後課題",
			StatusNameEn: "Task After Second Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_3),
			},
			StatusNameJa: "3次面接後課題",
			StatusNameEn: "Task After Third Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_4),
			},
			StatusNameJa: "4次面接後課題",
			StatusNameEn: "Task After Fourth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_5),
			},
			StatusNameJa: "5次面接後課題",
			StatusNameEn: "Task After Fifth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_6),
			},
			StatusNameJa: "6次面接後課題",
			StatusNameEn: "Task After Sixth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_7),
			},
			StatusNameJa: "7次面接後課題",
			StatusNameEn: "Task After Seventh Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_8),
			},
			StatusNameJa: "8次面接後課題",
			StatusNameEn: "Task After Eighth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_9),
			},
			StatusNameJa: "9次面接後課題",
			StatusNameEn: "Task After Ninth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.TASK_AFTER_INTERVIEW_10),
			},
			StatusNameJa: "10次面接後課題",
			StatusNameEn: "Task After Tenth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_1),
			},
			StatusNameJa: "1次面接落ち",
			StatusNameEn: "Failed First Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_2),
			},
			StatusNameJa: "2次面接落ち",
			StatusNameEn: "Failed Second Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_3),
			},
			StatusNameJa: "3次面接落ち",
			StatusNameEn: "Failed Third Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_4),
			},
			StatusNameJa: "4次面接落ち",
			StatusNameEn: "Failed Fourth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_5),
			},
			StatusNameJa: "5次面接落ち",
			StatusNameEn: "Failed Fifth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_6),
			},
			StatusNameJa: "6次面接落ち",
			StatusNameEn: "Failed Sixth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_7),
			},
			StatusNameJa: "7次面接落ち",
			StatusNameEn: "Failed Seventh Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_8),
			},
			StatusNameJa: "8次面接落ち",
			StatusNameEn: "Failed Eighth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_9),
			},
			StatusNameJa: "9次面接落ち",
			StatusNameEn: "Failed Ninth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_INTERVIEW_10),
			},
			StatusNameJa: "10次面接落ち",
			StatusNameEn: "Failed Tenth Interview",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.OFFER),
			},
			StatusNameJa: "内定",
			StatusNameEn: "Offer",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.OFFER_COMMITMENT),
			},
			StatusNameJa: "内定承諾",
			StatusNameEn: "Offer Commitment",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.Failing_TO_PASS_DOCUMENTS),
			},
			StatusNameJa: "書類選考落ち",
			StatusNameEn: "Failed Document Screening",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.WITHDRAWAL),
			},
			StatusNameJa: "選考辞退",
			StatusNameEn: "Withdrawal from Selection",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.OFFER_DISMISSAL),
			},
			StatusNameJa: "内定辞退",
			StatusNameEn: "Offer Rejection",
		},
	)
	r.InsertApplicantStatus(
		&model.ApplicantStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.OFFER_COMMITMENT_DISMISSAL),
			},
			StatusNameJa: "内定承諾後辞退",
			StatusNameEn: "Post-Acceptance Withdrawal",
		},
	)

	// m_calendar_freq_status
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_NONE),
			},
			Freq:   "",
			NameJa: "なし",
			NameEn: "None",
		},
	)
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_DAILY),
			},
			Freq:   "daily",
			NameJa: "毎日",
			NameEn: "Daily",
		},
	)
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_WEEKLY),
			},
			Freq:   "weekly",
			NameJa: "毎週",
			NameEn: "Weekly",
		},
	)
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_MONTHLY),
			},
			Freq:   "monthly",
			NameJa: "毎月",
			NameEn: "Monthly",
		},
	)
	r.InsertCalendarFreqStatus(
		&model.CalendarFreqStatus{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_YEARLY),
			},
			Freq:   "yearly",
			NameJa: "毎年",
			NameEn: "Yearly",
		},
	)
}

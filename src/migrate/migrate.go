package main

import (
	"api/src/infra"
	"api/src/model"
	"api/src/model/enum"
	"api/src/repository"
	"api/src/service"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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
			&model.OperationLogEvent{},
			&model.NoticeType{},
			&model.AnalysisTerm{},
			&model.HashKeyPre{},
			&model.S3NamePre{},
			// t
			&model.Company{},
			&model.CustomRole{},
			&model.RoleAssociation{},
			&model.User{},
			&model.UserGroup{},
			&model.UserGroupAssociation{},
			&model.UserSchedule{},
			&model.UserScheduleAssociation{},
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
			"id":        "ID",
			"name_ja":   "ロール名_日本語",
			"name_en":   "ロール名_英語",
			"role_type": "ロール種別",
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
			"name":       "ロール名",
			"company_id": "企業ID",
			"edit_flg":   "編集保護",
			"delete_flg": "削除保護",
			"created_at": "登録日時",
			"updated_at": "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_role", role); err != nil {
			log.Println(err)
		}

		// t_role_association
		if err := AddTableComment(dbConn, "t_role_association", "付与ロール"); err != nil {
			log.Println(err)
		}
		roleAssociation := map[string]string{
			"role_id":        "ロールID",
			"master_role_id": "マスターロールID",
		}
		if err := AddColumnComments(dbConn, "t_role_association", roleAssociation); err != nil {
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
			"user_type":     "ユーザー種別",
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
			"id":            "ID",
			"hash_key":      "ハッシュキー",
			"title":         "タイトル",
			"freq_id":       "頻度ID",
			"interview_flg": "面接フラグ",
			"company_id":    "企業ID",
			"start":         "開始時刻",
			"end":           "終了時刻",
			"created_at":    "登録日時",
			"updated_at":    "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_user_schedule", userSchedule); err != nil {
			log.Println(err)
		}

		// t_user_schedule_association
		if err := AddTableComment(dbConn, "t_user_schedule_association", "ユーザー予定紐づけ"); err != nil {
			log.Println(err)
		}
		userScheduleAssociation := map[string]string{
			"user_schedule_id": "ユーザー予定ID",
			"user_id":          "ユーザーID",
		}
		if err := AddColumnComments(dbConn, "t_user_schedule_association", userScheduleAssociation); err != nil {
			log.Println(err)
		}

		// t_applicant
		year := time.Now().Year()
		month := time.Now().Month()
		if err := AddTableComment(
			dbConn,
			fmt.Sprintf("t_applicant_%d_%02d", year, month),
			fmt.Sprintf("応募者_%d_%02d", year, month),
		); err != nil {
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
			"calendar_id":      "カレンダーID",
			"company_id":       "企業ID",
			"created_at":       "登録日時",
			"updated_at":       "更新日時",
		}
		if err := AddColumnComments(dbConn, fmt.Sprintf("t_applicant_%d_%02d", year, month), applicant); err != nil {
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
			"edit_flg":   "編集保護",
			"delete_flg": "削除保護",
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
			"edit_flg":   "編集保護",
			"delete_flg": "削除保護",
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
			"edit_flg":    "編集保護",
			"delete_flg":  "削除保護",
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
			&model.OperationLogEvent{},
			&model.NoticeType{},
			&model.AnalysisTerm{},
			&model.HashKeyPre{},
			&model.S3NamePre{},
			// t
			&model.Company{},
			&model.CustomRole{},
			&model.RoleAssociation{},
			&model.User{},
			&model.UserGroup{},
			&model.UserGroupAssociation{},
			&model.UserSchedule{},
			&model.UserScheduleAssociation{},
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

// 初期データ作成
func CreateData(db *gorm.DB) {
	master := repository.NewMasterRepository(db)
	admin := repository.NewAdminRepository(db)
	role := repository.NewRoleRepository(db)
	user := repository.NewUserRepository(db)

	tx := db.Begin()
	if err := tx.Error; err != nil {
		log.Printf("%v", err)
		return
	}

	// m_login_type
	loginTypes := []*model.LoginType{
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.LOGIN_TYPE_ADMIN),
			},
			Type: "システム管理者",
			Path: "admin",
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.LOGIN_TYPE_MANAGEMENT),
			},
			Type: "一般",
			Path: "management",
		},
	}
	for _, row := range loginTypes {
		if err := master.InsertLoginType(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_site
	sites := []*model.Site{
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.RECRUIT),
			},
			SiteName: "リクナビNEXT",
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.MYNAVI),
			},
			SiteName: "マイナビ",
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.DODA),
			},
			SiteName: "DODA",
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.OTHER),
			},
			SiteName: "その他",
		},
	}
	for _, row := range sites {
		if err := master.InsertSite(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_role
	roles := []*model.Role{
		// admin_ロール関連
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_ROLE_CREATE),
			},
			NameJa:   "システム管理者ロール作成",
			NameEn:   "AdminRoleCreate",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_ROLE_READ),
			},
			NameJa:   "システム管理者ロール閲覧",
			NameEn:   "AdminRoleRead",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_ROLE_DETAIL_READ),
			},
			NameJa:   "システム管理者ロール詳細閲覧",
			NameEn:   "AdminRoleDetailRead",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_ROLE_EDIT),
			},
			NameJa:   "システム管理者ロール編集",
			NameEn:   "AdminRoleEdit",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_ROLE_DELETE),
			},
			NameJa:   "システム管理者ロール削除",
			NameEn:   "AdminRoleDelete",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_ROLE_ASSIGN),
			},
			NameJa:   "システム管理者ロール変更",
			NameEn:   "AdminRoleAssign",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		// admin_企業関連
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_COMPANY_CREATE),
			},
			NameJa:   "システム管理者企業作成",
			NameEn:   "AdminCompanyCreate",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_COMPANY_READ),
			},
			NameJa:   "システム管理者企業閲覧",
			NameEn:   "AdminCompanyRead",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_COMPANY_DETAIL_READ),
			},
			NameJa:   "システム管理者企業詳細閲覧",
			NameEn:   "AdminCompanyDetailRead",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_COMPANY_EDIT),
			},
			NameJa:   "システム管理者企業編集",
			NameEn:   "AdminCompanyEdit",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_ADMIN_COMPANY_DELETE),
			},
			NameJa:   "システム管理者企業削除",
			NameEn:   "AdminCompanyDelete",
			RoleType: uint(enum.LOGIN_TYPE_ADMIN),
		},
		// management_ロール関連
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_ROLE_CREATE),
			},
			NameJa:   "管理者ロール作成",
			NameEn:   "ManagementRoleCreate",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_ROLE_READ),
			},
			NameJa:   "管理者ロール閲覧",
			NameEn:   "ManagementRoleRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_ROLE_DETAIL_READ),
			},
			NameJa:   "管理者ロール詳細閲覧",
			NameEn:   "ManagementRoleDetailRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_ROLE_EDIT),
			},
			NameJa:   "管理者ロール編集",
			NameEn:   "ManagementRoleEdit",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_ROLE_DELETE),
			},
			NameJa:   "管理者ロール削除",
			NameEn:   "ManagementRoleDelete",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_ROLE_ASSIGN),
			},
			NameJa:   "管理者ロール割振",
			NameEn:   "ManagementRoleAssign",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		// management_ユーザー関連
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_USER_CREATE),
			},
			NameJa:   "管理者ユーザー作成",
			NameEn:   "ManagementUserCreate",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_USER_READ),
			},
			NameJa:   "管理者ユーザー閲覧",
			NameEn:   "ManagementUserRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_USER_DETAIL_READ),
			},
			NameJa:   "管理者ユーザー詳細閲覧",
			NameEn:   "ManagementUserDetailRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_USER_EDIT),
			},
			NameJa:   "管理者ユーザー編集",
			NameEn:   "ManagementUserEdit",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_USER_DELETE),
			},
			NameJa:   "管理者ユーザー削除",
			NameEn:   "ManagementUserDelete",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		// management_チーム関連
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_TEAM_CREATE),
			},
			NameJa:   "管理者チーム作成",
			NameEn:   "ManagementTeamCreate",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_TEAM_READ),
			},
			NameJa:   "管理者チーム閲覧",
			NameEn:   "ManagementTeamRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_TEAM_DETAIL_READ),
			},
			NameJa:   "管理者チーム詳細閲覧",
			NameEn:   "ManagementTeamDetailRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_TEAM_EDIT),
			},
			NameJa:   "管理者チーム編集",
			NameEn:   "ManagementTeamEdit",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_TEAM_DELETE),
			},
			NameJa:   "管理者チーム削除",
			NameEn:   "ManagementTeamDelete",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		// management_カレンダー関連
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_CALENDAR_CREATE),
			},
			NameJa:   "管理者カレンダー作成",
			NameEn:   "ManagementCalendarCreate",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_CALENDAR_READ),
			},
			NameJa:   "管理者カレンダー閲覧",
			NameEn:   "ManagementCalendarRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_CALENDAR_DETAIL_READ),
			},
			NameJa:   "管理者カレンダー詳細閲覧",
			NameEn:   "ManagementCalendarDetailRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_CALENDAR_EDIT),
			},
			NameJa:   "管理者カレンダー編集",
			NameEn:   "ManagementCalendarEdit",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_CALENDAR_DELETE),
			},
			NameJa:   "管理者カレンダー削除",
			NameEn:   "ManagementCalendarDelete",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		// management_応募者関連
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_APPLICANT_CREATE),
			},
			NameJa:   "管理者応募者作成",
			NameEn:   "ManagementApplicantCreate",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_APPLICANT_READ),
			},
			NameJa:   "管理者応募者閲覧",
			NameEn:   "ManagementApplicantRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_APPLICANT_DETAIL_READ),
			},
			NameJa:   "管理者応募者詳細閲覧",
			NameEn:   "ManagementApplicantDetailRead",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_APPLICANT_DOWNLOAD),
			},
			NameJa:   "管理者応募者ダウンロード",
			NameEn:   "ManagementApplicantDownload",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_APPLICANT_CREATE_MEET_URL),
			},
			NameJa:   "管理者面接URL作成",
			NameEn:   "ManagementApplicantCreateMeetURL",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.ROLE_MANAGEMENT_APPLICANT_ASSIGN_USER),
			},
			NameJa:   "管理者応募者割振",
			NameEn:   "ManagementApplicantAssignUser",
			RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
		},
	}
	for _, row := range roles {
		if err := master.InsertRole(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_calendar_freq_status
	calendarFreqStatus := []*model.CalendarFreqStatus{
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_NONE),
			},
			Freq:   "",
			NameJa: "なし",
			NameEn: "None",
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_DAILY),
			},
			Freq:   "daily",
			NameJa: "毎日",
			NameEn: "Daily",
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_WEEKLY),
			},
			Freq:   "weekly",
			NameJa: "毎週",
			NameEn: "Weekly",
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_MONTHLY),
			},
			Freq:   "monthly",
			NameJa: "毎月",
			NameEn: "Monthly",
		},
		{
			AbstractMasterModel: model.AbstractMasterModel{
				ID: uint(enum.FREQ_YEARLY),
			},
			Freq:   "yearly",
			NameJa: "毎年",
			NameEn: "Yearly",
		},
	}
	for _, row := range calendarFreqStatus {
		if err := master.InsertCalendarFreqStatus(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// t_company
	companies := []*model.Company{
		{
			Name: "管理者",
		},
	}
	for _, row := range companies {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = string(enum.PRE_COMPANY) + "_" + *hash

		if err := admin.InsertCompany(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// t_role
	customRoles := []*model.CustomRole{
		{
			AbstractTransactionModel: model.AbstractTransactionModel{
				CompanyID: 1,
			},
			AbstractTransactionFlgModel: model.AbstractTransactionFlgModel{
				EditFlg:   uint(enum.ON),
				DeleteFlg: uint(enum.ON),
			},
			Name: "Initial full-rights role (name change recommended)",
		},
	}
	for _, row := range customRoles {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = string(enum.PRE_ROLE) + "_" + *hash

		if err := role.Insert(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// t_role_association
	for _, row := range roles {
		if row.RoleType == uint(enum.LOGIN_TYPE_ADMIN) {
			if err := role.InsertAssociation(tx, &model.RoleAssociation{
				RoleID:       1,
				MasterRoleID: row.ID,
			}); err != nil {
				if err := tx.Rollback().Error; err != nil {
					log.Printf("%v", err)
					return
				}
				return
			}
		}
	}

	// t_user
	users := []*model.User{
		{
			AbstractTransactionModel: model.AbstractTransactionModel{
				CompanyID: 1,
			},
			Name:     "Initial user (name change recommended)",
			Email: os.Getenv("INIT_USER_EMAIL"),
			RoleID:   1,
			UserType: uint(enum.LOGIN_TYPE_ADMIN),
		},
	}
	for index, row := range users {
		password, hashPassword, _ := service.GenerateHash(8, 16)
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = string(enum.PRE_USER) + "_" + *hash
		row.Password = *hashPassword
		row.InitPassword = *hashPassword

		if err := user.Insert(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
		log.Printf("init password for user%v: %v", index+1, *password)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("%v", err)
		return
	}
}

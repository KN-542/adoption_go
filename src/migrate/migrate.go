package main

import (
	"api/src/infra"
	"api/src/model/ddl"
	"api/src/model/static"
	"api/src/repository"
	"api/src/service"
	"flag"
	"fmt"
	"log"
	"os"

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
			&ddl.LoginType{},
			&ddl.Site{},
			&ddl.Role{},
			&ddl.Sidebar{},
			&ddl.SidebarRoleAssociation{},
			&ddl.ScheduleFreqStatus{},
			&ddl.ApplyVariable{},
			&ddl.OperationLogEvent{},
			&ddl.NoticeType{},
			&ddl.AnalysisTerm{},
			&ddl.HashKeyPre{},
			&ddl.S3NamePre{},
			&ddl.SelectStatusEvent{},
			&ddl.AssignRule{},
			&ddl.AutoAssignRule{},
			&ddl.DocumentRule{},
			&ddl.Occupation{},
			&ddl.Processing{},
			// t
			&ddl.Company{},
			&ddl.CustomRole{},
			&ddl.RoleAssociation{},
			&ddl.User{},
			&ddl.UserRefreshTokenAssociation{},
			&ddl.Team{},
			&ddl.TeamAssociation{},
			&ddl.SelectStatus{},
			&ddl.TeamEvent{},
			&ddl.TeamEventEachInterview{},
			&ddl.TeamAutoAssignRule{},
			&ddl.TeamAssignPriority{},
			&ddl.TeamPerInterview{},
			&ddl.TeamAssignPossible{},
			&ddl.Schedule{},
			&ddl.ScheduleAssociation{},
			&ddl.Applicant{},
			&ddl.ApplicantUserAssociation{},
			&ddl.ApplicantType{},
			&ddl.ApplicantTypeAssociation{},
			&ddl.ApplicantScheduleAssociation{},
			&ddl.ApplicantResumeAssociation{},
			&ddl.ApplicantCurriculumVitaeAssociation{},
			&ddl.ApplicantURLAssociation{},
			&ddl.Manuscript{},
			&ddl.ManuscriptTeamAssociation{},
			&ddl.ManuscriptSiteAssociation{},
			&ddl.ManuscriptApplicantAssociation{},
			&ddl.MailTemplate{},
			&ddl.Variable{},
			&ddl.MailPreview{},
			&ddl.Notice{},
			&ddl.OperationLog{},
			&ddl.HistoryOfUploadApplicant{},
		)

		/*
			論理名追加
		*/

		// m_login_type
		if err := AddTableComment(dbConn, "m_login_type", "ログイン種別マスタ"); err != nil {
			log.Println(err)
		}
		mLoginType := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"type":     "ログイン種別",
			"path":     "遷移パス",
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
			"hash_key":  "ハッシュキー",
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
			"id":               "ID",
			"hash_key":         "ハッシュキー",
			"site_name":        "媒体名",
			"outer_id_index":   "媒体側ID_index",
			"name_index":       "氏名_index",
			"email_index":      "メールアドレス_index",
			"tel_index":        "TEL_index",
			"age_index":        "年齢_index",
			"manuscript_index": "原稿_index",
			"name_check_type":  "氏名_チェックタイプ",
			"num_of_column":    "カラム数",
		}
		if err := AddColumnComments(dbConn, "m_site", mSite); err != nil {
			log.Println(err)
		}

		// m_sidebar
		if err := AddTableComment(dbConn, "m_sidebar", "サイドバーマスタ"); err != nil {
			log.Println(err)
		}
		mSidebar := map[string]string{
			"id":        "ID",
			"hash_key":  "ハッシュキー",
			"name_ja":   "機能名_日本語",
			"name_en":   "機能名_英語",
			"path":      "遷移パス",
			"func_type": "機能種別",
		}
		if err := AddColumnComments(dbConn, "m_sidebar", mSidebar); err != nil {
			log.Println(err)
		}

		// m_sidebar_role_association
		if err := AddTableComment(dbConn, "m_sidebar_role_association", "サイドバーロール紐づけマスタ"); err != nil {
			log.Println(err)
		}
		mSidebarRoleAssociation := map[string]string{
			"sidebar_id": "サイドバーID",
			"role_id":    "操作可能ロールID",
		}
		if err := AddColumnComments(dbConn, "m_sidebar_role_association", mSidebarRoleAssociation); err != nil {
			log.Println(err)
		}

		// m_schedule_freq_status
		if err := AddTableComment(dbConn, "m_schedule_freq_status", "予定頻度マスタ"); err != nil {
			log.Println(err)
		}
		mScheduleFreqStatus := map[string]string{
			"id":        "ID",
			"hash_key":  "ハッシュキー",
			"freq_name": "頻度名",
			"name_ja":   "名前_日本語",
			"name_en":   "名前_英語",
		}
		if err := AddColumnComments(dbConn, "m_schedule_freq_status", mScheduleFreqStatus); err != nil {
			log.Println(err)
		}

		// m_apply_variable
		if err := AddTableComment(dbConn, "m_apply_variable", "適用変数種別マスタ"); err != nil {
			log.Println(err)
		}
		mApplyVariable := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"name":     "種別名",
		}
		if err := AddColumnComments(dbConn, "m_apply_variable", mApplyVariable); err != nil {
			log.Println(err)
		}

		// m_operation_log_event
		if err := AddTableComment(dbConn, "m_operation_log_event", "操作ログイベントマスタ"); err != nil {
			log.Println(err)
		}
		mOperationLogEvent := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"event":    "通知内容",
		}
		if err := AddColumnComments(dbConn, "m_operation_log_event", mOperationLogEvent); err != nil {
			log.Println(err)
		}

		// m_notice
		if err := AddTableComment(dbConn, "m_notice", "通知マスタ"); err != nil {
			log.Println(err)
		}
		mNotice := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"notice":   "通知内容",
		}
		if err := AddColumnComments(dbConn, "m_notice", mNotice); err != nil {
			log.Println(err)
		}

		// m_analysis_term
		if err := AddTableComment(dbConn, "m_analysis_term", "分析項目マスタ"); err != nil {
			log.Println(err)
		}
		mAnalysisTerm := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"term_ja":  "項目_日本語",
			"term_en":  "項目_英語",
		}
		if err := AddColumnComments(dbConn, "m_analysis_term", mAnalysisTerm); err != nil {
			log.Println(err)
		}

		// m_hash_key_pre
		if err := AddTableComment(dbConn, "m_hash_key_pre", "ハッシュキープレビューマスタ"); err != nil {
			log.Println(err)
		}
		mHashKeyPre := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"pre":      "プレビュー",
		}
		if err := AddColumnComments(dbConn, "m_hash_key_pre", mHashKeyPre); err != nil {
			log.Println(err)
		}

		// m_s3_name_pre
		if err := AddTableComment(dbConn, "m_s3_name_pre", "S3ファイル名プレビューマスタ"); err != nil {
			log.Println(err)
		}
		mS3NamePre := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"pre":      "プレビュー",
		}
		if err := AddColumnComments(dbConn, "m_s3_name_pre", mS3NamePre); err != nil {
			log.Println(err)
		}

		// m_select_status_event
		if err := AddTableComment(dbConn, "m_select_status_event", "応募者ステータスイベントマスタ"); err != nil {
			log.Println(err)
		}
		mSelectStatusEvent := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"desc_ja":  "説明_日本語",
			"desc_en":  "説明_英語",
		}
		if err := AddColumnComments(dbConn, "m_select_status_event", mSelectStatusEvent); err != nil {
			log.Println(err)
		}

		// m_assign_rule
		if err := AddTableComment(dbConn, "m_assign_rule", "面接アサインルールマスタ"); err != nil {
			log.Println(err)
		}
		mAssignRule := map[string]string{
			"id":                       "ID",
			"hash_key":                 "ハッシュキー",
			"desc_ja":                  "説明_日本語",
			"desc_en":                  "説明_英語",
			"additional_configuration": "追加設定必要性",
		}
		if err := AddColumnComments(dbConn, "m_assign_rule", mAssignRule); err != nil {
			log.Println(err)
		}

		// m_auto_assign_rule
		if err := AddTableComment(dbConn, "m_auto_assign_rule", "面接自動割り当てルールマスタ"); err != nil {
			log.Println(err)
		}
		mAutoAssignRule := map[string]string{
			"id":                       "ID",
			"hash_key":                 "ハッシュキー",
			"desc_ja":                  "説明_日本語",
			"desc_en":                  "説明_英語",
			"additional_configuration": "追加設定必要性",
		}
		if err := AddColumnComments(dbConn, "m_auto_assign_rule", mAutoAssignRule); err != nil {
			log.Println(err)
		}

		// m_document_rule
		if err := AddTableComment(dbConn, "m_document_rule", "書類提出ルールマスタ"); err != nil {
			log.Println(err)
		}
		mDocumentRule := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"rule_ja":  "ルール_日本語",
			"rule_en":  "ルール_英語",
		}
		if err := AddColumnComments(dbConn, "m_document_rule", mDocumentRule); err != nil {
			log.Println(err)
		}

		// m_occupation
		if err := AddTableComment(dbConn, "m_occupation", "職種マスタ"); err != nil {
			log.Println(err)
		}
		mOccupation := map[string]string{
			"id":       "ID",
			"hash_key": "ハッシュキー",
			"name_ja":  "職種名_日本語",
			"name_en":  "職種名_英語",
		}
		if err := AddColumnComments(dbConn, "m_occupation", mOccupation); err != nil {
			log.Println(err)
		}

		// m_interview_processing
		if err := AddTableComment(dbConn, "m_interview_processing", "面接過程マスタ"); err != nil {
			log.Println(err)
		}
		mProcessing := map[string]string{
			"id":         "ID",
			"hash_key":   "ハッシュキー",
			"processing": "過程",
			"desc_ja":    "説明_日本語",
			"desc_en":    "説明_英語",
		}
		if err := AddColumnComments(dbConn, "m_interview_processing", mProcessing); err != nil {
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
			"company_id":    "企業ID",
			"created_at":    "登録日時",
			"updated_at":    "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_user", user); err != nil {
			log.Println(err)
		}

		// t_user_refresh_token_association
		if err := AddTableComment(dbConn, "t_user_refresh_token_association", "リフレッシュトークン紐づけ"); err != nil {
			log.Println(err)
		}
		userRefreshTokenAssociation := map[string]string{
			"user_id":       "ユーザーID",
			"refresh_token": "リフレッシュトークン",
		}
		if err := AddColumnComments(dbConn, "t_user_refresh_token_association", userRefreshTokenAssociation); err != nil {
			log.Println(err)
		}

		// t_team
		if err := AddTableComment(dbConn, "t_team", "チーム"); err != nil {
			log.Println(err)
		}
		team := map[string]string{
			"id":               "ID",
			"hash_key":         "ハッシュキー",
			"name":             "チーム名",
			"num_of_interview": "最大面接回数",
			"rule_id":          "ルールID",
			"company_id":       "企業ID",
			"created_at":       "登録日時",
			"updated_at":       "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_team", team); err != nil {
			log.Println(err)
		}

		// t_team_association
		if err := AddTableComment(dbConn, "t_team_association", "チーム紐づけ"); err != nil {
			log.Println(err)
		}
		teamAssociation := map[string]string{
			"team_id": "チームID",
			"user_id": "ユーザーID",
		}
		if err := AddColumnComments(dbConn, "t_team_association", teamAssociation); err != nil {
			log.Println(err)
		}

		// t_select_status
		if err := AddTableComment(dbConn, "t_select_status", "選考状況"); err != nil {
			log.Println(err)
		}
		selectStatus := map[string]string{
			"team_id":     "チームID",
			"id":          "ID",
			"hash_key":    "ハッシュキー",
			"company_id":  "企業ID",
			"status_name": "ステータス名",
			"created_at":  "登録日時",
			"updated_at":  "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_select_status", selectStatus); err != nil {
			log.Println(err)
		}

		// t_team_event
		if err := AddTableComment(dbConn, "t_team_event", "チームイベント"); err != nil {
			log.Println(err)
		}
		teamEvent := map[string]string{
			"team_id":   "チームID",
			"event_id":  "イベントID",
			"status_id": "ステータスID",
		}
		if err := AddColumnComments(dbConn, "t_team_event", teamEvent); err != nil {
			log.Println(err)
		}

		// t_team_event_each_interview
		if err := AddTableComment(dbConn, "t_team_event_each_interview", "チーム面接毎イベント"); err != nil {
			log.Println(err)
		}
		teamEventEachInterview := map[string]string{
			"team_id":          "チームID",
			"num_of_interview": "面接回数",
			"process_id":       "面接過程",
			"status_id":        "ステータスID",
		}
		if err := AddColumnComments(dbConn, "t_team_event_each_interview", teamEventEachInterview); err != nil {
			log.Println(err)
		}

		// t_team_auto_assign_rule_association
		if err := AddTableComment(dbConn, "t_team_auto_assign_rule_association", "チーム面接自動割り当てルール紐づけ"); err != nil {
			log.Println(err)
		}
		teamAutoAssignRule := map[string]string{
			"team_id": "チームID",
			"rule_id": "ルールID",
		}
		if err := AddColumnComments(dbConn, "t_team_auto_assign_rule_association", teamAutoAssignRule); err != nil {
			log.Println(err)
		}

		// t_team_assign_priority
		if err := AddTableComment(dbConn, "t_team_assign_priority", "面接割り振り優先順位"); err != nil {
			log.Println(err)
		}
		teamAssignPriority := map[string]string{
			"team_id":  "チームID",
			"user_id":  "ユーザーID",
			"priority": "優先順位",
		}
		if err := AddColumnComments(dbConn, "t_team_assign_priority", teamAssignPriority); err != nil {
			log.Println(err)
		}

		// t_team_per_interview
		if err := AddTableComment(dbConn, "t_team_per_interview", "面接毎設定"); err != nil {
			log.Println(err)
		}
		teamPerInterview := map[string]string{
			"team_id":          "チームID",
			"num_of_interview": "面接回数",
			"user_min":         "最低人数",
		}
		if err := AddColumnComments(dbConn, "t_team_per_interview", teamPerInterview); err != nil {
			log.Println(err)
		}

		// t_team_assign_possible
		if err := AddTableComment(dbConn, "t_team_assign_possible", "面接毎参加可能者"); err != nil {
			log.Println(err)
		}
		teamAssignPossible := map[string]string{
			"team_id":          "チームID",
			"num_of_interview": "面接回数",
			"user_id":          "ユーザーID",
		}
		if err := AddColumnComments(dbConn, "t_team_assign_possible", teamAssignPossible); err != nil {
			log.Println(err)
		}

		// t_schedule
		if err := AddTableComment(dbConn, "t_schedule", "予定"); err != nil {
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
			"team_id":       "チームID",
			"created_at":    "登録日時",
			"updated_at":    "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_schedule", userSchedule); err != nil {
			log.Println(err)
		}

		// t_schedule_association
		if err := AddTableComment(dbConn, "t_schedule_association", "予定紐づけ"); err != nil {
			log.Println(err)
		}
		userScheduleAssociation := map[string]string{
			"schedule_id": "予定ID",
			"user_id":     "ユーザーID",
		}
		if err := AddColumnComments(dbConn, "t_schedule_association", userScheduleAssociation); err != nil {
			log.Println(err)
		}

		// t_applicant
		if err := AddTableComment(dbConn, "t_applicant", "応募者"); err != nil {
			log.Println(err)
		}
		applicant := map[string]string{
			"id":         "ID",
			"outer_id":   "媒体側ID",
			"hash_key":   "ハッシュキー",
			"site_id":    "サイトID",
			"status":     "ステータス",
			"name":       "氏名",
			"email":      "メールアドレス",
			"tel":        "TEL",
			"age":        "年齢",
			"company_id": "企業ID",
			"created_at": "登録日時",
			"updated_at": "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_applicant", applicant); err != nil {
			log.Println(err)
		}

		// t_applicant_user_association
		if err := AddTableComment(dbConn, "t_applicant_user_association", "応募者ユーザー紐づけ"); err != nil {
			log.Println(err)
		}
		applicantUserAssociation := map[string]string{
			"applicant_id": "応募者ID",
			"user_id":      "ユーザーID",
		}
		if err := AddColumnComments(dbConn, "t_applicant_user_association", applicantUserAssociation); err != nil {
			log.Println(err)
		}

		// t_applicant_type
		if err := AddTableComment(dbConn, "t_applicant_type", "応募者種別"); err != nil {
			log.Println(err)
		}
		applicantType := map[string]string{
			"id":            "ID",
			"hash_key":      "ハッシュキー",
			"name":          "種別名",
			"team_id":       "チームID",
			"rule_id":       "書類提出ルールID",
			"occupation_id": "職種ID",
			"company_id":    "企業ID",
			"created_at":    "登録日時",
			"updated_at":    "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_applicant_type", applicantType); err != nil {
			log.Println(err)
		}

		// t_applicant_type_association
		if err := AddTableComment(dbConn, "t_applicant_type_association", "応募者種別紐づけ"); err != nil {
			log.Println(err)
		}
		applicantTypeAssociation := map[string]string{
			"applicant_id": "応募者ID",
			"type_id":      "種別ID",
		}
		if err := AddColumnComments(dbConn, "t_applicant_type_association", applicantTypeAssociation); err != nil {
			log.Println(err)
		}

		// t_applicant_schedule_association
		if err := AddTableComment(dbConn, "t_applicant_schedule_association", "応募者面接予定紐づけ"); err != nil {
			log.Println(err)
		}
		applicantScheduleAssociation := map[string]string{
			"applicant_id": "応募者ID",
			"schedule_id":  "予定ID",
		}
		if err := AddColumnComments(dbConn, "t_applicant_schedule_association", applicantScheduleAssociation); err != nil {
			log.Println(err)
		}

		// t_applicant_resume_association
		if err := AddTableComment(dbConn, "t_applicant_resume_association", "応募者履歴書紐づけ"); err != nil {
			log.Println(err)
		}
		applicantResumeAssociation := map[string]string{
			"applicant_id": "応募者ID",
			"extension":    "拡張子",
		}
		if err := AddColumnComments(dbConn, "t_applicant_resume_association", applicantResumeAssociation); err != nil {
			log.Println(err)
		}

		// t_applicant_curriculum_vitae_association
		if err := AddTableComment(dbConn, "t_applicant_curriculum_vitae_association", "応募者職務経歴書紐づけ"); err != nil {
			log.Println(err)
		}
		applicantCurriculumVitaeAssociation := map[string]string{
			"applicant_id": "応募者ID",
			"extension":    "拡張子",
		}
		if err := AddColumnComments(dbConn, "t_applicant_curriculum_vitae_association", applicantCurriculumVitaeAssociation); err != nil {
			log.Println(err)
		}

		// t_applicant_url_association
		if err := AddTableComment(dbConn, "t_applicant_url_association", "応募者面接用URL紐づけ"); err != nil {
			log.Println(err)
		}
		applicantURLAssociation := map[string]string{
			"applicant_id": "応募者ID",
			"url":          "URL",
		}
		if err := AddColumnComments(dbConn, "t_applicant_url_association", applicantURLAssociation); err != nil {
			log.Println(err)
		}

		// t_manuscript
		if err := AddTableComment(dbConn, "t_manuscript", "原稿"); err != nil {
			log.Println(err)
		}
		manuscript := map[string]string{
			"id":         "ID",
			"hash_key":   "ハッシュキー",
			"company_id": "企業ID",
			"content":    "原稿内容",
			"created_at": "登録日時",
			"updated_at": "更新日時",
		}
		if err := AddColumnComments(dbConn, "t_manuscript", manuscript); err != nil {
			log.Println(err)
		}

		// t_manuscript_team_association
		if err := AddTableComment(dbConn, "t_manuscript_team_association", "原稿チーム紐づけ"); err != nil {
			log.Println(err)
		}
		manuscriptTeamAssociation := map[string]string{
			"manuscript_id": "原稿ID",
			"team_id":       "チームID",
		}
		if err := AddColumnComments(dbConn, "t_manuscript_team_association", manuscriptTeamAssociation); err != nil {
			log.Println(err)
		}

		// t_manuscript_applicant_association
		if err := AddTableComment(dbConn, "t_manuscript_applicant_association", "原稿応募者紐づけ"); err != nil {
			log.Println(err)
		}
		manuscriptApplicantAssociation := map[string]string{
			"manuscript_id": "原稿ID",
			"applicant_id":  "応募者ID",
		}
		if err := AddColumnComments(dbConn, "t_manuscript_applicant_association", manuscriptApplicantAssociation); err != nil {
			log.Println(err)
		}

		// t_manuscript_site_association
		if err := AddTableComment(dbConn, "t_manuscript_site_association", "原稿サイト紐づけ"); err != nil {
			log.Println(err)
		}
		manuscriptSiteAssociation := map[string]string{
			"manuscript_id": "原稿ID",
			"site_id":       "サイトID",
		}
		if err := AddColumnComments(dbConn, "t_manuscript_site_association", manuscriptSiteAssociation); err != nil {
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
			"commit_id":  "コミットID",
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
			&ddl.LoginType{},
			&ddl.Site{},
			&ddl.Role{},
			&ddl.Sidebar{},
			&ddl.SidebarRoleAssociation{},
			&ddl.ScheduleFreqStatus{},
			&ddl.ApplyVariable{},
			&ddl.OperationLogEvent{},
			&ddl.NoticeType{},
			&ddl.AnalysisTerm{},
			&ddl.HashKeyPre{},
			&ddl.S3NamePre{},
			&ddl.SelectStatusEvent{},
			&ddl.AssignRule{},
			&ddl.AutoAssignRule{},
			&ddl.DocumentRule{},
			&ddl.Occupation{},
			&ddl.Processing{},
			// t
			&ddl.Company{},
			&ddl.CustomRole{},
			&ddl.RoleAssociation{},
			&ddl.User{},
			&ddl.UserRefreshTokenAssociation{},
			&ddl.Team{},
			&ddl.TeamAssociation{},
			&ddl.SelectStatus{},
			&ddl.TeamEvent{},
			&ddl.TeamEventEachInterview{},
			&ddl.TeamAutoAssignRule{},
			&ddl.TeamAssignPriority{},
			&ddl.TeamPerInterview{},
			&ddl.TeamAssignPossible{},
			&ddl.Schedule{},
			&ddl.ScheduleAssociation{},
			&ddl.Applicant{},
			&ddl.ApplicantUserAssociation{},
			&ddl.ApplicantType{},
			&ddl.ApplicantTypeAssociation{},
			&ddl.ApplicantScheduleAssociation{},
			&ddl.ApplicantResumeAssociation{},
			&ddl.ApplicantCurriculumVitaeAssociation{},
			&ddl.ApplicantURLAssociation{},
			&ddl.Manuscript{},
			&ddl.ManuscriptTeamAssociation{},
			&ddl.ManuscriptSiteAssociation{},
			&ddl.ManuscriptApplicantAssociation{},
			&ddl.MailTemplate{},
			&ddl.Variable{},
			&ddl.MailPreview{},
			&ddl.Notice{},
			&ddl.OperationLog{},
			&ddl.HistoryOfUploadApplicant{},
		)

		defer fmt.Println("Successfully Deleted")
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
	loginTypes := []*ddl.LoginType{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.LOGIN_TYPE_ADMIN),
			},
			Type: "システム管理者",
			Path: "admin",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.LOGIN_TYPE_MANAGEMENT),
			},
			Type: "一般",
			Path: "management",
		},
	}
	for _, row := range loginTypes {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_login_type" + "_" + *hash
		if err := master.InsertLoginType(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_site
	sites := []*ddl.Site{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.RECRUIT),
			},
			SiteName:        "リクナビNEXT",
			FileName:        static.FILE_NAME_RECRUIT,
			OuterIDIndex:    static.INDEX_RECRUIT_OUTER_ID,
			NameIndex:       static.INDEX_RECRUIT_NAME,
			EmailIndex:      static.INDEX_RECRUIT_EMAIL,
			TELIndex:        static.INDEX_RECRUIT_TEL,
			AgeIndex:        static.INDEX_RECRUIT_AGE,
			ManuscriptIndex: static.INDEX_RECRUIT_MANUSCRIPT,
			NameCheckType:   1,
			NumOfColumn:     static.COLUMNS_RECRUIT,
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.MYNAVI),
			},
			SiteName:        "マイナビ",
			FileName:        static.FILE_NAME_MYNAVI,
			OuterIDIndex:    static.INDEX_MYNAVI_OUTER_ID,
			NameIndex:       static.INDEX_MYNAVI_NAME,
			EmailIndex:      static.INDEX_MYNAVI_EMAIL,
			TELIndex:        static.INDEX_MYNAVI_TEL,
			AgeIndex:        static.INDEX_MYNAVI_AGE,
			ManuscriptIndex: static.INDEX_MYNAVI_MANUSCRIPT,
			NumOfColumn:     static.COLUMNS_MYNAVI,
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.DODA),
			},
			SiteName:        "DODA",
			FileName:        static.FILE_NAME_DODA,
			OuterIDIndex:    static.INDEX_DODA_OUTER_ID,
			NameIndex:       static.INDEX_DODA_NAME,
			EmailIndex:      static.INDEX_DODA_EMAIL,
			TELIndex:        static.INDEX_DODA_TEL,
			AgeIndex:        static.INDEX_DODA_AGE,
			ManuscriptIndex: static.INDEX_DODA_MANUSCRIPT,
			NumOfColumn:     static.COLUMNS_DODA,
		},
	}
	for _, row := range sites {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_site" + "_" + *hash
		if err := master.InsertSite(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_select_status_event
	events := []*ddl.SelectStatusEvent{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.STATUS_EVENT_DECIDE_SCHEDULE),
			},
			DescJa: "応募者が日程調整のフォームを入力した時(初回面接前)",
			DescEn: "When an applicant fills out the scheduling form (initial interview)",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.STATUS_EVENT_SUBMIT_DOCUMENTS),
			},
			DescJa: "応募者がフォームから必要書類を提出した時(初回面接前)",
			DescEn: "When the applicant submits the required documents via the form (initial interview)",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.STATUS_EVENT_SUBMIT_DOCUMENTS_NOT_PASS),
			},
			DescJa: "書類不採用時(初回面接前)",
			DescEn: "When documents are not accepted (initial interview)",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.STATUS_EVENT_SUBMIT_DOCUMENTS_PASS),
			},
			DescJa: "書類通過時(初回面接前)",
			DescEn: "When documents are passed (initial interview)",
		},
	}
	for _, row := range events {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_select_status_event" + "_" + *hash
		if err := master.InsertSelectStatusEvent(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_assign_rule
	rules := []*ddl.AssignRule{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ASSIGN_RULE_MANUAL),
			},
			DescJa:                  "手動で割り当て",
			DescEn:                  "Manual assignment",
			AdditionalConfiguration: static.ASSIGN_RULE_CONFIG_UNREQUIRED,
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ASSIGN_RULE_AUTO),
			},
			DescJa:                  "自動で割り当て",
			DescEn:                  "Auto assignment",
			AdditionalConfiguration: static.ASSIGN_RULE_CONFIG_REQUIRED,
		},
	}
	for _, row := range rules {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_assign_rule" + "_" + *hash
		if err := master.InsertAssignRule(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_auto_assign_rule
	autoAssignRules := []*ddl.AutoAssignRule{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.AUTO_ASSIGN_RULE_RANDOM),
			},
			DescJa:                  "ランダム",
			DescEn:                  "Random",
			AdditionalConfiguration: static.AUTO_ASSIGN_RULE_CONFIG_UNREQUIRED,
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.AUTO_ASSIGN_RULE_ASC),
			},
			DescJa:                  "優先順位を事前決定",
			DescEn:                  "Pre-determine priorities",
			AdditionalConfiguration: static.AUTO_ASSIGN_RULE_CONFIG_REQUIRED,
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.AUTO_ASSIGN_RULE_DESC_SCHEDULE),
			},
			DescJa:                  "予定が少ない順",
			DescEn:                  "In order of least scheduled",
			AdditionalConfiguration: static.AUTO_ASSIGN_RULE_CONFIG_UNREQUIRED,
		},
	}
	for _, row := range autoAssignRules {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_auto_assign_rule" + "_" + *hash
		if err := master.InsertAutoAssignRule(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_interview_processing
	interviewResult := []*ddl.Processing{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.INTERVIEW_PROCESSING_NOW,
			},
			Processing: "面接予定",
			DescJa:     "日程確定時",
			DescEn:     "The schedule is finalized",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.INTERVIEW_PROCESSING_PASS,
			},
			Processing: "通過",
			DescJa:     "面接通過時",
			DescEn:     "Time to pass face-to-face contact",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.INTERVIEW_PROCESSING_FAIL,
			},
			Processing: "不採用",
			DescJa:     "面接不通過時",
			DescEn:     "the face connection fails",
		},
	}
	for _, row := range interviewResult {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_interview_processing" + "_" + *hash
		if err := master.InsertProcessing(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_document_rule
	documentRules := []*ddl.DocumentRule{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.DOCUMENT_RULE_REPUDIATE,
			},
			RuleJa: "提出、確認不要",
			RuleEn: "Submission and confirmation not required",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.DOCUMENT_RULE_REPUDIATE_CONFIRM,
			},
			RuleJa: "提出不要、確認必須",
			RuleEn: "Submission not required and confirmation required",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.DOCUMENT_RULE_REQUIRED_CONFIRM,
			},
			RuleJa: "提出、確認必須",
			RuleEn: "Submission and confirmation required",
		},
	}
	for _, row := range documentRules {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_document_rule" + "_" + *hash
		if err := master.InsertDocumentRule(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_occupation
	occupations := []*ddl.Occupation{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_ENGINEER,
			},
			NameJa: "エンジニア",
			NameEn: "Engineer",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_PROJECT_MANAGER,
			},
			NameJa: "プロジェクトマネージャー",
			NameEn: "Project Manager",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_DESIGNER,
			},
			NameJa: "デザイナー",
			NameEn: "Designer",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_ACCOUNTANT,
			},
			NameJa: "会計士",
			NameEn: "Accountant",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_HR_MANAGER,
			},
			NameJa: "人事マネージャー",
			NameEn: "HR Manager",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_CONSULTANT,
			},
			NameJa: "コンサルタント",
			NameEn: "Consultant",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_SALES_REPRESENTATIVE,
			},
			NameJa: "営業担当",
			NameEn: "Sales Representative",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_MARKETING_SPECIALIST,
			},
			NameJa: "マーケティングスペシャリスト",
			NameEn: "Marketing Specialist",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_CUSTOMER_SUPPORT,
			},
			NameJa: "カスタマーサポート",
			NameEn: "Customer Support",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: static.OCCUPATION_CEO,
			},
			NameJa: "最高経営責任者",
			NameEn: "CEO",
		},
	}
	for _, row := range occupations {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_occupation" + "_" + *hash
		if err := master.InsertOccupation(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_role
	roles := []*ddl.Role{
		// admin_ロール関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_ROLE_CREATE),
			},
			NameJa:   "システム管理者ロール作成",
			NameEn:   "AdminRoleCreate",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_ROLE_READ),
			},
			NameJa:   "システム管理者ロール閲覧",
			NameEn:   "AdminRoleRead",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_ROLE_DETAIL_READ),
			},
			NameJa:   "システム管理者ロール詳細閲覧",
			NameEn:   "AdminRoleDetailRead",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_ROLE_EDIT),
			},
			NameJa:   "システム管理者ロール編集",
			NameEn:   "AdminRoleEdit",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_ROLE_DELETE),
			},
			NameJa:   "システム管理者ロール削除",
			NameEn:   "AdminRoleDelete",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_ROLE_ASSIGN),
			},
			NameJa:   "システム管理者ロール変更",
			NameEn:   "AdminRoleAssign",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		// admin_企業関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_COMPANY_CREATE),
			},
			NameJa:   "システム管理者企業作成",
			NameEn:   "AdminCompanyCreate",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_COMPANY_READ),
			},
			NameJa:   "システム管理者企業閲覧",
			NameEn:   "AdminCompanyRead",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_COMPANY_DETAIL_READ),
			},
			NameJa:   "システム管理者企業詳細閲覧",
			NameEn:   "AdminCompanyDetailRead",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_COMPANY_EDIT),
			},
			NameJa:   "システム管理者企業編集",
			NameEn:   "AdminCompanyEdit",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_COMPANY_DELETE),
			},
			NameJa:   "システム管理者企業削除",
			NameEn:   "AdminCompanyDelete",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		// admin_ユーザー関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_USER_CREATE),
			},
			NameJa:   "システム管理者ユーザー作成",
			NameEn:   "AdminUserCreate",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_USER_READ),
			},
			NameJa:   "システム管理者ユーザー閲覧",
			NameEn:   "AdminUserRead",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_USER_DETAIL_READ),
			},
			NameJa:   "システム管理者ユーザー詳細閲覧",
			NameEn:   "AdminUserDetailRead",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_USER_EDIT),
			},
			NameJa:   "システム管理者ユーザー編集",
			NameEn:   "AdminUserEdit",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_USER_DELETE),
			},
			NameJa:   "システム管理者ユーザー削除",
			NameEn:   "AdminUserDelete",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		// admin_操作ログ関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_LOG_READ),
			},
			NameJa:   "システム管理者操作ログ閲覧",
			NameEn:   "AdminLogRead",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_ADMIN_LOG_DETAIL_READ),
			},
			NameJa:   "システム管理者操作ログ詳細閲覧",
			NameEn:   "AdminLogDetailRead",
			RoleType: uint(static.LOGIN_TYPE_ADMIN),
		},
		// management_ロール関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_ROLE_CREATE),
			},
			NameJa:   "管理者ロール作成",
			NameEn:   "ManagementRoleCreate",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_ROLE_READ),
			},
			NameJa:   "管理者ロール閲覧",
			NameEn:   "ManagementRoleRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_ROLE_DETAIL_READ),
			},
			NameJa:   "管理者ロール詳細閲覧",
			NameEn:   "ManagementRoleDetailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_ROLE_EDIT),
			},
			NameJa:   "管理者ロール編集",
			NameEn:   "ManagementRoleEdit",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_ROLE_DELETE),
			},
			NameJa:   "管理者ロール削除",
			NameEn:   "ManagementRoleDelete",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_ROLE_ASSIGN),
			},
			NameJa:   "管理者ロール割振",
			NameEn:   "ManagementRoleAssign",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_ユーザー関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_USER_CREATE),
			},
			NameJa:   "管理者ユーザー作成",
			NameEn:   "ManagementUserCreate",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_USER_READ),
			},
			NameJa:   "管理者ユーザー閲覧",
			NameEn:   "ManagementUserRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_USER_DETAIL_READ),
			},
			NameJa:   "管理者ユーザー詳細閲覧",
			NameEn:   "ManagementUserDetailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_USER_EDIT),
			},
			NameJa:   "管理者ユーザー編集",
			NameEn:   "ManagementUserEdit",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_USER_DELETE),
			},
			NameJa:   "管理者ユーザー削除",
			NameEn:   "ManagementUserDelete",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_チーム関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_TEAM_CREATE),
			},
			NameJa:   "管理者チーム作成",
			NameEn:   "ManagementTeamCreate",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_TEAM_READ),
			},
			NameJa:   "管理者チーム閲覧",
			NameEn:   "ManagementTeamRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_TEAM_DETAIL_READ),
			},
			NameJa:   "管理者チーム詳細閲覧",
			NameEn:   "ManagementTeamDetailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_TEAM_EDIT),
			},
			NameJa:   "管理者チーム編集",
			NameEn:   "ManagementTeamEdit",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_TEAM_DELETE),
			},
			NameJa:   "管理者チーム削除",
			NameEn:   "ManagementTeamDelete",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_予定関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_SCHEDULE_CREATE),
			},
			NameJa:   "管理者予定作成",
			NameEn:   "ManagementScheduleCreate",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_SCHEDULE_READ),
			},
			NameJa:   "管理者予定閲覧",
			NameEn:   "ManagementScheduleRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_SCHEDULE_DETAIL_READ),
			},
			NameJa:   "管理者予定詳細閲覧",
			NameEn:   "ManagementScheduleDetailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_SCHEDULE_EDIT),
			},
			NameJa:   "管理者予定編集",
			NameEn:   "ManagementScheduleEdit",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_SCHEDULE_DELETE),
			},
			NameJa:   "管理者予定削除",
			NameEn:   "ManagementScheduleDelete",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_応募者関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_APPLICANT_CREATE),
			},
			NameJa:   "管理者応募者作成",
			NameEn:   "ManagementApplicantCreate",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_APPLICANT_READ),
			},
			NameJa:   "管理者応募者閲覧",
			NameEn:   "ManagementApplicantRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_APPLICANT_DETAIL_READ),
			},
			NameJa:   "管理者応募者詳細閲覧",
			NameEn:   "ManagementApplicantDetailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_APPLICANT_DOWNLOAD),
			},
			NameJa:   "管理者応募者ダウンロード",
			NameEn:   "ManagementApplicantDownload",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_APPLICANT_CREATE_MEET_URL),
			},
			NameJa:   "管理者面接URL作成",
			NameEn:   "ManagementApplicantCreateMeetURL",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_APPLICANT_ASSIGN_USER),
			},
			NameJa:   "管理者応募者割振",
			NameEn:   "ManagementApplicantAssignUser",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_APPLICANT_SETTING_MANUSCRIPT),
			},
			NameJa:   "管理者応募者原稿割振",
			NameEn:   "ManagementApplicantAssignManuscript",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_APPLICANT_SETTING_TYPE),
			},
			NameJa:   "管理者応募者種別割振",
			NameEn:   "ManagementApplicantAssignType",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_APPLICANT_SETTING_STATUS),
			},
			NameJa:   "管理者応募者ステータス割振",
			NameEn:   "ManagementApplicantAssignStatus",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_原稿関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MANUSCRIPT_CREATE),
			},
			NameJa:   "管理者原稿作成",
			NameEn:   "ManagementManuscriptCreate",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MANUSCRIPT_READ),
			},
			NameJa:   "管理者原稿閲覧",
			NameEn:   "ManagementManuscriptRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MANUSCRIPT_DETAIL_READ),
			},
			NameJa:   "管理者原稿詳細閲覧",
			NameEn:   "ManagementManuscriptDetailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MANUSCRIPT_EDIT),
			},
			NameJa:   "管理者原稿編集",
			NameEn:   "ManagementManuscriptEdit",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MANUSCRIPT_DELETE),
			},
			NameJa:   "管理者原稿削除",
			NameEn:   "ManagementManuscriptDelete",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_メール関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MAIL_CREATE),
			},
			NameJa:   "管理者メール作成",
			NameEn:   "ManagementMailCreate",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MAIL_READ),
			},
			NameJa:   "管理者メール閲覧",
			NameEn:   "ManagementMailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MAIL_DETAIL_READ),
			},
			NameJa:   "管理者メール詳細閲覧",
			NameEn:   "ManagementMailDetailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MAIL_EDIT),
			},
			NameJa:   "管理者メール編集",
			NameEn:   "ManagementMailEdit",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_MAIL_DELETE),
			},
			NameJa:   "管理者メール削除",
			NameEn:   "ManagementMailDelete",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_変数関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_VARIABLE_CREATE),
			},
			NameJa:   "管理者変数作成",
			NameEn:   "ManagementVariableCreate",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_VARIABLE_READ),
			},
			NameJa:   "管理者変数閲覧",
			NameEn:   "ManagementVariableRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_VARIABLE_DETAIL_READ),
			},
			NameJa:   "管理者変数詳細閲覧",
			NameEn:   "ManagementVariableDetailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_VARIABLE_EDIT),
			},
			NameJa:   "管理者変数編集",
			NameEn:   "ManagementVariableEdit",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_VARIABLE_DELETE),
			},
			NameJa:   "管理者変数削除",
			NameEn:   "ManagementVariableDelete",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_分析関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_ANALYSIS_READ),
			},
			NameJa:   "管理者分析閲覧",
			NameEn:   "ManagementAnalysisRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_操作ログ関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_LOG_READ),
			},
			NameJa:   "管理者操作ログ閲覧",
			NameEn:   "ManagementLogRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_LOG_DETAIL_READ),
			},
			NameJa:   "管理者操作ログ詳細閲覧",
			NameEn:   "ManagementLogDetailRead",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		// management_設定関連
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_SETTING_COMPANY),
			},
			NameJa:   "管理者企業設定",
			NameEn:   "ManagementSettingCompany",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.ROLE_MANAGEMENT_SETTING_TEAM),
			},
			NameJa:   "管理者チーム設定",
			NameEn:   "ManagementSettingTeam",
			RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
	}
	for _, row := range roles {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_role" + "_" + *hash
		if err := master.InsertRole(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_sidebar
	sidebar := []*ddl.Sidebar{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_ADMIN_COMPANY),
			},
			NameJa:   "企業",
			NameEn:   "Companies",
			Path:     "/admin/company",
			FuncType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_ADMIN_USER),
			},
			NameJa:   "ユーザー",
			NameEn:   "Users",
			Path:     "/admin/user",
			FuncType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_ADMIN_ROLE),
			},
			NameJa:   "ロール",
			NameEn:   "Roles",
			Path:     "/admin/role",
			FuncType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_ADMIN_LOG),
			},
			NameJa:   "操作ログ",
			NameEn:   "Logs",
			Path:     "/admin/log",
			FuncType: uint(static.LOGIN_TYPE_ADMIN),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			},
			NameJa:   "応募者",
			NameEn:   "Applicant",
			Path:     "/management/applicant",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_USER),
			},
			NameJa:   "ユーザー",
			NameEn:   "Users",
			Path:     "/management/user",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_TEAM),
			},
			NameJa:   "チーム",
			NameEn:   "Teams",
			Path:     "/management/team",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_SCHEDULE),
			},
			NameJa:   "予定",
			NameEn:   "Schedules",
			Path:     "/management/schedule",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_ROLE),
			},
			NameJa:   "ロール",
			NameEn:   "Roles",
			Path:     "/management/role",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_MANUSCRIPT),
			},
			NameJa:   "原稿",
			NameEn:   "Manuscript",
			Path:     "/management/manuscript",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_MAIL),
			},
			NameJa:   "メールテンプレート",
			NameEn:   "Mail Templates",
			Path:     "/management/mail",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_VARIABLE),
			},
			NameJa:   "変数",
			NameEn:   "Variables",
			Path:     "/management/variable",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_ANALYSIS),
			},
			NameJa:   "分析",
			NameEn:   "Analysis",
			Path:     "/management/analysis",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.SIDEBAR_MANAGEMENT_LOG),
			},
			NameJa:   "操作ログ",
			NameEn:   "Logs",
			Path:     "/management/log",
			FuncType: uint(static.LOGIN_TYPE_MANAGEMENT),
		},
	}
	for _, row := range sidebar {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_sidebar" + "_" + *hash
		if err := master.InsertSidebar(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_sidebar_role_association
	sidebarRoleAssociation := []*ddl.SidebarRoleAssociation{
		// admin_企業関連
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_COMPANY),
			RoleID:    uint(static.ROLE_ADMIN_COMPANY_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_COMPANY),
			RoleID:    uint(static.ROLE_ADMIN_COMPANY_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_COMPANY),
			RoleID:    uint(static.ROLE_ADMIN_COMPANY_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_COMPANY),
			RoleID:    uint(static.ROLE_ADMIN_COMPANY_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_COMPANY),
			RoleID:    uint(static.ROLE_ADMIN_COMPANY_DELETE),
		},
		// admin_ユーザー関連
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_USER),
			RoleID:    uint(static.ROLE_ADMIN_USER_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_USER),
			RoleID:    uint(static.ROLE_ADMIN_USER_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_USER),
			RoleID:    uint(static.ROLE_ADMIN_USER_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_USER),
			RoleID:    uint(static.ROLE_ADMIN_USER_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_USER),
			RoleID:    uint(static.ROLE_ADMIN_USER_DELETE),
		},
		// admin_ロール関連
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_ROLE),
			RoleID:    uint(static.ROLE_ADMIN_ROLE_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_ROLE),
			RoleID:    uint(static.ROLE_ADMIN_ROLE_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_ROLE),
			RoleID:    uint(static.ROLE_ADMIN_ROLE_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_ROLE),
			RoleID:    uint(static.ROLE_ADMIN_ROLE_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_ROLE),
			RoleID:    uint(static.ROLE_ADMIN_ROLE_DELETE),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_ROLE),
			RoleID:    uint(static.ROLE_ADMIN_ROLE_ASSIGN),
		},
		// admin_操作ログ関連
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_LOG),
			RoleID:    uint(static.ROLE_ADMIN_LOG_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_ADMIN_LOG),
			RoleID:    uint(static.ROLE_ADMIN_LOG_DETAIL_READ),
		},
		// management_ロール関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_ROLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_ROLE_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_ROLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_ROLE_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_ROLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_ROLE_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_ROLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_ROLE_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_ROLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_ROLE_DELETE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_ROLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_ROLE_ASSIGN),
		},
		// management_ユーザー関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_USER),
			RoleID:    uint(static.ROLE_MANAGEMENT_USER_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_USER),
			RoleID:    uint(static.ROLE_MANAGEMENT_USER_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_USER),
			RoleID:    uint(static.ROLE_MANAGEMENT_USER_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_USER),
			RoleID:    uint(static.ROLE_MANAGEMENT_USER_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_USER),
			RoleID:    uint(static.ROLE_MANAGEMENT_USER_DELETE),
		},
		// management_チーム関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_TEAM),
			RoleID:    uint(static.ROLE_MANAGEMENT_TEAM_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_TEAM),
			RoleID:    uint(static.ROLE_MANAGEMENT_TEAM_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_TEAM),
			RoleID:    uint(static.ROLE_MANAGEMENT_TEAM_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_TEAM),
			RoleID:    uint(static.ROLE_MANAGEMENT_TEAM_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_TEAM),
			RoleID:    uint(static.ROLE_MANAGEMENT_TEAM_DELETE),
		},
		// management_予定関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_SCHEDULE),
			RoleID:    uint(static.ROLE_MANAGEMENT_SCHEDULE_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_SCHEDULE),
			RoleID:    uint(static.ROLE_MANAGEMENT_SCHEDULE_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_SCHEDULE),
			RoleID:    uint(static.ROLE_MANAGEMENT_SCHEDULE_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_SCHEDULE),
			RoleID:    uint(static.ROLE_MANAGEMENT_SCHEDULE_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_SCHEDULE),
			RoleID:    uint(static.ROLE_MANAGEMENT_SCHEDULE_DELETE),
		},
		// management_応募者関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			RoleID:    uint(static.ROLE_MANAGEMENT_APPLICANT_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			RoleID:    uint(static.ROLE_MANAGEMENT_APPLICANT_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			RoleID:    uint(static.ROLE_MANAGEMENT_APPLICANT_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			RoleID:    uint(static.ROLE_MANAGEMENT_APPLICANT_DOWNLOAD),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			RoleID:    uint(static.ROLE_MANAGEMENT_APPLICANT_CREATE_MEET_URL),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			RoleID:    uint(static.ROLE_MANAGEMENT_APPLICANT_ASSIGN_USER),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			RoleID:    uint(static.ROLE_MANAGEMENT_APPLICANT_SETTING_MANUSCRIPT),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			RoleID:    uint(static.ROLE_MANAGEMENT_APPLICANT_SETTING_TYPE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_APPLICANT),
			RoleID:    uint(static.ROLE_MANAGEMENT_APPLICANT_SETTING_STATUS),
		},
		// management_原稿関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MANUSCRIPT),
			RoleID:    uint(static.ROLE_MANAGEMENT_MANUSCRIPT_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MANUSCRIPT),
			RoleID:    uint(static.ROLE_MANAGEMENT_MANUSCRIPT_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MANUSCRIPT),
			RoleID:    uint(static.ROLE_MANAGEMENT_MANUSCRIPT_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MANUSCRIPT),
			RoleID:    uint(static.ROLE_MANAGEMENT_MANUSCRIPT_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MANUSCRIPT),
			RoleID:    uint(static.ROLE_MANAGEMENT_MANUSCRIPT_DELETE),
		},
		// management_メール関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MAIL),
			RoleID:    uint(static.ROLE_MANAGEMENT_MAIL_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MAIL),
			RoleID:    uint(static.ROLE_MANAGEMENT_MAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MAIL),
			RoleID:    uint(static.ROLE_MANAGEMENT_MAIL_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MAIL),
			RoleID:    uint(static.ROLE_MANAGEMENT_MAIL_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_MAIL),
			RoleID:    uint(static.ROLE_MANAGEMENT_MAIL_DELETE),
		},
		// management_変数関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_VARIABLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_VARIABLE_CREATE),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_VARIABLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_VARIABLE_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_VARIABLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_VARIABLE_DETAIL_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_VARIABLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_VARIABLE_EDIT),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_VARIABLE),
			RoleID:    uint(static.ROLE_MANAGEMENT_VARIABLE_DELETE),
		},
		// management_分析関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_ANALYSIS),
			RoleID:    uint(static.ROLE_MANAGEMENT_ANALYSIS_READ),
		},
		// management_操作ログ関連
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_LOG),
			RoleID:    uint(static.ROLE_MANAGEMENT_LOG_READ),
		},
		{
			SidebarID: uint(static.SIDEBAR_MANAGEMENT_LOG),
			RoleID:    uint(static.ROLE_MANAGEMENT_LOG_DETAIL_READ),
		},
	}
	for _, row := range sidebarRoleAssociation {
		if err := master.InsertSidebarRoleAssociation(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_hash_key_pre
	preHashKeys := []*ddl.HashKeyPre{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: 1,
			},
			Pre: string(static.PRE_COMPANY),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: 2,
			},
			Pre: string(static.PRE_ROLE),
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: 3,
			},
			Pre: string(static.PRE_USER),
		},
	}
	for _, row := range preHashKeys {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_hash_key_pre" + "_" + *hash
		if err := master.InsertHashKeyPre(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// m_schedule_freq_status
	scheduleFreqStatus := []*ddl.ScheduleFreqStatus{
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.FREQ_NONE),
			},
			FreqName: "",
			NameJa:   "なし",
			NameEn:   "None",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.FREQ_DAILY),
			},
			FreqName: "daily",
			NameJa:   "毎日",
			NameEn:   "Daily",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.FREQ_WEEKLY),
			},
			FreqName: "weekly",
			NameJa:   "毎週",
			NameEn:   "Weekly",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.FREQ_MONTHLY),
			},
			FreqName: "monthly",
			NameJa:   "毎月",
			NameEn:   "Monthly",
		},
		{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: uint(static.FREQ_YEARLY),
			},
			FreqName: "yearly",
			NameJa:   "毎年",
			NameEn:   "Yearly",
		},
	}
	for _, row := range scheduleFreqStatus {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = "m_schedule_freq_status" + "_" + *hash
		if err := master.InsertScheduleFreqStatus(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// t_company
	companies := []*ddl.Company{
		{
			Name: "管理者",
		},
	}
	for _, row := range companies {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = string(static.PRE_COMPANY) + "_" + *hash

		if err := admin.Insert(tx, row); err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// t_role
	customRoles := []*ddl.CustomRole{
		{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				CompanyID: 1,
			},
			AbstractTransactionFlgModel: ddl.AbstractTransactionFlgModel{
				EditFlg:   uint(static.ON),
				DeleteFlg: uint(static.ON),
			},
			Name: "Initial role",
		},
	}
	for _, row := range customRoles {
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = string(static.PRE_ROLE) + "_" + *hash

		_, err := role.Insert(tx, row)
		if err != nil {
			if err := tx.Rollback().Error; err != nil {
				log.Printf("%v", err)
				return
			}
			return
		}
	}

	// t_role_association
	for _, row := range roles {
		if row.RoleType == uint(static.LOGIN_TYPE_ADMIN) {
			if err := role.InsertAssociation(tx, &ddl.RoleAssociation{
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
	users := []*ddl.User{
		{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				CompanyID: 1,
			},
			Name:     "Initial user",
			Email:    os.Getenv("INIT_USER_EMAIL"),
			RoleID:   1,
			UserType: uint(static.LOGIN_TYPE_ADMIN),
		},
	}
	for index, row := range users {
		password, hashPassword, _ := service.GenerateHash(8, 16)
		_, hash, _ := service.GenerateHash(1, 25)
		row.HashKey = string(static.PRE_USER) + "_" + *hash
		row.Password = *hashPassword
		row.InitPassword = *hashPassword

		_, err := user.Insert(tx, row)
		if err != nil {
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

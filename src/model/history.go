package model

/*
t_operation_log
操作ログ
*/
type OperationLog struct {
	AbstractTransactionModel
	// イベントID
	EventID uint `json:"event_id"`
	// 対象ユーザーID
	UserID uint `json:"user_id"`
	// ログ
	Log string `json:"log" gorm:"type:text"`
	// イベント
	Event OperationLogEvent `gorm:"foreignKey:event_id;references:id"`
	// ユーザー(外部キー)
	User User `gorm:"foreignKey:user_id;references:id"`
}

/*
t_history_of_upload_applicant
応募者アップロード履歴
*/
type HistoryOfUploadApplicant struct {
	// 履歴ID
	HistoryID uint64 `json:"history_id" gorm:"primaryKey;AUTO_INCREMENT"`
	// アップロードcsv
	CSV string `json:"csv" gorm:"not null;check:csv <> '';type:text"`
	// 操作ログ(外部キー)
	Log OperationLog `gorm:"foreignKey:history_id;references:id"`
}

func (t OperationLog) TableName() string {
	return "t_operation_log"
}
func (t HistoryOfUploadApplicant) TableName() string {
	return "t_history_of_upload_applicant"
}

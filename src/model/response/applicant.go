package response

import "api/src/model/entity"

// 検索
type ApplicantSearch struct {
	List []entity.ApplicantSearch `json:"list"`
}

// サイト一覧取得
type ApplicantSites struct {
	List []entity.Site `json:"list"`
}

// 応募者ステータス一覧取得
type ApplicantStatusList struct {
	List []entity.ApplicantStatus `json:"list"`
}

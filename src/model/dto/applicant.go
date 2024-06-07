package dto

import "api/src/model/request"

type ApplicantSearch struct {
	request.ApplicantSearch
	// ユーザー
	UserIDs []uint64
}

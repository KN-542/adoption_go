package service

import (
	"api/src/model/ddl"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/repository"
	"api/src/validator"
	"context"
	"log"
	"net/http"
	"strconv"
)

type IRoleService interface {
	// ロールチェック
	Check(req *request.CheckRole) (bool, *response.Error)
	// 検索_企業ID
	SearchRoleByComapny(req *request.SearchRoleByComapny) (*response.SearchRoleByComapny, *response.Error)
}

type RoleService struct {
	role  repository.IRoleRepository
	redis repository.IRedisRepository
	v     validator.IRoleValidator
}

func NewRoleService(
	role repository.IRoleRepository,
	redis repository.IRedisRepository,
	v validator.IRoleValidator,
) IRoleService {
	return &RoleService{role, redis, v}
}

// ロールチェック
func (r *RoleService) Check(req *request.CheckRole) (bool, *response.Error) {
	// バリデーション
	if err := r.v.Check(req); err != nil {
		log.Printf("%v", err)
		return false, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// ロールID取得
	ctx := context.Background()
	role, roleErr := r.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_ROLE)
	if roleErr != nil {
		return false, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	roleID, roleParseErr := strconv.ParseUint(*role, 10, 64)
	if roleParseErr != nil {
		log.Printf("%v", roleParseErr)
		return false, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 該当ロールのマスタID取得
	roles, masterRoleErr := r.role.GetRoleIDs(&ddl.CustomRole{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: roleID,
		},
	})
	if masterRoleErr != nil {
		return false, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ロールの存在チェック
	for _, row := range roles {
		if row.MasterRoleID == req.ID {
			return true, nil
		}
	}

	return false, nil
}

// 検索_企業ID
func (r *RoleService) SearchRoleByComapny(req *request.SearchRoleByComapny) (*response.SearchRoleByComapny, *response.Error) {
	// ロールID取得
	ctx := context.Background()
	company, companyErr := r.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
	if companyErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	companyID, companyParseErr := strconv.ParseUint(*company, 10, 64)
	if companyParseErr != nil {
		log.Printf("%v", companyParseErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 検索
	roles, err := r.role.SearchByCompanyID(&ddl.CustomRole{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			CompanyID: companyID,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.SearchRoleByComapny{
		List: roles,
	}, nil
}

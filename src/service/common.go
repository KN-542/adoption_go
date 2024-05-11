package service

import (
	"api/resources/static"
	"api/src/model/ddl"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/repository"
	"api/src/validator"
	"context"
	"log"
	"net/http"
	"strconv"
)

type ICommonService interface {
	// サイドバー表示
	Sidebar(req *request.Sidebar) (*response.Sidebar, *response.Error)
	// 使用可能ロール一覧
	Roles(req *request.Roles) (*response.Roles, *response.Error)
}

type CommonService struct {
	master repository.IMasterRepository
	role   repository.IRoleRepository
	v      validator.ICommonValidator
	redis  repository.IRedisRepository
}

func NewCommonService(
	master repository.IMasterRepository,
	role repository.IRoleRepository,
	v validator.ICommonValidator,
	redis repository.IRedisRepository,
) ICommonService {
	return &CommonService{master, role, v, redis}
}

// サイドバー表示
func (c *CommonService) Sidebar(req *request.Sidebar) (*response.Sidebar, *response.Error) {
	// バリデーション
	if err := c.v.Sidebar(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ロールID、ログイン種別取得
	ctx := context.Background()
	role, err := c.redis.Get(ctx, req.HashKey, static.REDIS_USER_ROLE)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}
	roleID, err := strconv.ParseUint(*role, 10, 64)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	loginType_0, err := c.redis.Get(ctx, req.HashKey, static.REDIS_USER_LOGIN_TYPE)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}
	loginType, err := strconv.ParseUint(*loginType_0, 10, 64)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 該当ロールのマスタID取得
	roleAssociate, err := c.role.GetRoleIDs(&ddl.CustomRole{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: roleID,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// サイドバー一覧
	var roles []ddl.Role
	for _, row := range roleAssociate {
		roles = append(roles, ddl.Role{
			AbstractMasterModel: ddl.AbstractMasterModel{
				ID: row.MasterRoleID,
			},
		})
	}
	sidebars, err := c.master.ListSidebar(roles, &ddl.LoginType{
		AbstractMasterModel: ddl.AbstractMasterModel{
			ID: uint(loginType),
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.Sidebar{
		Sidebars: sidebars,
	}, nil
}

// 使用可能ロール一覧
func (c *CommonService) Roles(req *request.Roles) (*response.Roles, *response.Error) {
	// バリデーション
	if err := c.v.Roles(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ロールID、ログイン種別取得
	ctx := context.Background()
	role, err := c.redis.Get(ctx, req.HashKey, static.REDIS_USER_ROLE)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}
	roleID, err := strconv.ParseUint(*role, 10, 64)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	loginType_0, err := c.redis.Get(ctx, req.HashKey, static.REDIS_USER_LOGIN_TYPE)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}
	loginType, err := strconv.ParseUint(*loginType_0, 10, 64)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 該当ロールのマスタID取得
	roleAssociates, err := c.role.GetRoleIDs(&ddl.CustomRole{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: roleID,
		},
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ロール一覧
	roles, err := c.master.ListRole(&ddl.Role{
		RoleType: uint(loginType),
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// レスポンス
	m := make(map[string]bool)
	for _, row := range roles {
		found := false
		for _, row2 := range roleAssociates {
			if row2.MasterRoleID == row.ID {
				found = true
			}
		}
		m[row.NameEn] = found
	}

	return &response.Roles{
		Map: m,
	}, nil
}

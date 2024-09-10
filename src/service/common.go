package service

import (
	"api/src/model/ddl"
	"api/src/model/entity"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/repository"
	"api/src/validator"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ICommonService interface {
	// サイドバー表示
	Sidebar(req *request.Sidebar) (*response.Sidebar, *response.Error)
	// 使用可能ロール一覧
	Roles(req *request.Roles) (*response.Roles, *response.Error)
	// 所属チーム一覧
	Teams(req *request.TeamsBelong) (*response.TeamsBelong, *response.Error)
	// チーム変更
	ChangeTeam(req *request.ChangeTeam) *response.Error
}

type CommonService struct {
	master repository.IMasterRepository
	role   repository.IRoleRepository
	user   repository.IUserRepository
	team   repository.ITeamRepository
	v      validator.ICommonValidator
	redis  repository.IRedisRepository
}

func NewCommonService(
	master repository.IMasterRepository,
	role repository.IRoleRepository,
	user repository.IUserRepository,
	team repository.ITeamRepository,
	v validator.ICommonValidator,
	redis repository.IRedisRepository,
) ICommonService {
	return &CommonService{master, role, user, team, v, redis}
}

// サイドバー表示
func (c *CommonService) Sidebar(req *request.Sidebar) (*response.Sidebar, *response.Error) {
	// バリデーション
	if err := c.v.Sidebar(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// ロールID、ログイン種別取得
	ctx := context.Background()
	role, roleErr := c.redis.Get(ctx, req.HashKey, static.REDIS_USER_ROLE)
	if roleErr != nil {
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}
	roleID, roleParseErr := strconv.ParseUint(*role, 10, 64)
	if roleParseErr != nil {
		log.Printf("%v", roleParseErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	loginType_0, loginTypeErr := c.redis.Get(ctx, req.HashKey, static.REDIS_USER_LOGIN_TYPE)
	if loginTypeErr != nil {
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}
	loginType, loginTypeParseErr := strconv.ParseUint(*loginType_0, 10, 64)
	if loginTypeParseErr != nil {
		log.Printf("%v", loginTypeParseErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 該当ロールのマスタID取得
	roleAssociate, rolesErr := c.role.GetRoleIDs(&ddl.CustomRole{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: roleID,
		},
	})
	if rolesErr != nil {
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
	sidebars, sidebarsErr := c.master.ListSidebar(roles, &ddl.LoginType{
		AbstractMasterModel: ddl.AbstractMasterModel{
			ID: uint(loginType),
		},
	})
	if sidebarsErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var res []entity.Sidebar
	for _, row := range sidebars {
		flg := false
		for _, row2 := range res {
			if row.NameJa == row2.NameJa {
				flg = true
			}
		}

		if !flg {
			res = append(res, row)
		}
	}

	return &response.Sidebar{
		Sidebars: res,
	}, nil
}

// 使用可能ロール一覧
func (c *CommonService) Roles(req *request.Roles) (*response.Roles, *response.Error) {
	// バリデーション
	if err := c.v.Roles(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// ロールID、ログイン種別取得
	ctx := context.Background()
	role, roleErr := c.redis.Get(ctx, req.HashKey, static.REDIS_USER_ROLE)
	if roleErr != nil {
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}
	roleID, roleParseErr := strconv.ParseUint(*role, 10, 64)
	if roleParseErr != nil {
		log.Printf("%v", roleParseErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	loginType_0, loginTypeErr := c.redis.Get(ctx, req.HashKey, static.REDIS_USER_LOGIN_TYPE)
	if loginTypeErr != nil {
		return nil, &response.Error{
			Status: http.StatusUnauthorized,
			Code:   static.CODE_LOGIN_REQUIRED,
		}
	}
	loginType, loginTypeParseErr := strconv.ParseUint(*loginType_0, 10, 64)
	if loginTypeParseErr != nil {
		log.Printf("%v", loginTypeParseErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 該当ロールのマスタID取得
	roleAssociates, rolesErr := c.role.GetRoleIDs(&ddl.CustomRole{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			ID: roleID,
		},
	})
	if rolesErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ロール一覧
	roles, masterRoleErr := c.master.ListRole(&ddl.Role{
		RoleType: uint(loginType),
	})
	if masterRoleErr != nil {
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

// 所属チーム一覧
func (c *CommonService) Teams(req *request.TeamsBelong) (*response.TeamsBelong, *response.Error) {
	// バリデーション
	if err := c.v.Teams(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// ユーザーID取得
	user, userErr := c.user.Get(&ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if userErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 所属チーム一覧
	teams, teamsErr := c.team.ListBelongTeam(&ddl.TeamAssociation{
		UserID: user.ID,
	})
	if teamsErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.TeamsBelong{List: teams}, nil
}

// チーム変更
func (c *CommonService) ChangeTeam(req *request.ChangeTeam) *response.Error {
	// バリデーション
	if err := c.v.ChangeTeam(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// チーム取得
	team, teamErr := c.team.Get(&ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if teamErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	ctx := context.Background()
	teamID := strconv.FormatUint(team.ID, 10)
	if err := c.redis.Set(
		ctx,
		req.UserHashKey,
		static.REDIS_USER_TEAM_ID,
		&teamID,
		24*time.Hour,
	); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

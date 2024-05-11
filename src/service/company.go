package service

import (
	"api/resources/static"
	"api/src/model/ddl"
	"api/src/model/enum"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/repository"
	"api/src/validator"
	"log"
	"net/http"
)

type ICompanyService interface {
	// 登録
	Create(req *request.CompanyCreate) (*response.CompanyCreate, *response.Error)
}

type CompanyService struct {
	company repository.ICompanyRepository
	master  repository.IMasterRepository
	role    repository.IRoleRepository
	user    repository.IUserRepository
	v       validator.ICompanyValidator
	db      repository.IDBRepository
}

func NewCompanyService(
	company repository.ICompanyRepository,
	master repository.IMasterRepository,
	role repository.IRoleRepository,
	user repository.IUserRepository,
	v validator.ICompanyValidator,
	db repository.IDBRepository,
) ICompanyService {
	return &CompanyService{company, master, role, user, v, db}
}

// 登録
func (c *CompanyService) Create(req *request.CompanyCreate) (*response.CompanyCreate, *response.Error) {
	// バリデーション
	if err := c.v.Create(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// 企業名重複確認
	if err := c.company.IsDuplName(&req.Company); err != nil {
		return nil, &response.Error{
			Status: http.StatusConflict,
			Code:   static.CODE_COMPANY_NAME_DUPL,
		}
	}

	// ハッシュキー生成
	_, hash, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	req.HashKey = string(enum.PRE_COMPANY) + "_" + *hash

	tx, err := c.db.TxStart()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	company, err := c.company.Insert(tx, &req.Company)
	if err != nil {
		if err := c.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ロール用ハッシュキー生成
	_, roleHash, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ロール登録
	role, err := c.role.Insert(tx, &ddl.CustomRole{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(enum.PRE_ROLE) + "_" + *roleHash,
			CompanyID: uint(company.ID),
		},
		Name: "Initial role " + company.Name,
	})
	if err != nil {
		if err := c.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 付与ロール登録
	roles, err := c.master.ListRole(&ddl.Role{
		RoleType: uint(enum.LOGIN_TYPE_MANAGEMENT),
	})
	if err != nil {
		if err := c.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	for _, row := range roles {
		if err := c.role.InsertAssociation(tx, &ddl.RoleAssociation{
			RoleID:       uint(role.ID),
			MasterRoleID: row.ID,
		}); err != nil {
			if err := c.db.TxRollback(tx); err != nil {
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	// ユーザー用モデル生成
	_, userHash, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	password, hashPassword, err := GenerateHash(8, 16)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー登録
	userModel := ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(enum.PRE_USER) + "_" + *userHash,
			CompanyID: uint(company.ID),
		},
		Name:         "Initial user " + company.Name,
		Email:        req.Email,
		RoleID:       uint(role.ID),
		UserType:     enum.LOGIN_TYPE_MANAGEMENT,
		Password:     *hashPassword,
		InitPassword: *hashPassword,
	}
	user, err := c.user.Insert(tx, &userModel)
	if err != nil {
		if err := c.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム用モデル生成
	_, teamHash, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム登録
	team, err := c.user.InsertTeam(tx, &ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(enum.PRE_TEAM) + "_" + *teamHash,
			CompanyID: uint(company.ID),
		},
		Name: "Initial team " + company.Name,
	})
	if err != nil {
		if err := c.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム紐づけ登録
	if err := c.user.InsertTeamAssociation(tx, &ddl.TeamAssociation{
		TeamID: uint(team.ID),
		UserID: uint(user.ID),
	}); err != nil {
		if err := c.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := c.db.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.CompanyCreate{
		Password: *password,
	}, nil
}

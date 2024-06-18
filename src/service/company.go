package service

import (
	"api/src/model/ddl"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/repository"
	"api/src/validator"
	"log"
	"net/http"
)

type ICompanyService interface {
	// 登録
	Create(req *request.CreateCompany) (*response.CreateCompany, *response.Error)
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
func (c *CompanyService) Create(req *request.CreateCompany) (*response.CreateCompany, *response.Error) {
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
	_, hash, hashErr := GenerateHash(1, 25)
	if hashErr != nil {
		log.Printf("%v", hashErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	req.HashKey = string(static.PRE_COMPANY) + "_" + *hash

	tx, txErr := c.db.TxStart()
	if txErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	company, companyErr := c.company.Insert(tx, &req.Company)
	if companyErr != nil {
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
	_, roleHash, roleHashErr := GenerateHash(1, 25)
	if roleHashErr != nil {
		log.Printf("%v", roleHashErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ロール登録
	role, roleErr := c.role.Insert(tx, &ddl.CustomRole{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(static.PRE_ROLE) + "_" + *roleHash,
			CompanyID: company.ID,
		},
		Name: "Initial role " + company.Name,
	})
	if roleErr != nil {
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
	roles, rolesErr := c.master.ListRole(&ddl.Role{
		RoleType: uint(static.LOGIN_TYPE_MANAGEMENT),
	})
	if rolesErr != nil {
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
			RoleID:       role.ID,
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
	_, userHash, userHashErr := GenerateHash(1, 25)
	if userHashErr != nil {
		log.Printf("%v", userHashErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	password, hashPassword, passwordErr := GenerateHash(8, 16)
	if passwordErr != nil {
		log.Printf("%v", passwordErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー登録
	userModel := ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(static.PRE_USER) + "_" + *userHash,
			CompanyID: company.ID,
		},
		Name:         "Initial user " + company.Name,
		Email:        req.Email,
		RoleID:       role.ID,
		UserType:     static.LOGIN_TYPE_MANAGEMENT,
		Password:     *hashPassword,
		InitPassword: *hashPassword,
	}
	user, userErr := c.user.Insert(tx, &userModel)
	if userErr != nil {
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
	_, teamHash, teamHashErr := GenerateHash(1, 25)
	if teamHashErr != nil {
		log.Printf("%v", teamHashErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム登録
	team, teamErr := c.user.InsertTeam(tx, &ddl.Team{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(static.PRE_TEAM) + "_" + *teamHash,
			CompanyID: company.ID,
		},
		Name: "Initial team " + company.Name,
	})
	if teamErr != nil {
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
		TeamID: team.ID,
		UserID: user.ID,
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

	// 選考状況用モデル作成
	_, selectHash, selectHashErr := GenerateHash(1, 25)
	if selectHashErr != nil {
		log.Printf("%v", selectHashErr)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	_, selectHash2, selectHash2Err := GenerateHash(1, 25)
	if selectHash2Err != nil {
		log.Printf("%v", selectHash2Err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 選考状況登録
	if err := c.user.InsertSelectStatus(tx, &ddl.SelectStatus{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(static.PRE_SELECT_STATUS) + "_" + *selectHash,
			CompanyID: company.ID,
		},
		TeamID:     team.ID,
		StatusName: "日程未回答",
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
	if err := c.user.InsertSelectStatus(tx, &ddl.SelectStatus{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   string(static.PRE_SELECT_STATUS) + "_" + *selectHash2,
			CompanyID: company.ID,
		},
		TeamID:     team.ID,
		StatusName: "日程回答済み",
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

	return &response.CreateCompany{
		Password: *password,
	}, nil
}

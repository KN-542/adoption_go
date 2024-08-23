package service

import (
	"api/src/model/ddl"
	"api/src/model/request"
	"api/src/model/response"
	"api/src/model/static"
	"api/src/repository"
	"api/src/validator"
	"fmt"
	"log"
	"net/http"
)

type ICompanyService interface {
	// 登録
	Create(req *request.CreateCompany) (*response.CreateCompany, *response.Error)
	// 検索
	Search(req *request.SearchCompany) (*response.SearchCompany, *response.Error)
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

	var list []*ddl.RoleAssociation
	for _, row := range roles {
		list = append(list, &ddl.RoleAssociation{
			RoleID:       role.ID,
			MasterRoleID: row.ID,
		})
	}

	if err := c.role.InsertsAssociation(tx, list); err != nil {
		if err := c.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
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
		Name:           "Initial team " + company.Name,
		NumOfInterview: 3,
		RuleID:         static.ASSIGN_RULE_MANUAL,
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

	// 面接毎設定登録 & 面接毎参加可能者登録(全員参加可能)
	var perList []*ddl.TeamPerInterview
	var possibleList []*ddl.TeamAssignPossible
	for i := 1; i <= int(team.NumOfInterview); i++ {
		perList = append(perList, &ddl.TeamPerInterview{
			TeamID:         team.ID,
			NumOfInterview: uint(i),
			UserMin:        1,
		})
		possibleList = append(possibleList, &ddl.TeamAssignPossible{
			TeamID:         team.ID,
			UserID:         user.ID,
			NumOfInterview: uint(i),
		})
	}
	if err := c.user.InsertsPerInterview(tx, perList); err != nil {
		if err := c.db.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if err := c.user.InsertsAssignPossible(tx, possibleList); err != nil {
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

	fmt.Print(*password)

	return &response.CreateCompany{
		Password: *password,
	}, nil
}

// 検索
func (c *CompanyService) Search(req *request.SearchCompany) (*response.SearchCompany, *response.Error) {
	// バリデーション
	if err := c.v.Search(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 検索
	res, err := c.company.Search(&req.Company)
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.SearchCompany{
		List: res,
	}, nil
}

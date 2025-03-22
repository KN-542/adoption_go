package service

import (
	"api/src/model/ddl"
	"api/src/model/dto"
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
)

type IManuscriptService interface {
	// 検索
	Search(req *request.SearchManuscript) (*response.SearchManuscript, *response.Error)
	// 登録
	Create(req *request.CreateManuscript) *response.Error
	// 応募者紐づけ登録
	CreateApplicantAssociation(req *request.CreateApplicantAssociation) *response.Error
	// 検索_同一チーム
	SearchManuscriptByTeam(req *request.SearchManuscriptByTeam) (*response.SearchManuscriptByTeam, *response.Error)
	// 削除
	Delete(req *request.DeleteManuscriptRequest) *response.Error
	// 取得
	Get(req *request.GetManuscript) (*response.GetManuscript, *response.Error)
	// 更新
	Update(req *request.UpdateManuscript) *response.Error
}

type ManuscriptService struct {
	manuscript repository.IManuscriptRepository
	master     repository.IMasterRepository
	user       repository.IUserRepository
	team       repository.ITeamRepository
	applicant  repository.IApplicantRepository
	db         repository.IDBRepository
	redis      repository.IRedisRepository
	validate   validator.IManuscriptValidator
}

func NewManuscriptService(
	manuscript repository.IManuscriptRepository,
	master repository.IMasterRepository,
	user repository.IUserRepository,
	team repository.ITeamRepository,
	applicant repository.IApplicantRepository,
	db repository.IDBRepository,
	redis repository.IRedisRepository,
	validate validator.IManuscriptValidator,
) IManuscriptService {
	return &ManuscriptService{manuscript, master, user, team, applicant, db, redis, validate}
}

// 検索
func (s *ManuscriptService) Search(req *request.SearchManuscript) (*response.SearchManuscript, *response.Error) {
	// バリデーション
	if err := s.validate.Search(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 企業ID取得
	ctx := context.Background()
	company, companyErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
	if companyErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	companyID, companyIDErr := strconv.ParseUint(*company, 10, 64)
	if companyIDErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 検索
	dto := dto.SearchManuscript{
		SearchManuscript: *req,
		CompanyID:        companyID,
	}
	manuscripts, count, manuscriptsErr := s.manuscript.Search(&dto)
	if manuscriptsErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var res []entity.SearchManuscript
	for _, manuscript := range manuscripts {
		manuscript.ID = 0
		for _, site := range manuscript.Sites {
			site.ID = 0
		}
		res = append(res, *manuscript)
	}

	return &response.SearchManuscript{
		List: res,
		Num:  uint64(count),
	}, nil
}

// 登録
func (s *ManuscriptService) Create(req *request.CreateManuscript) *response.Error {
	// バリデーション
	if err := s.validate.Create(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 企業ID取得
	ctx := context.Background()
	company, companyErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
	if companyErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	companyID, companyIDErr := strconv.ParseUint(*company, 10, 64)
	if companyIDErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 内容重複チェック
	count, countErr := s.manuscript.CheckDuplicateContent(&ddl.Manuscript{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			CompanyID: companyID,
		},
		Content: req.Content,
	})
	if countErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	if *count > 0 {
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_MANUSCRIPT_DUPLICATE_CONTENT,
		}
	}

	// チームID取得
	teamIDs, teamIDsErr := s.team.GetIDs(req.Teams)
	if teamIDsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// サイトID取得
	siteIDs, siteIDsErr := s.master.SelectSiteIDs(req.Sites)
	if siteIDsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ハッシュキー生成
	_, hashKey, hashErr := GenerateHash(1, 25)
	if hashErr != nil {
		log.Printf("%v", hashErr)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, txErr := s.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	manuscript, manuscriptErr := s.manuscript.Insert(tx, &ddl.Manuscript{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   static.PRE_MANUSCRIPT + "_" + *hashKey,
			CompanyID: companyID,
		},
		Content: req.Content,
	})
	if manuscriptErr != nil {
		if err := s.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var teamAssociations []*ddl.ManuscriptTeamAssociation
	for _, teamID := range teamIDs {
		teamAssociations = append(teamAssociations, &ddl.ManuscriptTeamAssociation{
			ManuscriptID: manuscript.ID,
			TeamID:       teamID,
		})
	}

	var siteAssociations []*ddl.ManuscriptSiteAssociation
	for _, siteID := range siteIDs {
		siteAssociations = append(siteAssociations, &ddl.ManuscriptSiteAssociation{
			ManuscriptID: manuscript.ID,
			SiteID:       siteID,
		})
	}

	// チーム紐づけ登録
	if err := s.manuscript.InsertTeamAssociation(tx, teamAssociations); err != nil {
		if err := s.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// サイト紐づけ登録
	if err := s.manuscript.InsertSiteAssociation(tx, siteAssociations); err != nil {
		if err := s.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := s.db.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 応募者紐づけ登録
func (s *ManuscriptService) CreateApplicantAssociation(req *request.CreateApplicantAssociation) *response.Error {
	// バリデーション
	if err := s.validate.CreateApplicantAssociation(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 原稿取得
	manuscript, manuscriptErr := s.manuscript.Get(&ddl.Manuscript{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.ManuscriptHash,
		},
	})
	if manuscriptErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 応募者ID取得
	ids, idsErr := s.applicant.GetIDs(req.Applicants)
	if idsErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	var associations []*ddl.ManuscriptApplicantAssociation
	for _, id := range ids {
		associations = append(associations, &ddl.ManuscriptApplicantAssociation{
			ManuscriptID: manuscript.ID,
			ApplicantID:  id,
		})
	}

	tx, txErr := s.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 削除
	if err := s.manuscript.DeleteApplicantAssociation(tx, ids); err != nil {
		if err := s.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	if err := s.manuscript.InsertsApplicantAssociation(tx, associations); err != nil {
		if err := s.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := s.db.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 削除処理
func (s *ManuscriptService) Delete(req *request.DeleteManuscriptRequest) *response.Error {
	// 原稿ID取得
	manuscriptIDs, manuscriptErr := s.manuscript.GetManuscriptIDsByHashKeys(req.ManuscriptHashKeys)
	if manuscriptErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 応募者に紐づいている原稿IDがあるかチェック
	count, countErr := s.manuscript.CheckManuscriptAssociationByApplicant(manuscriptIDs)
	if countErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	//  応募者に紐づいている原稿IDがあれば、削除不可(400 エラーを返す)
	if count > 0 {
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_MANUSCRIPT_CANNOT_DELETE_APPLICANT,
		}
	}

	// トランザクションの開始
	tx, txErr := s.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム紐づけの削除
	if err := s.manuscript.DeleteTeamAssociation(tx, manuscriptIDs); err != nil {
		if err := s.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// サイト紐づけの削除
	if err := s.manuscript.DeleteSiteAssociation(tx, manuscriptIDs); err != nil {
		if err := s.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 原稿の削除
	if err := s.manuscript.Delete(tx, req.ManuscriptHashKeys); err != nil {
		if err := s.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// トランザクションのコミット
	if err := s.db.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 検索_同一チーム
func (s *ManuscriptService) SearchManuscriptByTeam(req *request.SearchManuscriptByTeam) (*response.SearchManuscriptByTeam, *response.Error) {
	// チームID取得
	ctx := context.Background()
	team, teamErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_TEAM_ID)
	if teamErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	teamID, teamIDErr := strconv.ParseUint(*team, 10, 64)
	if teamIDErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 検索
	res, err := s.manuscript.SearchByTeam(&ddl.ManuscriptTeamAssociation{
		TeamID: teamID,
	})
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &response.SearchManuscriptByTeam{
		List: res,
	}, nil
}

// 取得
func (s *ManuscriptService) Get(req *request.GetManuscript) (*response.GetManuscript, *response.Error) {
	// バリデーション
	if err := s.validate.Get(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
		}
	}
	// 原稿取得
	manuscript, manuscriptErr := s.manuscript.Get(&ddl.Manuscript{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if manuscriptErr != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	//　ID隠蔽処理
	for _, site := range manuscript.Sites {
		site.ID = 0
	}
	for _, team := range manuscript.Teams {
		team.ID = 0
	}

	return &response.GetManuscript{
		Manuscript: entity.Manuscript{
			Manuscript: ddl.Manuscript{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: manuscript.HashKey,
				},
				Content: manuscript.Content,
			},
			Teams: manuscript.Teams,
			Sites: manuscript.Sites,
		},
	}, nil
}

// 更新
func (s *ManuscriptService) Update(req *request.UpdateManuscript) *response.Error {
	// バリデーション
	if err := s.validate.Update(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
		}
	}

	// 原稿取得
	manuscript, manuscriptErr := s.manuscript.Get(&ddl.Manuscript{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
	})
	if manuscriptErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 内容が更新されている場合、他のものと重複していないかチェック
	if req.Content != manuscript.Content {
		// 企業ID取得
		ctx := context.Background()
		company, companyErr := s.redis.Get(ctx, req.UserHashKey, static.REDIS_USER_COMPANY_ID)
		if companyErr != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		companyID, companyIDErr := strconv.ParseUint(*company, 10, 64)
		if companyIDErr != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}

		// 内容重複チェック
		count, countErr := s.manuscript.CheckDuplicateContent(&ddl.Manuscript{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				CompanyID: companyID,
			},
			Content: req.Content,
		})
		if countErr != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		if *count > 0 {
			return &response.Error{
				Status: http.StatusBadRequest,
				Code:   static.CODE_MANUSCRIPT_DUPLICATE_CONTENT,
			}
		}
	}

	// チームの紐付けが新たにリクエストされているか
	isUpdateTeams := len(req.Teams) > 0
	// チームの紐付けが新たにリクエストされている場合、チームID取得
	var teamIDs []uint64 = nil
	if isUpdateTeams {
		var err error
		teamIDs, err = s.team.GetIDs(req.Teams)
		if err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	// サイトの紐付けが新たにリクエストされているか
	isUpdateSites := len(req.Sites) > 0
	// サイトの紐付けが新たにリクエストされている場合、サイトID取得
	var siteIDs []uint = nil
	if isUpdateSites {
		var err error
		siteIDs, err = s.master.SelectSiteIDs(req.Sites)
		if err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	// トランザクションの開始
	tx, txErr := s.db.TxStart()
	if txErr != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 更新処理
	manuscriptUpdateErr := s.manuscript.Update(tx, &ddl.Manuscript{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: req.HashKey,
		},
		Content: req.Content,
	})
	if manuscriptUpdateErr != nil {
		if err := s.db.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チームの紐付けが新たにリクエストされている場合、チーム紐付けの新規割当て
	if isUpdateTeams {
		//　チーム紐付け情報作成
		var teamAssociations []*ddl.ManuscriptTeamAssociation
		for _, teamID := range teamIDs {
			teamAssociations = append(teamAssociations, &ddl.ManuscriptTeamAssociation{
				ManuscriptID: manuscript.ID,
				TeamID:       teamID,
			})
		}
		// DB更新
		if err := s.manuscript.InsertTeamAssociation(tx, teamAssociations); err != nil {
			if err := s.db.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	// サイトの紐付けが新たにリクエストされている場合、サイト紐付けの新規割当て
	if isUpdateSites {
		// サイト紐付け情報作成
		var siteAssociations []*ddl.ManuscriptSiteAssociation
		for _, siteID := range siteIDs {
			siteAssociations = append(siteAssociations, &ddl.ManuscriptSiteAssociation{
				ManuscriptID: manuscript.ID,
				SiteID:       siteID,
			})
		}
		// DB更新
		if err := s.manuscript.InsertSiteAssociation(tx, siteAssociations); err != nil {
			if err := s.db.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	// トランザクションのコミット
	if err := s.db.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

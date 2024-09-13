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
	// 検索_同一チーム
	SearchManuscriptByTeam(req *request.SearchManuscriptByTeam) (*response.SearchManuscriptByTeam, *response.Error)
}

type ManuscriptService struct {
	manuscript repository.IManuscriptRepository
	master     repository.IMasterRepository
	user       repository.IUserRepository
	team       repository.ITeamRepository
	db         repository.IDBRepository
	redis      repository.IRedisRepository
	validate   validator.IManuscriptValidator
}

func NewManuscriptService(
	manuscript repository.IManuscriptRepository,
	master repository.IMasterRepository,
	user repository.IUserRepository,
	team repository.ITeamRepository,
	db repository.IDBRepository,
	redis repository.IRedisRepository,
	validate validator.IManuscriptValidator,
) IManuscriptService {
	return &ManuscriptService{manuscript, master, user, team, db, redis, validate}
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
	dto := dto.SearchManuscript{
		SearchManuscript: *req,
		TeamID:           teamID,
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

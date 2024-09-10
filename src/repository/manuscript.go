package repository

import (
	"api/src/model/ddl"
	"api/src/model/dto"
	"api/src/model/entity"
	"log"

	"gorm.io/gorm"
)

type IManuscriptRepository interface {
	// 登録
	Insert(tx *gorm.DB, m *ddl.Manuscript) (*entity.Manuscript, error)
	// チーム紐づけ登録
	InsertTeamAssociation(tx *gorm.DB, m []*ddl.ManuscriptTeamAssociation) error
	// サイト紐づけ登録
	InsertSiteAssociation(tx *gorm.DB, m []*ddl.ManuscriptSiteAssociation) error
	// 検索
	Search(m *dto.SearchManuscript) ([]*entity.SearchManuscript, int64, error)
	// 紐づけ取得
	GetAssociationByTeamID(m *ddl.ManuscriptTeamAssociation) ([]entity.ManuscriptTeamAssociation, error)
	// 内容重複チェック
	CheckDuplicateContent(m *ddl.Manuscript) (*int64, error)
}

type ManuscriptRepository struct {
	db *gorm.DB
}

func NewManuscriptRepository(db *gorm.DB) IManuscriptRepository {
	return &ManuscriptRepository{db}
}

// 登録
func (u *ManuscriptRepository) Insert(tx *gorm.DB, m *ddl.Manuscript) (*entity.Manuscript, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &entity.Manuscript{
		Manuscript: *m,
	}, nil
}

// チーム紐づけ登録
func (u *ManuscriptRepository) InsertTeamAssociation(tx *gorm.DB, m []*ddl.ManuscriptTeamAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// サイト紐づけ登録
func (u *ManuscriptRepository) InsertSiteAssociation(tx *gorm.DB, m []*ddl.ManuscriptSiteAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 検索
func (s *ManuscriptRepository) Search(m *dto.SearchManuscript) ([]*entity.SearchManuscript, int64, error) {
	var res []*entity.SearchManuscript
	var count int64

	query := s.db.Table("t_manuscript").
		Joins(`
			LEFT JOIN
				t_manuscript_team_association
			ON
				t_manuscript_team_association.manuscript_id = t_manuscript.id
		`).
		Where("t_manuscript_team_association.team_id = ?", m.TeamID)

	if len(m.Sites) > 0 {
		query = query.Joins(`
				INNER JOIN
					t_manuscript_site_association
				ON
					t_manuscript_site_association.manuscript_id = t_manuscript.id
			`).
			Joins(`
				INNER JOIN
					m_site
				ON
					t_manuscript_site_association.site_id = m_site.id
			`).
			Where("m_site.hash_key IN ?", m.Sites)
	}

	if err := query.Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return nil, 0, err
	}

	offset := (m.Page - 1) * m.PageSize
	if err := query.Select(`
			t_manuscript.id,
			t_manuscript.hash_key,
			t_manuscript.content
		`).
		Offset(offset).
		Limit(m.PageSize).
		Find(&res).
		Error; err != nil {
		log.Printf("%v", err)
		return nil, 0, err
	}

	if len(res) > 0 {
		var siteAssociations []struct {
			ManuscriptID uint64
			SiteID       uint
			HashKey      string
			SiteName     string
		}
		var manuscriptIDs []uint64
		for _, manuscript := range res {
			manuscriptIDs = append(manuscriptIDs, manuscript.ID)
		}

		if err := s.db.Table("t_manuscript_site_association").
			Select(`
			t_manuscript_site_association.manuscript_id,
			m_site.id AS site_id,
			m_site.hash_key,
			m_site.site_name
		`).
			Joins(`
			INNER JOIN
				m_site
			ON
				t_manuscript_site_association.site_id = m_site.id
		`).
			Where("t_manuscript_site_association.manuscript_id IN ?", manuscriptIDs).
			Find(&siteAssociations).Error; err != nil {
			log.Printf("%v", err)
			return nil, 0, err
		}

		siteMap := make(map[uint64][]ddl.Site)
		for _, assoc := range siteAssociations {
			site := ddl.Site{
				AbstractMasterModel: ddl.AbstractMasterModel{
					ID:      assoc.SiteID,
					HashKey: assoc.HashKey,
				},
				SiteName: assoc.SiteName,
			}
			siteMap[assoc.ManuscriptID] = append(siteMap[assoc.ManuscriptID], site)
		}

		for _, r := range res {
			r.Sites = siteMap[r.ID]
		}
	}

	return res, count, nil
}

// 紐づけ取得
func (s *ManuscriptRepository) GetAssociationByTeamID(m *ddl.ManuscriptTeamAssociation) ([]entity.ManuscriptTeamAssociation, error) {
	var res []entity.ManuscriptTeamAssociation

	query := s.db.Model(&ddl.ManuscriptTeamAssociation{}).Where(
		&ddl.ManuscriptTeamAssociation{
			TeamID: m.TeamID,
		},
	)

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return res, nil
}

// 内容重複チェック
func (s *ManuscriptRepository) CheckDuplicateContent(m *ddl.Manuscript) (*int64, error) {
	var count int64

	if err := s.db.Model(&ddl.Manuscript{}).Where(
		&ddl.Manuscript{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				CompanyID: m.CompanyID,
			},
			Content: m.Content,
		},
	).Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &count, nil
}

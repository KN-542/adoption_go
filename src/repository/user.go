package repository

import (
	"api/src/model/ddl"
	"api/src/model/dto"
	"api/src/model/entity"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type IUserRepository interface {
	// ログイン認証
	Login(m *ddl.User) ([]entity.Login, error)
	// 登録
	Insert(tx *gorm.DB, m *ddl.User) (*entity.User, error)
	// 検索
	Search(m *dto.SearchUser) ([]entity.SearchUser, error)
	// 検索_同一企業
	SearchByCompany(m *ddl.User) ([]entity.SearchUser, error)
	// 取得
	Get(m *ddl.User) (*entity.User, error)
	// 取得_PK
	GetByPrimary(m *ddl.User) (*entity.User, error)
	// 更新
	Update(tx *gorm.DB, m *ddl.User) error
	// 削除
	Delete(tx *gorm.DB, m []string) error
	// リフレッシュトークン紐づけ登録
	InsertUserRefreshTokenAssociation(tx *gorm.DB, m *ddl.UserRefreshTokenAssociation) error
	// リフレッシュトークン紐づけ取得
	GetUserRefreshTokenAssociation(m *ddl.UserRefreshTokenAssociation) (*entity.UserRefreshTokenAssociation, error)
	// リフレッシュトークン紐づけ取得_ハッシュキー
	GetUserRefreshTokenAssociationByHashKey(m *ddl.User) ([]entity.UserRefreshTokenAssociation, error)
	// メールアドレス重複チェック
	EmailDuplCheck(m *ddl.User) error
	// ID取得
	GetIDs(m []string) ([]uint64, error)
	// 取得_ハッシュキー配列
	GetByHashKeys(m []string) ([]entity.User, error)
	// ユーザー取得_予定数順
	GetUsersSortedByScheduleCount(m *ddl.Schedule) ([]entity.User, error)
	// ユーザーと紐づいている応募者数を取得
	CountApplicantUserAssociation(m []uint64) (int64, error)
	// ユーザーと紐づいているスケジュール数を取得
	CountScheduleAssociation(m []uint64) (int64, error)
	// 削除_通知
	DeleteNotice(tx *gorm.DB, m []uint64) error
	// 削除_面接毎参加可能者
	DeleteTeamAssignPossible(tx *gorm.DB, m []uint64) error
	// 削除_面接割り振り優先順位
	DeleteTeamAssignPriority(tx *gorm.DB, m []uint64) error
	// 削除_チーム紐づけ
	DeleteTeamAssociation(tx *gorm.DB, m []uint64) error
	// 削除_リフレッシュトークン紐づけ
	DeleteUserRefreshTokenAssociation(tx *gorm.DB, m []uint64) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

// ログイン認証
func (u *UserRepository) Login(m *ddl.User) ([]entity.Login, error) {
	var res []entity.Login

	if err := u.db.Model(&entity.Login{}).
		Select(`
			t_user.hash_key,
			t_user.name,
			t_user.company_id,
			t_user.password,
			t_user.init_password,
			t_user.role_id,
			t_user.user_type
		`).
		Where(
			&ddl.User{
				Email: m.Email,
			},
		).Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// 登録
func (u *UserRepository) Insert(tx *gorm.DB, m *ddl.User) (*entity.User, error) {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return &entity.User{
		User: *m,
	}, nil
}

// 検索
func (u *UserRepository) Search(m *dto.SearchUser) ([]entity.SearchUser, error) {
	var l []entity.SearchUser

	// TODO
	query := u.db.Model(&entity.SearchUser{}).
		Select(`
			t_user.hash_key,
			t_user.name,
			t_user.email,
			t_role.name as role_name
		`).
		Joins("LEFT JOIN t_role ON t_role.id = t_user.role_id").
		Where("t_user.company_id = ?", m.CompanyID)

	if err := query.Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// 検索_同一企業
func (u *UserRepository) SearchByCompany(m *ddl.User) ([]entity.SearchUser, error) {
	var l []entity.SearchUser

	if err := u.db.Model(&entity.SearchUser{}).
		Select("hash_key, name, email").
		Where("company_id = ?", m.CompanyID).
		Find(&l).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return l, nil
}

// 取得
func (u *UserRepository) Get(m *ddl.User) (*entity.User, error) {
	var res entity.User
	if err := u.db.Where(
		&ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// 取得_PK
func (u *UserRepository) GetByPrimary(m *ddl.User) (*entity.User, error) {
	var res entity.User
	if err := u.db.Where(
		&ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				ID: m.ID,
			},
		},
	).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// 更新
func (u *UserRepository) Update(tx *gorm.DB, m *ddl.User) error {
	user := ddl.User{
		Name:     m.Name,
		Email:    m.Email,
		Password: m.Password,
		RoleID:   m.RoleID,
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			UpdatedAt: time.Now(),
		},
	}
	if err := tx.Model(&ddl.User{}).Where(
		&ddl.User{
			AbstractTransactionModel: ddl.AbstractTransactionModel{
				HashKey: m.HashKey,
			},
		},
	).Updates(user).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}

// 削除
func (u *UserRepository) Delete(tx *gorm.DB, m []string) error {
	if err := tx.
		Where("hash_key IN ?", m).
		Delete(&ddl.User{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// リフレッシュトークン紐づけ登録
func (u *UserRepository) InsertUserRefreshTokenAssociation(tx *gorm.DB, m *ddl.UserRefreshTokenAssociation) error {
	if err := tx.Create(m).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// リフレッシュトークン紐づけ取得
func (u *UserRepository) GetUserRefreshTokenAssociation(m *ddl.UserRefreshTokenAssociation) (*entity.UserRefreshTokenAssociation, error) {
	var res entity.UserRefreshTokenAssociation

	if err := u.db.Model(&ddl.UserRefreshTokenAssociation{}).Where(
		&ddl.UserRefreshTokenAssociation{
			UserID: m.UserID,
		},
	).First(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &res, nil
}

// リフレッシュトークン紐づけ取得_ハッシュキー
func (u *UserRepository) GetUserRefreshTokenAssociationByHashKey(m *ddl.User) ([]entity.UserRefreshTokenAssociation, error) {
	var res []entity.UserRefreshTokenAssociation

	if err := u.db.Model(&ddl.UserRefreshTokenAssociation{}).
		Joins("LEFT JOIN t_user ON t_user.id = t_user_refresh_token_association.user_id").
		Where("t_user.hash_key = ?", m.HashKey).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// メールアドレス重複チェック
func (u *UserRepository) EmailDuplCheck(m *ddl.User) error {
	var count int64
	if err := u.db.Model(&ddl.User{}).Where(
		&ddl.User{
			Email: m.Email,
		},
	).Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("duplicate Email Address")
	}

	return nil
}

// ID取得
func (u *UserRepository) GetIDs(m []string) ([]uint64, error) {
	var res []entity.User
	if err := u.db.Table("t_user").
		Select("id").
		Where("hash_key IN ?", m).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	var IDs []uint64
	for _, row := range res {
		IDs = append(IDs, row.ID)
	}

	return IDs, nil
}

// 取得_ハッシュキー配列
func (u *UserRepository) GetByHashKeys(m []string) ([]entity.User, error) {
	var res []entity.User
	if err := u.db.Model(&ddl.User{}).
		Select("id, hash_key, name, email").
		Where("hash_key IN ?", m).
		Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// ユーザー取得_予定数順
func (u *UserRepository) GetUsersSortedByScheduleCount(m *ddl.Schedule) ([]entity.User, error) {
	var res []entity.User

	query := u.db.Table("t_user").
		Joins("LEFT JOIN t_team_association ON t_team_association.user_id = t_user.id").
		Joins("LEFT JOIN t_schedule_association ON t_schedule_association.user_id = t_user.id").
		Where("t_team_association.team_id = ?", m.TeamID).
		Group("t_user.id").
		Order("COUNT(DISTINCT t_schedule_association.schedule_id) ASC")

	if err := query.Find(&res).Error; err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return res, nil
}

// ユーザーと紐づいている応募者数を取得
func (u *UserRepository) CountApplicantUserAssociation(m []uint64) (int64, error) {
	var count int64
	if err := u.db.Model(&ddl.ApplicantUserAssociation{}).
		Where("user_id IN ?", m).
		Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return 0, err
	}
	return count, nil
}

// ユーザーと紐づいているスケジュール数を取得
func (u *UserRepository) CountScheduleAssociation(m []uint64) (int64, error) {
	var count int64
	if err := u.db.Model(&ddl.ScheduleAssociation{}).
		Where("user_id IN ?", m).
		Count(&count).Error; err != nil {
		log.Printf("%v", err)
		return 0, err
	}
	return count, nil
}

// 削除_通知
func (u *UserRepository) DeleteNotice(tx *gorm.DB, m []uint64) error {
	// 通知_通知元ユーザー
	if err := tx.
		Where("from_user_id IN ?", m).
		Delete(&ddl.Notice{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	// 通知_通知先ユーザー
	if err := tx.
		Where("to_user_id IN ?", m).
		Delete(&ddl.Notice{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 削除_面接毎参加可能者
func (u *UserRepository) DeleteTeamAssignPossible(tx *gorm.DB, m []uint64) error {
	if err := tx.
		Where("user_id IN ?", m).
		Delete(&ddl.TeamAssignPossible{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 削除_面接割り振り優先順位
func (u *UserRepository) DeleteTeamAssignPriority(tx *gorm.DB, m []uint64) error {
	if err := tx.
		Where("user_id IN ?", m).
		Delete(&ddl.TeamAssignPriority{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 削除_チーム紐づけ
func (u *UserRepository) DeleteTeamAssociation(tx *gorm.DB, m []uint64) error {
	if err := tx.
		Where("user_id IN ?", m).
		Delete(&ddl.TeamAssociation{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

// 削除_リフレッシュトークン紐づけ
func (u *UserRepository) DeleteUserRefreshTokenAssociation(tx *gorm.DB, m []uint64) error {
	if err := tx.
		Where("user_id IN ?", m).
		Delete(&ddl.UserRefreshTokenAssociation{}).Error; err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil
}

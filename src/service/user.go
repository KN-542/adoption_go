package service

import (
	"api/resources/static"
	"api/src/model"
	"api/src/repository"
	"api/src/validator"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	// 一覧
	List() (*model.UsersResponse, *model.ErrorResponse)
	// 登録
	Create(req *model.User) (*model.UserResponse, *model.ErrorResponse)
	// 取得
	Get(req *model.User) (*model.UserResponse, *model.ErrorResponse)
	// ロール一覧
	RoleList() (*model.UserRoles, *model.ErrorResponse)
	// 検索(グループ)
	SearchGroups() (*model.UserGroupsResponse, *model.ErrorResponse)
	// グループ登録
	CreateGroup(req *model.UserGroup) *model.ErrorResponse
	// スケジュール登録種別一覧
	ListScheduleType() (*model.CalendarsFreqStatus, *model.ErrorResponse)
	// スケジュール登録
	CreateSchedule(req *model.UserScheduleRequest) *model.ErrorResponse
	// スケジュール一覧
	Schedules() (*model.UserSchedulesResponse, *model.ErrorResponse)
	// スケジュール削除
	DeleteSchedule(req *model.UserSchedule) *model.ErrorResponse
	// 予約表提示
	DispReserveTable() (*model.ReserveTableResponse, *model.ErrorResponse)
}

type UserService struct {
	r repository.IUserRepository
	m repository.IMasterRepository
	v validator.IUserValidator
	d repository.IDBRepository
	o repository.IOuterIFRepository
}

func NewUserService(
	r repository.IUserRepository,
	m repository.IMasterRepository,
	v validator.IUserValidator,
	d repository.IDBRepository,
	o repository.IOuterIFRepository,
) IUserService {
	return &UserService{r, m, v, d, o}
}

// 一覧
func (u *UserService) List() (*model.UsersResponse, *model.ErrorResponse) {
	users, err := u.r.List()
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.UsersResponse{
		Users: users,
	}, nil
}

// 登録
func (u *UserService) Create(req *model.User) (*model.UserResponse, *model.ErrorResponse) {
	// バリデーション
	if err := u.v.CreateValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// メールアドレス重複チェック
	if err := u.r.EmailDuplCheck(req); err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusConflict,
			Code:   static.CODE_USER_EMAIL_DUPL,
		}
	}

	// 初回パスワード発行
	password, hashPassword, err := GenerateHash(8, 16)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// ハッシュキー生成
	_, hashKey, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	user := model.User{
		HashKey:      *hashKey,
		Name:         req.Name,
		Email:        req.Email,
		Password:     *hashPassword,
		InitPassword: *hashPassword,
		RoleID:       req.RoleID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.r.Insert(tx, &user); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return nil, &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	res := model.UserResponse{
		Email:        user.Email,
		InitPassword: *password,
	}
	return &res, nil
}

// 取得
func (u *UserService) Get(req *model.User) (*model.UserResponse, *model.ErrorResponse) {
	if err := u.v.HashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	res, err := u.r.Get(req)
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &*model.ConvertUser(res), nil
}

// ロール一覧
func (u *UserService) RoleList() (*model.UserRoles, *model.ErrorResponse) {
	roles, err := u.m.SelectRole()
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.UserRoles{Roles: *roles}, nil
}

// 検索(グループ)
func (u *UserService) SearchGroups() (*model.UserGroupsResponse, *model.ErrorResponse) {
	userGroups, err := u.r.SearchGroup()
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー存在確認
	for index, userGroup := range userGroups {
		var l []string
		users, err := u.r.ConfirmUserByHashKeys(strings.Split(userGroup.Users, ","))
		if err != nil {
			if err != nil {
				return nil, &model.ErrorResponse{
					Status: http.StatusInternalServerError,
				}
			}
		}
		if users == nil || len(users) == 0 {
			userGroups[index].Users = ""
			continue
		}

		for _, user := range users {
			l = append(l, user.Name)
		}
		userGroups[index].Users = strings.Join(l, ",")
	}

	return &model.UserGroupsResponse{
		UserGroups: userGroups,
	}, nil
}

// グループ登録
func (u *UserService) CreateGroup(req *model.UserGroup) *model.ErrorResponse {
	// バリデーション
	if err := u.v.CreateGroupValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ユーザー存在確認
	users, err := u.r.ConfirmUserByHashKeys(strings.Split(req.Users, ","))
	if err != nil {
		if err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
	}

	var l []string
	for _, row := range users {
		l = append(l, row.HashKey)
	}
	req.Users = strings.Join(l, ",")

	// ハッシュキー生成
	_, hashKey, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// グループ登録
	req.HashKey = *hashKey
	if err := u.r.InsertGroup(tx, req); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
		if err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// スケジュール登録種別一覧
func (u *UserService) ListScheduleType() (*model.CalendarsFreqStatus, *model.ErrorResponse) {
	res, err := u.m.SelectCalendarFreqStatus()
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.CalendarsFreqStatus{
		List: *res,
	}, nil
}

// スケジュール登録
func (u *UserService) CreateSchedule(req *model.UserScheduleRequest) *model.ErrorResponse {
	// バリデーション
	if err := u.v.CreateScheduleValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ハッシュキー生成
	_, hashKey, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// スケジュール登録
	if err := u.r.InsertSchedule(tx, &model.UserSchedule{
		HashKey:      *hashKey,
		UserHashKeys: req.UserHashKeys,
		FreqID:       req.FreqID,
		Start:        req.Start,
		End:          req.End,
		Title:        req.Title,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// スケジュール一覧
func (u *UserService) Schedules() (*model.UserSchedulesResponse, *model.ErrorResponse) {
	res, err := u.r.ListSchedule()
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.UserSchedulesResponse{List: res}, nil
}

// スケジュール削除
func (u *UserService) DeleteSchedule(req *model.UserSchedule) *model.ErrorResponse {
	if err := u.v.ScheduleHashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// 削除
	if err := u.r.DeleteSchedule(tx, req); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &model.ErrorResponse{
				Status: http.StatusInternalServerError,
			}
		}
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 予約表提示
func (u *UserService) DispReserveTable() (*model.ReserveTableResponse, *model.ErrorResponse) {
	const WEEKS = 7
	const RESERVE_DURATION = 2 * WEEKS

	// TZをAsia/Tokyoに
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// TODO 応募者の面接可能グループ取得(一旦全グループ可能であるとする)
	var availabilityGroups []model.UserGroupResponse
	list, err := u.r.SearchGroup()
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}
	for _, row := range list {
		availabilityGroups = append(availabilityGroups, row)
	}

	// TODO 異なるグループメンバーでの面接可否設定取得(一旦不可能とする)
	isAvailabilityDifferentGroupMeeting := false

	// スケジュール一覧
	schedulesUTC, err := u.r.ListSchedule()
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}
	var schedules []model.UserScheduleResponse
	for _, row := range schedulesUTC {
		start := row.Start.In(jst)
		end := row.End.In(jst)
		row.Start = start
		row.End = end
		schedules = append(schedules, row)
	}

	// 日本の休日取得
	holidays, err := u.o.HolidaysJp(time.Now().Year())
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// 祝日かどうかの判定
	var reserveTable []model.ReserveTable
	var reserveTime []model.ReserveTime
	start := time.Now().AddDate(0, 0, WEEKS).In(jst)
	for i := 0; i < RESERVE_DURATION; i++ {
		isReserve := true

		s := time.Date(
			start.AddDate(0, 0, i).Year(),
			start.AddDate(0, 0, i).Month(),
			start.AddDate(0, 0, i).Day(),
			0,
			0,
			0,
			0,
			jst,
		)

		for _, holiday := range holidays {
			y1, m1, d1 := s.Date()
			y2, m2, d2 := holiday.Date()
			if y1 == y2 && m1 == m2 && d1 == d2 {
				isReserve = false
				break
			}
		}

		for d := s.Add(9 * time.Hour); d.Day() == s.Day() && d.Hour() <= 20; d = d.Add(30 * time.Minute) {
			if !isReserve {
				reserveTime = append(reserveTime, model.ReserveTime{
					Time:      d,
					IsReserve: false,
				})
			} else {
				var reserveOfGroups []model.ReserveOfGroup

				// グループ毎の面接可能人数
				for _, schedule := range schedules {
					userHashKeys := strings.Split(schedule.UserHashKeys, ",")

					for _, group := range availabilityGroups {
						var min uint = uint(len(strings.Split(group.Users, ",")))
						for _, userHashKey := range userHashKeys {
							var sum uint = uint(len(strings.Split(group.Users, ",")))

							// 時刻dは対象範囲内か
							if d.After(schedule.Start.Add(-1*time.Minute)) && d.Before(schedule.End.Add(1*time.Minute)) {
								// ユーザーグループまたはユーザーのハッシュキーか
								if userHashKey == group.HashKey {
									sum -= sum
								} else if strings.Contains(group.Users, userHashKey) {
									sum--
								}
							}
							if min > sum {
								min = sum
							}
						}
						reserveOfGroups = append(reserveOfGroups, model.ReserveOfGroup{
							HashKey: group.HashKey,
							Count:   min,
						})
					}
				}

				fmt.Println(reserveOfGroups)
				if isAvailabilityDifferentGroupMeeting {
					var sum uint = 0
					for _, reserve := range reserveOfGroups {
						sum += reserve.Count
					}
					reserveTime = append(reserveTime, model.ReserveTime{
						Time:      d,
						IsReserve: sum > 1,
					})
				} else {
					isMore2 := false
					for _, reserve := range reserveOfGroups {
						if reserve.Count > 1 {
							isMore2 = true
							reserveTime = append(reserveTime, model.ReserveTime{
								Time:      d,
								IsReserve: true,
							})
							break
						}
					}
					if !isMore2 {
						reserveTime = append(reserveTime, model.ReserveTime{
							Time:      d,
							IsReserve: false,
						})
					}
				}
			}
		}

		reserveTable = append(reserveTable, model.ReserveTable{
			Date:    s,
			Options: reserveTime,
		})
		reserveTime = []model.ReserveTime{}
	}

	return &model.ReserveTableResponse{
		List: reserveTable,
	}, nil
}

func GenerateHash(minLength, maxLength int) (*string, *string, error) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	length, err := rand.Int(rand.Reader, big.NewInt(int64(maxLength-minLength+1)))
	if err != nil {
		return nil, nil, err
	}
	strLength := minLength + int(length.Int64())

	buffer := make([]byte, strLength)
	_, err = rand.Read(buffer)
	if err != nil {
		return nil, nil, err
	}
	for i := 0; i < strLength; i++ {
		buffer[i] = chars[int(buffer[i])%len(chars)]
	}
	str := string(buffer)

	buffer2, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return nil, nil, err
	}
	hash := string(buffer2)

	return &str, &hash, nil
}

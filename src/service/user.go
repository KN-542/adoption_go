package service

import (
	"api/resources/static"
	"api/src/model"
	"api/src/model/enum"
	"api/src/repository"
	"api/src/validator"
	"crypto/rand"
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
	CreateSchedule(req *model.UserScheduleRequest) (*string, *model.ErrorResponse)
	// スケジュール更新
	UpdateSchedule(req *model.UserScheduleRequest) *model.ErrorResponse
	// スケジュール一覧
	Schedules() (*model.UserSchedulesResponse, *model.ErrorResponse)
	// スケジュール削除
	DeleteSchedule(req *model.UserSchedule) *model.ErrorResponse
	// 予約表提示
	DispReserveTable() (*model.ReserveTable, *model.ErrorResponse)
}

type UserService struct {
	r  repository.IUserRepository
	ra repository.IApplicantRepository
	m  repository.IMasterRepository
	v  validator.IUserValidator
	d  repository.IDBRepository
	o  repository.IOuterIFRepository
}

func NewUserService(
	r repository.IUserRepository,
	ra repository.IApplicantRepository,
	m repository.IMasterRepository,
	v validator.IUserValidator,
	d repository.IDBRepository,
	o repository.IOuterIFRepository,
) IUserService {
	return &UserService{r, ra, m, v, d, o}
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
func (u *UserService) CreateSchedule(req *model.UserScheduleRequest) (*string, *model.ErrorResponse) {
	// バリデーション
	if err := u.v.CreateScheduleValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &model.ErrorResponse{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
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

	return hashKey, nil
}

// スケジュール更新
func (u *UserService) UpdateSchedule(req *model.UserScheduleRequest) *model.ErrorResponse {
	tx, err := u.d.TxStart()
	if err != nil {
		return &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// スケジュール更新
	if err := u.r.UpdateSchedule(tx, &model.UserSchedule{
		HashKey:      req.HashKey,
		UserHashKeys: req.UserHashKeys,
		UpdatedAt:    time.Now(),
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

	// 応募者側更新
	if err := u.ra.UpdateDesiredAt(tx, &model.Applicant{
		HashKey: req.ApplicantHashKey,
		Users:   req.UserHashKeys,
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

// スケジュール一覧 (バッチでも実行したい)
func (u *UserService) Schedules() (*model.UserSchedulesResponse, *model.ErrorResponse) {
	schedulesBefore, err := u.r.ListSchedule()
	if err != nil {
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

	// 日付が過去の場合、更新or削除
	for _, schedule := range schedulesBefore {
		if schedule.Start.Before(time.Now()) {
			// なしの場合
			if schedule.FreqID == uint(enum.FREQ_NONE) {
				if err := u.r.DeleteSchedule(tx, &model.UserSchedule{
					HashKey: schedule.HashKey,
				}); err != nil {
					if err := u.d.TxRollback(tx); err != nil {
						return nil, &model.ErrorResponse{
							Status: http.StatusInternalServerError,
						}
					}
					return nil, &model.ErrorResponse{
						Status: http.StatusInternalServerError,
					}
				}
			} else {
				s := schedule.Start
				e := schedule.End
				if schedule.FreqID == uint(enum.FREQ_DAILY) {
					s = s.AddDate(0, 0, 1)
					e = e.AddDate(0, 0, 1)
				}
				if schedule.FreqID == uint(enum.FREQ_WEEKLY) {
					s = s.AddDate(0, 0, 7)
					e = e.AddDate(0, 0, 7)
				}
				if schedule.FreqID == uint(enum.FREQ_MONTHLY) {
					s = s.AddDate(0, 1, 0)
					e = e.AddDate(0, 1, 0)
				}
				if schedule.FreqID == uint(enum.FREQ_YEARLY) {
					s = s.AddDate(1, 0, 0)
					e = e.AddDate(1, 0, 0)
				}

				if err := u.r.UpdatePastSchedule(tx, &model.UserSchedule{
					HashKey:   schedule.HashKey,
					Start:     s,
					End:       e,
					UpdatedAt: time.Now(),
				}); err != nil {
					if err := u.d.TxRollback(tx); err != nil {
						return nil, &model.ErrorResponse{
							Status: http.StatusInternalServerError,
						}
					}
					return nil, &model.ErrorResponse{
						Status: http.StatusInternalServerError,
					}
				}
			}
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	schedulesAfter, err := u.r.ListSchedule()
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	return &model.UserSchedulesResponse{List: schedulesAfter}, nil
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
func (u *UserService) DispReserveTable() (*model.ReserveTable, *model.ErrorResponse) {
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

	// TZを日本に
	var schedulesJST []model.UserScheduleResponse
	for _, row := range schedulesUTC {
		start := row.Start.In(jst)
		end := row.End.In(jst)
		row.Start = start
		row.End = end
		schedulesJST = append(schedulesJST, row)
	}

	// スケジュールの頻度が「毎日」と「毎週」の場合、コピー
	var schedules []model.UserScheduleResponse
	start := time.Now().AddDate(0, 0, WEEKS).In(jst)
	s := time.Date(
		start.Year(),
		start.Month(),
		start.Day(),
		0,
		0,
		0,
		0,
		jst,
	)
	for _, row := range schedulesJST {
		if row.FreqID == uint(enum.FREQ_NONE) || row.FreqID == uint(enum.FREQ_MONTHLY) || row.FreqID == uint(enum.FREQ_YEARLY) {
			schedules = append(schedules, row)
			continue
		}

		for i := 0; i < RESERVE_DURATION; i++ {
			s_0 := time.Date(
				s.AddDate(0, 0, i).Year(),
				s.AddDate(0, 0, i).Month(),
				s.AddDate(0, 0, i).Day(),
				row.Start.Hour(),
				row.Start.Minute(),
				row.Start.Second(),
				row.Start.Nanosecond(),
				jst,
			)
			e_0 := time.Date(
				s.AddDate(0, 0, i).Year(),
				s.AddDate(0, 0, i).Month(),
				s.AddDate(0, 0, i).Day(),
				row.End.Hour(),
				row.End.Minute(),
				row.End.Second(),
				row.End.Nanosecond(),
				jst,
			)
			if row.FreqID == uint(enum.FREQ_DAILY) || (row.FreqID == uint(enum.FREQ_WEEKLY) && s_0.Weekday() == row.Start.Weekday()) {
				schedules = append(schedules, model.UserScheduleResponse{
					HashKey:      row.HashKey,
					UserHashKeys: row.UserHashKeys,
					Title:        row.Title,
					FreqID:       row.FreqID,
					Freq:         row.Freq,
					Start:        s_0,
					End:          e_0,
				})
			}
		}
	}

	// 日本の休日取得
	holidays, err := u.o.HolidaysJp(time.Now().Year())
	if err != nil {
		return nil, &model.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}

	// 祝日かどうかの判定
	var times []time.Time
	var reserveTime []model.ReserveTime
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

		times = append(times, s)

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
				for _, group := range availabilityGroups {
					var sum uint = uint(len(strings.Split(group.Users, ",")))

					for _, schedule := range schedules {
						userHashKeys := strings.Split(schedule.UserHashKeys, ",")

						for _, userHashKey := range userHashKeys {
							// 時刻dは対象範囲内か
							if d.After(schedule.Start.Add(-1*time.Minute)) && d.Before(schedule.End.Add(1*time.Minute)) {
								// ユーザーグループまたはユーザーのハッシュキーか
								if userHashKey == group.HashKey {
									sum = 0
								} else if strings.Contains(group.Users, userHashKey) {
									if sum > 0 {
										sum--
									}
								} else {
									continue
								}
							}
						}
					}

					reserveOfGroups = append(reserveOfGroups, model.ReserveOfGroup{
						HashKey: group.HashKey,
						Count:   sum,
					})
				}

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
	}

	return &model.ReserveTable{
		Dates:   times,
		Options: reserveTime,
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

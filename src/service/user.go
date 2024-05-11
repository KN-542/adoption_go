package service

import (
	"api/resources/static"
	"api/src/model/ddl"
	"api/src/model/enum"
	"api/src/model/response"
	"api/src/repository"
	"api/src/validator"
	"log"
	"net/http"
	"strings"
	"time"
)

type IUserService interface {
	// 一覧
	List() (*ddl.UsersResponse, *response.Error)
	// 登録
	Create(req *ddl.User) (*ddl.UserResponse, *response.Error)
	// 取得
	Get(req *ddl.User) (*ddl.UserResponse, *response.Error)
	// 検索(チーム)
	SearchTeams() (*ddl.TeamsResponse, *response.Error)
	// チーム登録
	CreateTeam(req *ddl.TeamRequest) *response.Error
	// スケジュール登録種別一覧
	ListScheduleType() (*ddl.CalendarsFreqStatus, *response.Error)
	// スケジュール登録
	CreateSchedule(req *ddl.UserScheduleRequest) (*string, *response.Error)
	// スケジュール更新
	UpdateSchedule(req *ddl.UserScheduleRequest) *response.Error
	// スケジュール一覧
	Schedules() (*ddl.UserSchedulesResponse, *response.Error)
	// スケジュール削除
	DeleteSchedule(req *ddl.UserSchedule) *response.Error
	// 予約表提示
	DispReserveTable() (*ddl.ReserveTable, *response.Error)
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
func (u *UserService) List() (*ddl.UsersResponse, *response.Error) {
	users, err := u.r.List()
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &ddl.UsersResponse{
		Users: users,
	}, nil
}

// 登録
func (u *UserService) Create(req *ddl.User) (*ddl.UserResponse, *response.Error) {
	// バリデーション
	if err := u.v.CreateValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// メールアドレス重複チェック
	if err := u.r.EmailDuplCheck(req); err != nil {
		return nil, &response.Error{
			Status: http.StatusConflict,
			Code:   static.CODE_USER_EMAIL_DUPL,
		}
	}

	// 初回パスワード発行
	password, hashPassword, err := GenerateHash(8, 16)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ハッシュキー生成
	_, hashKey, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 登録
	user := ddl.User{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: *hashKey,
		},
		Name:         req.Name,
		Email:        req.Email,
		Password:     *hashPassword,
		InitPassword: *hashPassword,
		RoleID:       req.RoleID,
	}

	_, err2 := u.r.Insert(tx, &user)
	if err2 != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	res := ddl.UserResponse{
		Email:        user.Email,
		InitPassword: *password,
	}
	return &res, nil
}

// 取得
func (u *UserService) Get(req *ddl.User) (*ddl.UserResponse, *response.Error) {
	if err := u.v.HashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	user, err := u.r.Get(req)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &ddl.UserResponse{
		HashKey: user.HashKey,
		Name:    user.Name,
		Email:   user.Email,
		RoleID:  user.RoleID,
	}, nil
}

// 検索(チーム)
func (u *UserService) SearchTeams() (*ddl.TeamsResponse, *response.Error) {
	teams, err := u.r.SearchTeam()
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ユーザー存在確認
	for index, team := range teams {
		var l []string
		users, err := u.r.ConfirmUserByHashKeys(strings.Split(team.Users, ","))
		if err != nil {
			if err != nil {
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
		}
		if users == nil || len(users) == 0 {
			teams[index].Users = ""
			continue
		}

		for _, user := range users {
			l = append(l, user.Name)
		}
		teams[index].Users = strings.Join(l, ",")
	}

	return &ddl.TeamsResponse{
		Teams: teams,
	}, nil
}

// チーム登録
func (u *UserService) CreateTeam(req *ddl.TeamRequest) *response.Error {
	// バリデーション
	if err := u.v.CreateTeamValidate(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ユーザー存在確認
	users, err := u.r.GetUserBasicByHashKeys(strings.Split(req.Users, ","))
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ハッシュキー生成
	_, hashKey, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// チーム登録
	req.HashKey = *hashKey
	team, err := u.r.InsertTeam(tx, &req.Team)
	if err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	for _, row := range users {
		// チーム紐づけ登録
		if err := u.r.InsertTeamAssociation(tx, &ddl.TeamAssociation{
			TeamID: uint(team.ID),
			UserID: row.ID,
		}); err != nil {
			if err := u.d.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// スケジュール登録種別一覧
func (u *UserService) ListScheduleType() (*ddl.CalendarsFreqStatus, *response.Error) {
	res, err := u.m.SelectCalendarFreqStatus()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &ddl.CalendarsFreqStatus{
		List: res,
	}, nil
}

// スケジュール登録
func (u *UserService) CreateSchedule(req *ddl.UserScheduleRequest) (*string, *response.Error) {
	// バリデーション
	if err := u.v.CreateScheduleValidate(req); err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	// ユーザー存在確認
	users, err := u.r.GetUserBasicByHashKeys(strings.Split(req.UserHashKeys, ","))
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// ハッシュキー生成
	_, hashKey, err := GenerateHash(1, 25)
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// スケジュール登録
	id, err := u.r.InsertSchedule(tx, &ddl.UserSchedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey: *hashKey,
		},
		InterviewFlg: req.InterviewFlg,
		FreqID:       req.FreqID,
		Start:        req.Start,
		End:          req.End,
		Title:        req.Title,
	})
	if err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	for _, row := range users {
		// スケジュール紐づけ登録
		if err := u.r.InsertScheduleAssociation(tx, &ddl.UserScheduleAssociation{
			UserScheduleID: *id,
			UserID:         row.ID,
		}); err != nil {
			if err := u.d.TxRollback(tx); err != nil {
				return nil, &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return nil, &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return hashKey, nil
}

// スケジュール更新
func (u *UserService) UpdateSchedule(req *ddl.UserScheduleRequest) *response.Error {
	// ユーザー存在確認
	users, err := u.r.GetUserBasicByHashKeys(strings.Split(req.UserHashKeys, ","))
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// スケジュール更新
	id, err := u.r.UpdateSchedule(tx, &ddl.UserSchedule{
		AbstractTransactionModel: ddl.AbstractTransactionModel{
			HashKey:   req.HashKey,
			UpdatedAt: time.Now(),
		},
	})
	if err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	// 問答無用で紐づけテーブルの該当スケジュールIDのレコード削除
	if err := u.r.DeleteScheduleAssociation(tx, &ddl.UserScheduleAssociation{
		UserScheduleID: *id,
	}); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	for _, row := range users {
		// スケジュール紐づけ登録
		if err := u.r.InsertScheduleAssociation(tx, &ddl.UserScheduleAssociation{
			UserScheduleID: *id,
			UserID:         row.ID,
		}); err != nil {
			if err := u.d.TxRollback(tx); err != nil {
				return &response.Error{
					Status: http.StatusInternalServerError,
				}
			}
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// スケジュール一覧 (バッチでも実行したい)
func (u *UserService) Schedules() (*ddl.UserSchedulesResponse, *response.Error) {
	schedulesBefore, err := u.r.ListSchedule()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 日付が過去の場合、更新or削除
	for _, schedule := range schedulesBefore {
		if schedule.Start.Before(time.Now()) {
			// なしの場合
			if schedule.FreqID == uint(enum.FREQ_NONE) {
				if err := u.r.DeleteSchedule(tx, &ddl.UserSchedule{
					AbstractTransactionModel: ddl.AbstractTransactionModel{
						HashKey: schedule.HashKey,
					},
				}); err != nil {
					if err := u.d.TxRollback(tx); err != nil {
						return nil, &response.Error{
							Status: http.StatusInternalServerError,
						}
					}
					return nil, &response.Error{
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

				if err := u.r.UpdatePastSchedule(tx, &ddl.UserSchedule{
					AbstractTransactionModel: ddl.AbstractTransactionModel{
						HashKey:   schedule.HashKey,
						UpdatedAt: time.Now(),
					},
					Start: s,
					End:   e,
				}); err != nil {
					if err := u.d.TxRollback(tx); err != nil {
						return nil, &response.Error{
							Status: http.StatusInternalServerError,
						}
					}
					return nil, &response.Error{
						Status: http.StatusInternalServerError,
					}
				}
			}
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	schedulesAfter, err := u.r.ListSchedule()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return &ddl.UserSchedulesResponse{List: schedulesAfter}, nil
}

// スケジュール削除
func (u *UserService) DeleteSchedule(req *ddl.UserSchedule) *response.Error {
	if err := u.v.ScheduleHashKeyValidate(req); err != nil {
		log.Printf("%v", err)
		return &response.Error{
			Status: http.StatusBadRequest,
			Code:   static.CODE_BAD_REQUEST,
		}
	}

	tx, err := u.d.TxStart()
	if err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 削除
	if err := u.r.DeleteSchedule(tx, req); err != nil {
		if err := u.d.TxRollback(tx); err != nil {
			return &response.Error{
				Status: http.StatusInternalServerError,
			}
		}
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	if err := u.d.TxCommit(tx); err != nil {
		return &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

// 予約表提示
func (u *UserService) DispReserveTable() (*ddl.ReserveTable, *response.Error) {
	const WEEKS = 7
	const RESERVE_DURATION = 2 * WEEKS

	// TZをAsia/Tokyoに
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Printf("%v", err)
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// TODO 応募者の面接可能チーム取得(一旦全チーム可能であるとする)
	var availabilityTeams []ddl.TeamResponse
	list, err := u.r.SearchTeam()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}
	for _, row := range list {
		availabilityTeams = append(availabilityTeams, row)
	}

	// TODO 異なるチームメンバーでの面接可否設定取得(一旦不可能とする)
	isAvailabilityDifferentTeamMeeting := false

	// スケジュール一覧
	schedulesUTC, err := u.r.ListSchedule()
	if err != nil {
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// TZを日本に
	var schedulesJST []ddl.UserScheduleResponse
	for _, row := range schedulesUTC {
		start := row.Start.In(jst)
		end := row.End.In(jst)
		row.Start = start
		row.End = end
		schedulesJST = append(schedulesJST, row)
	}

	// スケジュールの頻度が「毎日」と「毎週」の場合、コピー
	var schedules []ddl.UserScheduleResponse
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
				schedules = append(schedules, ddl.UserScheduleResponse{
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
		return nil, &response.Error{
			Status: http.StatusInternalServerError,
		}
	}

	// 祝日かどうかの判定
	var times []time.Time
	var reserveTime []ddl.ReserveTime
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
				reserveTime = append(reserveTime, ddl.ReserveTime{
					Time:      d,
					IsReserve: false,
				})
			} else {
				var reserveOfTeams []ddl.ReserveOfTeam

				// チーム毎の面接可能人数
				for _, team := range availabilityTeams {
					var sum uint = uint(len(strings.Split(team.Users, ",")))

					for _, schedule := range schedules {
						userHashKeys := strings.Split(schedule.UserHashKeys, ",")

						for _, userHashKey := range userHashKeys {
							// 時刻dは対象範囲内か
							if d.After(schedule.Start.Add(-1*time.Minute)) && d.Before(schedule.End.Add(1*time.Minute)) {
								// チームまたはユーザーのハッシュキーか
								if userHashKey == team.HashKey {
									sum = 0
								} else if strings.Contains(team.Users, userHashKey) {
									if sum > 0 {
										sum--
									}
								} else {
									continue
								}
							}
						}
					}

					reserveOfTeams = append(reserveOfTeams, ddl.ReserveOfTeam{
						HashKey: team.HashKey,
						Count:   sum,
					})
				}

				if isAvailabilityDifferentTeamMeeting {
					var sum uint = 0
					for _, reserve := range reserveOfTeams {
						sum += reserve.Count
					}
					reserveTime = append(reserveTime, ddl.ReserveTime{
						Time:      d,
						IsReserve: sum > 1,
					})
				} else {
					isMore2 := false
					for _, reserve := range reserveOfTeams {
						if reserve.Count > 1 {
							isMore2 = true
							reserveTime = append(reserveTime, ddl.ReserveTime{
								Time:      d,
								IsReserve: true,
							})
							break
						}
					}
					if !isMore2 {
						reserveTime = append(reserveTime, ddl.ReserveTime{
							Time:      d,
							IsReserve: false,
						})
					}
				}
			}
		}
	}

	return &ddl.ReserveTable{
		Dates:   times,
		Options: reserveTime,
	}, nil
}

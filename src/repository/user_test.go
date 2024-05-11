package repository

import (
	"api/src/infra"
	"api/src/model/ddl"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestUserRepository_Insert(t *testing.T) {
	type args struct {
		m *ddl.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// ok
		{
			"ok",
			args{&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey:   "abc",
					CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
				Name:         "taro",
				Email:        "taro@au.com",
				Password:     "root",
				InitPassword: "root",
				RoleID:       1,
				RefreshToken: "token",
			}},
			false,
		},
		// ng_hash_key
		{
			"ng_hash_key",
			args{&ddl.User{
				Name:         "taro",
				Email:        "taro@au.com",
				Password:     "root",
				InitPassword: "root",
				RoleID:       1,
			}},
			true,
		},
		// ng_name
		{
			"ng_name",
			args{&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: "abc",
				},
				Email:        "taro@au.com",
				Password:     "root",
				InitPassword: "root",
				RoleID:       1,
			}},
			true,
		},
		// ng_email
		{
			"ng_email",
			args{&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: "abc",
				},
				Name:         "taro",
				Password:     "root",
				InitPassword: "root",
				RoleID:       1,
			}},
			true,
		},
		// ng_password
		{
			"ng_password",
			args{&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: "abc",
				},
				Name:     "taro",
				Email:    "taro@au.com",
				Password: "root",
				RoleID:   1,
			}},
			true,
		},
		// ng_init_password
		{
			"ng_init_password",
			args{&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: "abc",
				},
				Name:         "taro",
				Email:        "taro@au.com",
				InitPassword: "root",
				RoleID:       1,
			}},
			true,
		},
		// ng_role_id
		{
			"ng_role_id",
			args{&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: "abc",
				},
				Name:         "taro",
				Email:        "taro@au.com",
				Password:     "root",
				InitPassword: "root",
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: infra.NewDB(),
			}
			tx := u.db.Begin()
			if err := tx.Error; err != nil {
				t.Errorf("UserRepository.Insert() error = %v", err)
			}

			_, err := u.Insert(tx, tt.args.m)

			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					t.Errorf("UserRepository.Insert() error = %v", err)
				}
				if !tt.wantErr {
					t.Errorf("UserRepository.Insert() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			if err == nil {
				if err := tx.Commit().Error; err != nil {
					t.Errorf("UserRepository.Insert() error = %v", err)
				}

				_, err := u.Get(tt.args.m)
				if err != nil {
					if err := tx.Rollback().Error; err != nil {
						t.Errorf("UserRepository.Get() error = %v", err)
					}
					t.Errorf("UserRepository.Get() error = %v", err)
				}

				tx := u.db.Begin()
				if err := u.Delete(tx, tt.args.m); err != nil {
					if err := tx.Rollback().Error; err != nil {
						t.Errorf("UserRepository.Delete() error = %v", err)
					}
					t.Errorf("UserRepository.Delete() error = %v", err)
				}
				if err := tx.Commit().Error; err != nil {
					t.Errorf("UserRepository.Delete() error = %v", err)
				}

				if tt.wantErr {
					t.Errorf("UserRepository.Insert() error = %v, wantErr %v", nil, tt.wantErr)
				}
			}
		})
	}
}

// TODO
func TestUserRepository_List(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ddl.UserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			got, err := u.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_Get(t *testing.T) {
	type args struct {
		m *ddl.User
	}
	tests := []struct {
		name    string
		args    args
		want    *ddl.User
		wantErr bool
	}{
		// ok
		{
			"ok",
			args{&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: "abc",
				},
			}},
			&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey:   "abc",
					CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
				Name:         "taro",
				Email:        "taro@au.com",
				Password:     "root",
				InitPassword: "root",
				RoleID:       1,
				RefreshToken: "token",
			},
			false,
		},
		// ng 0件
		{
			"ng_0",
			args{&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey: "ng_ab",
				},
			}},
			&ddl.User{
				AbstractTransactionModel: ddl.AbstractTransactionModel{
					HashKey:   "ng_abc",
					CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
				Name:         "taro",
				Email:        "taro@au.com",
				Password:     "root",
				InitPassword: "root",
				RoleID:       1,
				RefreshToken: "token",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: infra.NewDB(),
			}
			tx := u.db.Begin()
			if err := tx.Error; err != nil {
				t.Errorf("UserRepository.Insert() error = %v", err)
			}

			_, err := u.Insert(tx, tt.want)
			if err != nil {
				if err := tx.Rollback().Error; err != nil {
					t.Errorf("UserRepository.Insert() error = %v", err)
				}
				t.Errorf("UserRepository.Insert() error = %v", err)
			}
			if err := tx.Commit().Error; err != nil {
				t.Errorf("UserRepository.Insert() error = %v", err)
			}

			got, err := u.Get(tt.args.m)
			if err != nil {
				tx2 := u.db.Begin()
				if err := u.Delete(tx2, tt.want); err != nil {
					if err := tx2.Rollback().Error; err != nil {
						t.Errorf("UserRepository.Delete() error = %v", err)
					}
					t.Errorf("UserRepository.Delete() error = %v", err)
				}
				if err := tx2.Commit().Error; err != nil {
					t.Errorf("UserRepository.Delete() error = %v", err)
				}

				if !tt.wantErr {
					t.Errorf("UserRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !(got == nil || (got.Name == tt.want.Name &&
				got.Email == tt.want.Email &&
				got.Password == tt.want.Password &&
				got.InitPassword == tt.want.InitPassword &&
				got.RoleID == tt.want.RoleID &&
				got.RefreshToken == tt.want.RefreshToken &&
				got.CreatedAt.Sub(tt.want.CreatedAt) < time.Second &&
				got.UpdatedAt.Sub(tt.want.UpdatedAt) < time.Second)) {
				tx2 := u.db.Begin()
				if err := u.Delete(tx2, tt.want); err != nil {
					if err := tx2.Rollback().Error; err != nil {
						t.Errorf("UserRepository.Delete() error = %v", err)
					}
					t.Errorf("UserRepository.Delete() error = %v", err)
				}
				if err := tx2.Commit().Error; err != nil {
					t.Errorf("UserRepository.Delete() error = %v", err)
				}

				t.Errorf("UserRepository.Get() = %v, want %v", got, tt.want)
			}

			tx2 := u.db.Begin()
			if err := u.Delete(tx2, tt.want); err != nil {
				if err := tx2.Rollback().Error; err != nil {
					t.Errorf("UserRepository.Delete() error = %v", err)
				}
				t.Errorf("UserRepository.Delete() error = %v", err)
			}
			if err := tx2.Commit().Error; err != nil {
				t.Errorf("UserRepository.Delete() error = %v", err)
			}

			if tt.wantErr {
				t.Errorf("UserRepository.Get() error = %v, wantErr %v", nil, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		m  *ddl.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			if err := u.Update(tt.args.tx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		m  *ddl.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			if err := u.Delete(tt.args.tx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_ConfirmUserByHashKeys(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		hashKeys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []ddl.UserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			got, err := u.ConfirmUserByHashKeys(tt.args.hashKeys)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.ConfirmUserByHashKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.ConfirmUserByHashKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_Login(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		m *ddl.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []ddl.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			got, err := u.Login(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_PasswordChange(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		m  *ddl.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			if err := u.PasswordChange(tt.args.tx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.PasswordChange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_ConfirmInitPassword(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		m *ddl.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *int8
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			got, err := u.ConfirmInitPassword(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.ConfirmInitPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserRepository.ConfirmInitPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_ConfirmInitPassword2(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		m *ddl.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			got, err := u.ConfirmInitPassword2(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.ConfirmInitPassword2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserRepository.ConfirmInitPassword2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_EmailDuplCheck(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		m *ddl.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			if err := u.EmailDuplCheck(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.EmailDuplCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_SearchTeam(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ddl.TeamResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			got, err := u.SearchTeam()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.SearchTeam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.SearchTeam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_InsertTeam(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		m  *ddl.Team
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			_, err := u.InsertTeam(tt.args.tx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.InsertTeam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_InsertSchedule(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		m  *ddl.UserSchedule
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			if _, err := u.InsertSchedule(tt.args.tx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.InsertSchedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_ListSchedule(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ddl.UserScheduleResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			got, err := u.ListSchedule()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.ListSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.ListSchedule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_DeleteSchedule(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		m  *ddl.UserSchedule
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			if err := u.DeleteSchedule(tt.args.tx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.DeleteSchedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_UpdatePastSchedule(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		m  *ddl.UserSchedule
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			if err := u.UpdatePastSchedule(tt.args.tx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.UpdatePastSchedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

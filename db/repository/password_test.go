package repository

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func Test_passwordRepository_AddPassword(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		params AddPasswordParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Password
		wantErr bool
	}{
		{
			name: "Add password",
			fields: fields{
				db: gdb.Debug().Begin(),
			},
			args: args{
				params: AddPasswordParams{
					Item:     "item",
					Login:    "login",
					Password: "password",
					UserId:   1,
				},
			},
			want: Password{
				Item:     "item",
				Login:    "login",
				Password: "password",
				UserId:   1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &passwordRepository{
				db: tt.fields.db,
			}
			got, err := repo.AddPassword(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got.Item, tt.want.Item)
			assert.Equal(t, got.Login, tt.want.Login)
			assert.Equal(t, got.Password, tt.want.Password)
			assert.Equal(t, got.UserId, tt.want.UserId)
		})
	}
}

func Test_passwordRepository_GetAllPasswords(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		params GetAllPasswordsParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Get all passwords",
			fields: fields{
				db: gdb.Debug().Begin(),
			},
			args: args{
				params: GetAllPasswordsParams{
					UserId: 1,
				},
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &passwordRepository{
				db: tt.fields.db,
			}
			got, err := repo.GetAllPasswords(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPasswords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("GetAllPasswords() got = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_passwordRepository_GetPassword(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		params GetPasswordParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Password
		wantErr bool
	}{
		{
			name: "Get password",
			fields: fields{
				db: gdb.Debug().Begin(),
			},
			args: args{
				params: GetPasswordParams{
					Item: "item",
				},
			},
			want: Password{
				Item:  "item",
				Login: "login",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &passwordRepository{
				db: tt.fields.db,
			}
			got, err := repo.GetPassword(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Login, tt.want.Login) {
				t.Errorf("GetPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

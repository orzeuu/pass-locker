package repository

import (
	"github.com/orzeuu/pass-locker/password"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

var gdb *gorm.DB

func TestMain(m *testing.M) {
	var err error

	dbName := "db_test.db"
	exec.Command("rm", "-f", dbName)

	gdb, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	gdb.AutoMigrate(&User{}, &Password{})

	hashedPassword, _ := password.HashPassword("test")
	gdb.Create([]User{{Login: "test", Password: hashedPassword}})

	gdb.Create([]Password{{Item: "item", Login: "login", Password: "password", UserId: 1}})
	gdb.Create([]Password{{Item: "item2", Login: "login2", Password: "password2", UserId: 1}})

	defer exec.Command("rm", "-f", dbName)
	os.Exit(m.Run())
}

func Test_userRepository_GetUser(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		params GetUserParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    User
		wantErr bool
	}{
		{
			name: "Not found user",
			fields: fields{
				db: gdb.Debug().Begin(),
			},
			args: args{
				params: GetUserParams{
					Login:    "test2",
					Password: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "Found with incorrect password",
			fields: fields{
				db: gdb.Debug().Begin(),
			},
			args: args{
				params: GetUserParams{
					Login:    "test",
					Password: "test2",
				},
			},
			wantErr: true,
		},
		{
			name: "Get user",
			fields: fields{
				db: gdb.Debug().Begin(),
			},
			args: args{
				params: GetUserParams{
					Login:    "test",
					Password: "test",
				},
			},
			want: User{
				Login: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepository{
				db: tt.fields.db,
			}
			got, err := repo.GetUser(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Login, tt.want.Login) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_InsertUser(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		params InsertUserParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    User
		wantErr bool
	}{
		{
			name: "Insert user",
			fields: fields{
				gdb.Debug().Begin(),
			},
			args: args{
				params: InsertUserParams{
					Login:    "test2",
					Password: "test2",
				},
			},
			want: User{
				Login: "test2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepository{
				db: tt.fields.db,
			}
			got, err := repo.InsertUser(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Login, tt.want.Login) {
				t.Errorf("InsertUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

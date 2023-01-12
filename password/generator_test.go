package password

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want *Generator
	}{
		{
			name: "default config",
			args: args{nil},
			want: func() *Generator {
				cfg := &DefaultConfig
				cfg.CharacterSet = getCharacterSet(cfg)
				cfg.Length = LengthStrong
				return &Generator{cfg}
			}(),
		},
		{
			name: "set config",
			args: args{&Config{
				IncludeLowercaseLetters: true,
			}},
			want: func() *Generator {
				cfg := &Config{IncludeLowercaseLetters: true}
				cfg.CharacterSet = "abcdefghijklmnopqrstuvwxyz"
				cfg.Length = LengthStrong
				return &Generator{cfg}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGenerator(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generator_Generate(t *testing.T) {
	type fields struct {
		Config *Config
	}
	tests := []struct {
		name    string
		fields  fields
		test    func(string, string)
		wantErr bool
	}{
		{
			name:   "valid",
			fields: fields{&DefaultConfig},
			test: func(pwd, characterSet string) {
				assert.Len(t, pwd, int(DefaultConfig.Length))
				err := stringMatchesCharacters(pwd, characterSet)
				if err != nil {
					t.Errorf("Generate() error = %v", err)
					return
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := NewGenerator(tt.fields.Config)
			got, err := pg.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.test(got, pg.CharacterSet)
		})
	}
}

func Test_getCharacterSet(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "exclude similar characters",
			args: args{
				config: &Config{
					IncludeSymbols:           true,
					IncludeAllSymbols:        true,
					IncludeNumbers:           true,
					IncludeLowercaseLetters:  true,
					ExcludeSimilarCharacters: true,
				},
			},
			want: "abcdefghkmnpqrstuvwxyz23456789~!@#$%^&*-=+`()_{}[]\\|:;\"'<>,.?/",
		},
		{
			name: "exclude numbers",
			args: args{
				config: &Config{
					IncludeLowercaseLetters:  true,
					IncludeSymbols:           true,
					IncludeNumbers:           false,
					ExcludeSimilarCharacters: true,
					IncludeAllSymbols:        true,
				},
			},
			want: "abcdefghkmnpqrstuvwxyz~!@#$%^&*-=+`()_{}[]\\|:;\"'<>,.?/",
		},
		{
			name: "full list",
			args: args{
				config: &Config{
					IncludeNumbers:           true,
					IncludeLowercaseLetters:  true,
					IncludeSymbols:           true,
					IncludeAllSymbols:        true,
					ExcludeSimilarCharacters: false,
				},
			},
			want: "abcdefghijklmnopqrstuvwxyz0123456789~!@#$%^&*-=+`()_{}[]\\|:;\"'<>,.?/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCharacterSet(tt.args.config); got != tt.want {
				t.Errorf("getCharacterSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringMatchesCharacters(str, characters string) error {
	set := strings.Split(characters, "")
	strSet := strings.Split(str, "")

	for _, strChr := range strSet {
		found := false
		for _, setChr := range set {
			if strChr == setChr {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("%v should not be in the str", strChr)
		}
	}

	return nil
}

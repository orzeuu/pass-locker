package password

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	LengthStrong     uint = 24
	DefaultLetterSet      = "abcdefghijklmnopqrstuvwxyz"
	DefaultNumberSet      = "0123456789"
	DefaultSymbolSet      = "~!@#$%^&*-=+"
)

var (
	DefaultConfig = Config{
		Length:                  LengthStrong,
		IncludeNumbers:          true,
		IncludeLowercaseLetters: true,
		IncludeUppercaseLetters: true,
		IncludeSymbols:          true,
	}
)

type Generator struct {
	*Config
}

type Config struct {
	Length                  uint
	CharacterSet            string
	IncludeSymbols          bool
	IncludeNumbers          bool
	IncludeLowercaseLetters bool
	IncludeUppercaseLetters bool
}

func NewGenerator(config *Config) *Generator {
	if config == nil {
		config = &DefaultConfig
	}

	if config.Length == 0 {
		config.Length = LengthStrong
	}

	if config.CharacterSet == "" {
		config.CharacterSet = getCharacterSet(config)
	}

	return &Generator{
		Config: config,
	}
}

func (pg *Generator) Generate() (string, error) {
	var password string
	characterSet := strings.Split(pg.CharacterSet, "")
	max := big.NewInt(int64(len(characterSet)))

	for i := uint(0); i < pg.Config.Length; i++ {
		val, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		password += characterSet[val.Int64()]
	}

	return password, nil
}

func getCharacterSet(config *Config) string {
	var characterSet string
	if config.IncludeLowercaseLetters {
		characterSet += DefaultLetterSet
	}

	if config.IncludeUppercaseLetters {
		characterSet += strings.ToUpper(DefaultLetterSet)
	}

	if config.IncludeNumbers {
		characterSet += DefaultNumberSet
	}

	if config.IncludeSymbols {
		characterSet += DefaultSymbolSet
	}

	return characterSet
}

package custerr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func CustomError(err error) error {
	customErr := errors.New("")
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				customErr = fmt.Errorf("%w"+fmt.Sprintf(", %s wajib diisi", err.Field()), customErr)
			case "email":
				customErr = fmt.Errorf("%w"+fmt.Sprintf(", format %s tidak valid", err.Field()), customErr)
			}
		}
	} else if code, ok := err.(*mysql.MySQLError); ok {
		switch code.Number {
		case 1062:
			if strings.Contains(code.Message, "email") {
				customErr = fmt.Errorf("%w Email sudah terpakai", customErr)
			} else if strings.Contains(code.Message, "username") {
				customErr = fmt.Errorf("%w Username sudah terpakai", customErr)
			}
		case 1264:
			customErr = fmt.Errorf("%w Jumlah karakter terlalu panjang", customErr)
		case 1406:
			customErr = fmt.Errorf("%w Jumlah karakter terlalu panjang", customErr)
		default:
			customErr = err
		}
	} else {
		customErr = err
	}

	return customErr
}
package custerr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

func CustomError(err error) error {
	var customErr []string
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, err := range castedObject {
			switch err.Tag() {
			case "required":
				customErr = append(customErr, fmt.Sprintf("%s wajib diisi, ", err.Field()))
			case "email":
				customErr = append(customErr, fmt.Sprintf("format %s tidak valid,", err.Field()))
			}
		}
	} else if code, ok := err.(*mysql.MySQLError); ok {
		switch code.Number {
		case 1062:
			if strings.Contains(code.Message, "email") {
				customErr = append(customErr, fmt.Sprintf("Email sudah terpakai"))
			} else if strings.Contains(code.Message, "username") {
				customErr = append(customErr, fmt.Sprintf("Username sudah terpakai"))
			} else if strings.Contains(code.Message, "slug") {
				customErr = append(customErr, fmt.Sprintf("Domain sudah terpakai"))
			} else {
				customErr = append(customErr, fmt.Sprintf("Nama sudah terpakai"))
			}
		case 1264:
			customErr = append(customErr, fmt.Sprintf("Jumlah karakter terlalu panjang"))
		case 1406:
			customErr = append(customErr, fmt.Sprintf("Jumlah karakter terlalu panjang"))
		default:
			customErr = append(customErr, fmt.Sprintf(err.Error()))
		}
	} else {
		customErr = append(customErr, fmt.Sprintf(err.Error()))
	}

	return errors.New(strings.Join(customErr, "\n"))
}

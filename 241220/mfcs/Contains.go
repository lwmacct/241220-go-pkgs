package mfcs

import (
	"github.com/thoas/go-funk"
)

func Contains(in interface{}, elem interface{}) bool {
	return funk.Contains(in, elem)
}

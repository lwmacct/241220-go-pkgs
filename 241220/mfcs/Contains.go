package mfcs

import (
	"github.com/thoas/go-funk"
)

type Ts struct{}

func New() *Ts {
	return &Ts{}
}

func (t *Ts) Contains(in interface{}, elem interface{}) bool {
	return funk.Contains(in, elem)
}

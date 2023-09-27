package model

import (
	"time"

	"github.com/google/uuid"
)

type Date string

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)

	t, err := time.Parse(`"2006-01-02"`, s)
	if err != nil {
		return err
	}

	*d = Date(t.Format("2006-01-02"))
	return nil
}

type Pessoa struct {
	ID         uuid.UUID `json:"id" bson:"_id"`
	Nome       string    `json:"nome" validate:"required,min=1,max=32"`
	CpfCnpj    string    `json:"cpfCnpj" bson:"cpf_cnpj" validate:"required,max=14"`
	Nascimento Date      `json:"nascimento" bson:"nascimento" validate:"required"`
	Seguros    []string  `json:"seguros" bson:"seguros" validate:"omitempty,min=1,dive,max=32"`
}

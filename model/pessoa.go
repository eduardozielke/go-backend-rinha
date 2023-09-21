package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pessoa struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Nome       string             `json:"nome" validate:"required,min=1,max=32"`
	CpfCnpj    string             `json:"cpfCnpj" bson:"cpf_cnpj" validate:"required,max=14"`
	Nascimento string             `json:"nascimento" bson:"nascimento" validate:"required"`
	Seguros    []string           `json:"seguros" bson:"seguros" validate:"omitempty,min=1,dive,max=32"`
}

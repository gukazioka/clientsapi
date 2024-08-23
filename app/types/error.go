package types

import "errors"

var ErrorAlreadyExists = errors.New("client already exists")
var ErrorInvalidDocument = errors.New("invalid document")
var ErrorInvalidCpf = errors.New("invalid cpf")
var ErrorInvalidCnpj = errors.New("invalid cnpj")

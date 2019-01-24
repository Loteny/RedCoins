package cadastro

import (
	"regexp"

	"github.com/loteny/redcoins/erros"
)

// email verifica se o e-mail é válido (formato regex /.+@.+/)
func email(email string) error {
	matched, err := regexp.MatchString(".+@.+", email)
	if err != nil {
		return erros.CriaInternoPadrao(err)
	}
	if !matched {
		return ErrEmailInvalido
	}
	return nil
}

// senha verifica se a senha possui pelo menos 6 caracteres e no máximo 64
func senha(senha string) error {
	matched, err := regexp.MatchString("^.{6,64}$", senha)
	if err != nil {
		return erros.CriaInternoPadrao(err)
	}
	if !matched {
		return ErrSenhaInvalida
	}
	return nil
}

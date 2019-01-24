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

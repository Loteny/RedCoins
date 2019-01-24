package cadastro

import (
	"regexp"

	"github.com/loteny/redcoins/erros"
)

// email verifica se o e-mail é válido (formato regex /.+@.+/)
func email(email string) error {
	return validacaoMatchSimples(email, "^.+@.+$", ErrEmailInvalido)
}

// senha verifica se a senha possui pelo menos 6 caracteres e no máximo 64
func senha(senha string) error {
	return validacaoMatchSimples(senha, "^.{6,64}$", ErrSenhaInvalida)
}

// nome verifica se o campo não está vazio
func nome(nome string) error {
	return validacaoMatchSimples(nome, "^.{6,64}$", ErrNomeInvalido)
}

// validacaoMatchSimples executa uma validação básica com um regex passado como
// argumento. Retorna o erro gerado pela função do regex, caso houve algum, ou o
// erro passado por argumento para essa função no caso de o regex não bater com
// o dado passado.
func validacaoMatchSimples(s string, reg string, e error) error {
	matched, err := regexp.MatchString(reg, s)
	if err != nil {
		return erros.CriaInternoPadrao(err)
	}
	if !matched {
		return e
	}
	return nil
}

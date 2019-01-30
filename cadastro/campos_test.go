package cadastro

import (
	"strings"
	"testing"

	"github.com/loteny/redcoins/erros"
)

func TestEmailValido(t *testing.T) {
	if err := email("teste@gmail.com"); !erros.Vazio(err) {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestEmailInvalido(t *testing.T) {
	if err := email("emailinvalido.com"); err.Error() != ErrEmailInvalido.Error() {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestSenha(t *testing.T) {
	// Número abaixo do mínimo de caracteres
	if err := senha("abc"); err.Error() != ErrSenhaInvalida.Error() {
		t.Errorf("Erro retornado: %v", err)
	}
	// Número mínimo de caracteres
	if err := senha("abcdef"); !erros.Vazio(err) {
		t.Errorf("Erro retornado: %v", err)
	}
	// Número normal de caracteres
	if err := senha("abcdefgh"); !erros.Vazio(err) {
		t.Errorf("Erro retornado: %v", err)
	}
	// Número acima do máximo de caracteres
	if err := senha(strings.Repeat("a", 65)); err.Error() != ErrSenhaMuitoLonga.Error() {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestNome(t *testing.T) {
	// Nome vazio
	if err := nome(""); err.Error() != ErrNomeInvalido.Error() {
		t.Errorf("Erro retornado: %v", err)
	}
	// Nome válido
	if err := nome("Ronnie James Dio"); !erros.Vazio(err) {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestNascimento(t *testing.T) {
	// Formatação incorreta
	if err := nascimento(""); err.Error() != ErrNascimentoInvalido.Error() {
		t.Errorf("Erro retornado: %v", err)
	}
	if err := nascimento("9231-332"); err.Error() != ErrNascimentoInvalido.Error() {
		t.Errorf("Erro retornado: %v", err)
	}
	// Datas absurdas (30 de fevereiro)
	if err := nascimento("2018-02-30"); err.Error() != ErrNascimentoInvalido.Error() {
		t.Errorf("Erro retornado: %v", err)
	}
	// Datas futuras
	if err := nascimento("5012-01-25"); err.Error() != ErrNascimentoInvalido.Error() {
		t.Errorf("Erro retornado: %v", err)
	}
	// Datas passadas
	if err := nascimento("2015-03-07"); !erros.Vazio(err) {
		t.Errorf("Erro retornado: %v", err)
	}
}

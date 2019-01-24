package cadastro

import (
	"strings"
	"testing"
)

func TestEmailValido(t *testing.T) {
	if err := email("teste@gmail.com"); err != nil {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestEmailInvalido(t *testing.T) {
	if err := email("emailinvalido.com"); err != ErrEmailInvalido {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestSenha(t *testing.T) {
	// Número abaixo do mínimo de caracteres
	if err := senha("abc"); err != ErrSenhaInvalida {
		t.Errorf("Erro retornado: %v", err)
	}
	// Número mínimo de caracteres
	if err := senha("abcdef"); err != nil {
		t.Errorf("Erro retornado: %v", err)
	}
	// Número normal de caracteres
	if err := senha("abcdefgh"); err != nil {
		t.Errorf("Erro retornado: %v", err)
	}
	// Número acima do máximo de caracteres
	if err := senha(strings.Repeat("a", 65)); err != ErrSenhaInvalida {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestNome(t *testing.T) {
	// Nome vazio
	if err := nome(""); err != ErrNomeInvalido {
		t.Errorf("Erro retornado: %v", err)
	}
	// Nome válido
	if err := nome("Ronnie James Dio"); err != nil {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestNascimento(t *testing.T) {
	// Formatação incorreta
	if err := nascimento(""); err != ErrNascimentoInvalido {
		t.Errorf("Erro retornado: %v", err)
	}
	if err := nascimento("9231-332"); err != ErrNascimentoInvalido {
		t.Errorf("Erro retornado: %v", err)
	}
	// Datas absurdas (30 de fevereiro)
	if err := nascimento("2018-02-30"); err != ErrNascimentoInvalido {
		t.Errorf("Erro retornado: %v", err)
	}
	// Datas futuras
	if err := nascimento("5012-01-25"); err != ErrNascimentoInvalido {
		t.Errorf("Erro retornado: %v", err)
	}
	// Datas passadas
	if err := nascimento("2015-03-07"); err != nil {
		t.Errorf("Erro retornado: %v", err)
	}
}

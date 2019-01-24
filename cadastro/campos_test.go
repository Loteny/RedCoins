package cadastro

import (
	"strings"
	"testing"
)

func TestEmailValido(t *testing.T) {
	err := email("teste@gmail.com")
	if err != nil {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestEmailInvalido(t *testing.T) {
	err := email("emailinvalido.com")
	if err != ErrEmailInvalido {
		t.Errorf("Erro retornado: %v", err)
	}
}

func TestSenha(t *testing.T) {
	// Número abaixo do mínimo de caracteres
	err := senha("abc")
	if err != ErrSenhaInvalida {
		t.Errorf("Erro retornado: %v", err)
	}
	// Número mínimo de caracteres
	err = senha("abcdef")
	if err != nil {
		t.Errorf("Erro retornado: %v", err)
	}
	// Número normal de caracteres
	err = senha("abcdefgh")
	if err != nil {
		t.Errorf("Erro retornado: %v", err)
	}
	// Número acima do máximo de caracteres
	err = senha(strings.Repeat("a", 65))
	if err != ErrSenhaInvalida {
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

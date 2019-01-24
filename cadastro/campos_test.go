package cadastro

import "testing"

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

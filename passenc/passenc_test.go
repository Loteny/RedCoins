package passenc

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGeraHashed(t *testing.T) {
	hashed, err := GeraHashed([]byte("123456"))
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	// Checa manualmente a validade do resultado da função
	err = bcrypt.CompareHashAndPassword(hashed, []byte("123456"))
	if err != nil {
		t.Errorf("Resultado do hashing inválido: %v", err)
	}
	// Checa a verificação para uma senha incorreta
	err = bcrypt.CompareHashAndPassword(hashed, []byte("1234567"))
	if err != bcrypt.ErrMismatchedHashAndPassword {
		t.Errorf("Resultado do hashing inválido: %v", err)
	}
}

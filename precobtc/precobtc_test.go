package precobtc

import "testing"

func TestPrecoUnidade(t *testing.T) {
	// Verifica se o preço retornado é lógico (não-negativo) e não absurdo
	if preco, err := PrecoUnidade(); err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	} else if preco < 0 || preco > 1000000 {
		t.Fatalf("Preço inesperado: %v", preco)
	}
}

func TestPreco(t *testing.T) {
	// Verifica se o preço retornado é lógico (não-negativo) e não absurdo
	if preco, err := Preco(2); err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	} else if preco < 0 || preco > 1000000 {
		t.Fatalf("Preço inesperado: %v", preco)
	}
}

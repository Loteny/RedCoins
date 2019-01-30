// Para o bom funcionamento desse módulo de testes, por enquanto, o banco de
// dados de teste deve existir e estar devidamente populado.
package transacao

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCompraHTTP(t *testing.T) {
	// Formulário válido
	form := url.Values{}
	form.Set("email", "teste@gmail.com")
	form.Set("senha", "123456")
	form.Set("qtd", "0.03")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP := func(w http.ResponseWriter, r *http.Request) {
		err := CompraHTTP(r)
		if err != nil {
			t.Errorf("Erro inesperado na transação: %v", err)
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)

	// Quantidade inválida
	form.Set("qtd", "-2")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		err := CompraHTTP(r)
		if err != ErrQtdInvalida {
			t.Errorf("Erro inesperado na transação: %v", err)
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)
}

// testRealizaRequestHTTPPostForm é uma função auxiliar para geração de requests
// HTTP com formulário POST
func testRealizaRequestHTTPPostForm(t *testing.T, form url.Values,
	f func(w http.ResponseWriter, r *http.Request)) {
	// Criação do request HTTP
	request, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Recorder para armazenar a resposta
	recorder := httptest.NewRecorder()

	// Upload do servidor de teste
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(recorder, request)
}

package cadastro

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCadastraHTTPS(t *testing.T) {
	// Formulário para enviar no POST
	form := url.Values{}
	form.Set("email", "teste@gmail.com")
	form.Set("senha", "123456")
	form.Set("nome", "Ronnie James Dio")
	form.Set("nascimento", "1942-07-10")

	// Criação do request HTTP
	request, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Recorder para armazenar a resposta
	recorder := httptest.NewRecorder()

	// Upload do servidor de teste
	handler := http.HandlerFunc(CadastraHTTPS)
	handler.ServeHTTP(recorder, request)

	// Checa os status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status code incorreto (adquirido %v, esperado %v).",
			status, http.StatusOK)
	}

	// Checa o corpo da mensagem
	esperado := ``
	if recorder.Body.String() != esperado {
		t.Errorf("Corpo da mensagem incorreto.\nAdquirido: %v\nDesejado: %v",
			recorder.Body.String(), esperado)
	}
}

func TestValidaDadosCadastroRequestHTTP(t *testing.T) {
	//
}

func TestValidaDadosCadastro(t *testing.T) {
	// Dados válidos
	dados := dadosCadastrais{
		email:      "teste@gmail.com",
		senha:      "123456",
		nome:       "Ronnie James Dio",
		nascimento: "1942-07-10",
	}
	err := validaDadosCadastro(&dados)
	if err != nil {
		t.Errorf("Erro: %v", err)
	}
}

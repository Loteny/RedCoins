package cadastro

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCadastraHTTPS(t *testing.T) {
	// Formulário válido
	form := url.Values{}
	form.Set("email", "teste@gmail.com")
	form.Set("senha", "123456")
	form.Set("nome", "Ronnie James Dio")
	form.Set("nascimento", "1942-07-10")
	testRequestHTTPPostForm(t, form, CadastraHTTPS, http.StatusOK, ``)

	// Nascimento inválido
	form.Set("nascimento", "194207-10")
	testRequestHTTPPostForm(t, form, CadastraHTTPS, http.StatusBadRequest,
		`{"erro":"nascimento invalido"}`)
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

// testRequestHTTP é uma função auxiliar de para geração e tratamento de
// requests HTTP com formulário POST
func testRequestHTTPPostForm(t *testing.T, form url.Values,
	f func(http.ResponseWriter, *http.Request),
	statusCodeEsperado int, respostaEsperada string) {
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

	// Checa os status code
	if status := recorder.Code; status != statusCodeEsperado {
		t.Errorf("Status code incorreto (adquirido %v, esperado %v).",
			status, statusCodeEsperado)
	}

	// Checa o corpo da mensagem
	if recorder.Body.String() != respostaEsperada {
		t.Errorf("Corpo da mensagem incorreto.\nAdquirido: %v\nDesejado: %v",
			recorder.Body.String(), respostaEsperada)
	}
}

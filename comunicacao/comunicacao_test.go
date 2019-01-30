package comunicacao

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/loteny/redcoins/erros"
)

// ConteudoStruct serve como estrutura para ser passada nos testes de
// comunicação em JSON
type ConteudoStruct struct {
	Funcionando string
}

// TestResponde envia uma mensagem em HTTP, responde utilizando a função
// Responde e verifica seu status code, Content-Type e corpo JSON
func TestResponde(t *testing.T) {
	// Conteúdo JSON para enviar responder no request
	conteudo, err := json.Marshal(ConteudoStruct{"sim"})
	if err != nil {
		t.Fatal(err)
	}

	// Criação do request HTTP
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Recorder para armazenar a resposta
	recorder := httptest.NewRecorder()

	// Upload do servidor de teste
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Responde(w, http.StatusOK, conteudo)
	})
	handler.ServeHTTP(recorder, request)

	// Checa os status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status code incorreto (adquirido %v, esperado %v).",
			status, http.StatusOK)
	}

	// Checa o corpo da mensagem
	esperado := `{"Funcionando":"sim"}`
	if recorder.Body.String() != esperado {
		t.Errorf("Corpo da mensagem incorreto.\nAdquirido: %v\nDesejado: %v",
			recorder.Body.String(), esperado)
	}
}

// TestRespondeSucesso é idência à TestResponde
func TestRespondeSucesso(t *testing.T) {
	// Conteúdo JSON para enviar responder no request
	conteudo, err := json.Marshal(ConteudoStruct{"sim"})
	if err != nil {
		t.Fatal(err)
	}

	// Criação do request HTTP
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Recorder para armazenar a resposta
	recorder := httptest.NewRecorder()

	// Upload do servidor de teste
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RespondeSucesso(w, conteudo)
	})
	handler.ServeHTTP(recorder, request)

	// Checa os status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status code incorreto (adquirido %v, esperado %v).",
			status, http.StatusOK)
	}

	// Checa o corpo da mensagem
	esperado := `{"Funcionando":"sim"}`
	if recorder.Body.String() != esperado {
		t.Errorf("Corpo da mensagem incorreto.\nAdquirido: %v\nDesejado: %v",
			recorder.Body.String(), esperado)
	}
}

func TestRespondeErro(t *testing.T) {
	// Criação do request HTTP
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Recorder para armazenar a resposta
	recorder := httptest.NewRecorder()

	// Upload do servidor de teste
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RespondeErro(w, http.StatusBadRequest, erros.Cria(false, 400, "erro teste"))
	})
	handler.ServeHTTP(recorder, request)

	// Checa os status code
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("Status code incorreto (adquirido %v, esperado %v).",
			status, http.StatusBadRequest)
	}

	// Checa o corpo da mensagem
	esperado := `{"erros":["erro teste"]}`
	if recorder.Body.String() != esperado {
		t.Errorf("Corpo da mensagem incorreto.\nAdquirido: %v\nDesejado: %v",
			recorder.Body.String(), esperado)
	}
}

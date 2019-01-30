package cadastro

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/loteny/redcoins/erros"
)

// TestRealizaCadastroRequestHTTP por enquanto exige que o banco de dados de
// teste esteja vazio mas com as tabelas criadas
func TestRealizaCadastroRequestHTTP(t *testing.T) {
	// Formulário válido
	form := url.Values{}
	form.Set("email", "teste@gmail.com")
	form.Set("senha", "123456")
	form.Set("nome", "Ronnie James Dio")
	form.Set("nascimento", "1942-07-10")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP := func(w http.ResponseWriter, r *http.Request) {
		err := RealizaCadastroRequestHTTP(r)
		if !erros.Vazio(err) {
			t.Errorf("Erro inesperado no cadastro: %v", err)
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)

	// Nascimento inválido
	form.Set("email", "segundo@gmail.com")
	form.Set("nascimento", "194207-10")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		err := RealizaCadastroRequestHTTP(r)
		if err.Error() != ErrNascimentoInvalido.Error() {
			t.Errorf("Erro inesperado no cadastro: %v", err)
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)

	// Nascimento e senha inválidos
	form.Set("senha", "123")
	erroEsperado := ErrSenhaInvalida
	erroEsperado = erros.JuntaErros(ErrSenhaInvalida, ErrNascimentoInvalido)
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		err := RealizaCadastroRequestHTTP(r)
		if err.Error() != erroEsperado.Error() {
			t.Errorf("Erro inesperado no cadastro: %v\n%v", err, erroEsperado.Error())
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)
}

// TestVerificaLoginRequestHTTP por enquanto exige que
// TestRealizaCadastroRequestHTTP seja executado antes
func TestVerificaLoginRequestHTTP(t *testing.T) {
	// Conta válida
	form := url.Values{}
	form.Set("email", "teste@gmail.com")
	form.Set("senha", "123456")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP := func(w http.ResponseWriter, r *http.Request) {
		logado, err := VerificaLoginRequestHTTP(r)
		if !erros.Vazio(err) {
			t.Errorf("Erro inesperado no login: %v", err)
		} else if !logado {
			t.Errorf("Usuário não foi logado quando deveria ter sido.")
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)

	// Senha incorreta
	form.Set("email", "teste@gmail.com")
	form.Set("senha", "123455")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		logado, err := VerificaLoginRequestHTTP(r)
		if !erros.Vazio(err) {
			t.Errorf("Erro inesperado no login: %v", err)
		} else if logado {
			t.Errorf("Usuário foi logado quando não deveria ter sido.")
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)

	// Conta não existente
	form.Set("email", "email-nao-cadastrado@gmail.com")
	form.Set("senha", "123456")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		logado, err := VerificaLoginRequestHTTP(r)
		if !erros.Vazio(err) {
			t.Errorf("Erro inesperado no login: %v", err)
		} else if logado {
			t.Errorf("Usuário foi logado quando não deveria ter sido.")
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

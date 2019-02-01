package cadastro

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/loteny/redcoins/database"
	"github.com/loteny/redcoins/erros"
	"github.com/loteny/redcoins/passenc"
)

// init deleta o banco de dados e cria um novo apenas com alguns usuários para
// testes
func init() {
	if err := database.DeletaDatabaseTeste(); err != nil {
		log.Fatalf("erro ao deletar database: %s", err)
	}
	if err := database.CriaDatabase(); err != nil {
		log.Fatalf("erro ao criar database: %s", err)
	}
	if err := testPopulaDatabase(); err != nil {
		log.Fatalf("erro ao popular banco de dados: %s", err)
	}
}

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
		// Verifica se o usuário foi cadastrado vendo a senha dele no banco
		// de dados
		senha, err2 := database.AdquireSenhaHashed("teste@gmail.com")
		if err2 != nil {
			t.Errorf("Erro inesperado na verificação de cadastro: %v", err2)
		}
		if sucesso, err2 := passenc.VerificaSenha([]byte("123456"), senha); err2 != nil {
			t.Errorf("Erro inesperado na verificação de cadastro: %v", err2)
		} else if !sucesso {
			t.Errorf("Usuário não foi cadastrado corretamente. Senha hash recebida: %v", senha)
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

func TestVerificaLoginRequestHTTP(t *testing.T) {
	// Usuário e senhas corretos - autenticação bem-sucedida
	form := url.Values{}
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP := func(w http.ResponseWriter, r *http.Request) {
		logado, _, err := VerificaLoginRequestHTTP(r)
		if !erros.Vazio(err) {
			t.Errorf("Erro inesperado no login: %v", err)
		} else if !logado {
			t.Errorf("Usuário não foi logado quando deveria ter sido.")
		}
	}
	testRealizaRequestHTTPPostFormAuth(t, form, rotaHTTP, "valido1@gmail.com", "senhavalido1")

	// Sem autenticação
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		logado, _, err := VerificaLoginRequestHTTP(r)
		if !erros.Vazio(err) {
			t.Errorf("Erro inesperado no login: %v", err)
		} else if logado {
			t.Errorf("Usuário foi logado quando não deveria ter sido (sem autenticação).")
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)

	// Conta inexistente
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		logado, _, err := VerificaLoginRequestHTTP(r)
		if !erros.Vazio(err) {
			t.Errorf("Erro inesperado no login: %v", err)
		} else if logado {
			t.Errorf("Usuário foi logado quando não deveria ter sido (conta inexistente).")
		}
	}
	testRealizaRequestHTTPPostFormAuth(t, form, rotaHTTP, "email-nao-cadastrado@gmail.com", "123456")

	// Senha incorreta
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		logado, _, err := VerificaLoginRequestHTTP(r)
		if !erros.Vazio(err) {
			t.Errorf("Erro inesperado no login: %v", err)
		} else if logado {
			t.Errorf("Usuário foi logado quando não deveria ter sido (senha incorreta).")
		}
	}
	testRealizaRequestHTTPPostFormAuth(t, form, rotaHTTP, "valido1@gmail.com", "senhaincorreta")
}

// testRealizaRequestHTTPPostForm é uma função auxiliar para geração de requests
// HTTP com formulário POST. A função 'f' vai receber e processar o request
// HTTP.
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

// testRealizaRequestHTTPPostFormAuth é idência à
// testRealizaRequestHTTPPostForm, mas utiliza HTTP Basic Auth
func testRealizaRequestHTTPPostFormAuth(t *testing.T, form url.Values,
	f func(w http.ResponseWriter, r *http.Request), usuario string, senha string) {
	// Criação do request HTTP
	request, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(usuario, senha)

	// Recorder para armazenar a resposta
	recorder := httptest.NewRecorder()

	// Upload do servidor de teste
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(recorder, request)
}

// testPopulaDatabase insere dados no banco de dados necessários para a
// realização dos testes
func testPopulaDatabase() error {
	senha, err := passenc.GeraHashed([]byte("senhavalido1"))
	if err != nil {
		return err
	}
	usr := database.Usuario{
		Email:      "valido1@gmail.com",
		Senha:      senha,
		Nome:       "Conta Válida 1",
		Nascimento: "1994-03-07",
	}
	return database.InsereUsuario(&usr)
}

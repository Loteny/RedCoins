package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/loteny/redcoins/database"
	"github.com/loteny/redcoins/passenc"
)

// init deleta o banco de dados e cria um novo apropriadamente populado
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

func TestRotaCadastro(t *testing.T) {
	// Cadastro válido
	form := url.Values{}
	form.Set("email", "testerotacadastro@gmail.com")
	form.Set("senha", "123456")
	form.Set("nome", "Teste Rota Cadastro")
	form.Set("nascimento", "1994-03-07")
	statusCode, body := testPostSimples(t, form, RotaCadastro)
	if statusCode != 201 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != "" {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}

	// Usuário duplicado
	statusCode, body = testPostSimples(t, form, RotaCadastro)
	if statusCode != 400 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `{"erros":["email_ja_cadastrado"]}` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}
}

func TestRotaCompra(t *testing.T) {
	// Compra válida
	form := url.Values{}
	form.Set("qtd", "0.0001")
	form.Set("data", "2000-01-01")
	statusCode, body := testPostAuth(t, form, RotaCompra, "valido4@gmail.com", "senhavalido4")
	if statusCode != 201 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != "" {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}

	// Autenticação falha
	statusCode, body = testPostAuth(t, form, RotaCompra, "valido4@gmail.com", "senhaincorreta")
	if statusCode != 403 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}

	// Data inválida
	form.Set("data", "200001-01")
	statusCode, body = testPostAuth(t, form, RotaCompra, "valido4@gmail.com", "senhavalido4")
	if statusCode != 400 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `{"erros":["data_invalida"]}` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}
}

func TestRotaVenda(t *testing.T) {
	// Venda válida
	form := url.Values{}
	form.Set("qtd", "0.0001")
	form.Set("data", "2000-01-01")
	statusCode, body := testPostAuth(t, form, RotaVenda, "valido3@gmail.com", "senhavalido3")
	if statusCode != 201 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != "" {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}

	// Autenticação falha
	statusCode, body = testPostAuth(t, form, RotaVenda, "valido3@gmail.com", "senhaincorreta")
	if statusCode != 403 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}

	// Data inválida
	form.Set("data", "200001-01")
	statusCode, body = testPostAuth(t, form, RotaVenda, "valido3@gmail.com", "senhavalido3")
	if statusCode != 400 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `{"erros":["data_invalida"]}` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}

	// Quantidade inválida
	form.Set("qtd", "0,4,4")
	form.Set("data", "2000-01-01")
	statusCode, body = testPostAuth(t, form, RotaVenda, "valido3@gmail.com", "senhavalido3")
	if statusCode != 400 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `{"erros":["qtd_invalida"]}` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}

	// Saldo insuficiente
	form.Set("qtd", "10000")
	statusCode, body = testPostAuth(t, form, RotaVenda, "valido3@gmail.com", "senhavalido3")
	if statusCode != 400 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `{"erros":["saldo_insuficiente"]}` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}
}

func TestRotaRelatorioDia(t *testing.T) {
	// Verificaremos as transações do dia 2018-01-02
	dados := map[string]string{"data": "2018-01-02"}
	statusCode, body := testGetAuth(t, dados, RotaRelatorioDia, "valido1@gmail.com", "senhavalido1")
	if statusCode != 200 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `{"transacoes":[`+
		`{"usuario":"valido1@gmail.com","compra":true,"creditos":30,"bitcoins":0.003,"dia":"2018-01-02"},`+
		`{"usuario":"valido2@gmail.com","compra":true,"creditos":40,"bitcoins":0.004,"dia":"2018-01-02"},`+
		`{"usuario":"valido2@gmail.com","compra":false,"creditos":15,"bitcoins":0.0005,"dia":"2018-01-02"}]}` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}

	// Autenticação falha
	statusCode, body = testGetAuth(t, dados, RotaRelatorioDia, "valido1@gmail.com", "senhaincorreta")
	if statusCode != 403 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}
}

func TestRotaRelatorioUsuario(t *testing.T) {
	// Verificaremos as transações do usuário valido1@gmail.com
	dados := map[string]string{"email": "valido1@gmail.com"}
	statusCode, body := testGetAuth(t, dados, RotaRelatorioUsuario, "valido1@gmail.com", "senhavalido1")
	if statusCode != 200 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `{"transacoes":[`+
		`{"usuario":"valido1@gmail.com","compra":true,"creditos":10,"bitcoins":0.001,"dia":"2018-01-01"},`+
		`{"usuario":"valido1@gmail.com","compra":true,"creditos":20,"bitcoins":0.002,"dia":"2018-01-01"},`+
		`{"usuario":"valido1@gmail.com","compra":true,"creditos":30,"bitcoins":0.003,"dia":"2018-01-02"}]}` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}

	// Autenticação falha
	statusCode, body = testGetAuth(t, dados, RotaRelatorioUsuario, "valido1@gmail.com", "senhaincorreta")
	if statusCode != 403 {
		t.Errorf("Status code inesperado: %v", statusCode)
	}
	if body != `` {
		t.Errorf("Corpo da resposta inesperado: %v", body)
	}
}

// testPostSimples realiza um request HTTP POST que a função 'f' vai receber.
// Retorna o status code e o corpo da mensagem de resposta.
func testPostSimples(t *testing.T, form url.Values, f func(w http.ResponseWriter, r *http.Request)) (int, string) {
	request, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Armazena a resposta da rota
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(recorder, request)

	result := recorder.Result()
	statusCode := result.StatusCode
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return statusCode, string(body)
}

// testPostAuth realiza um request HTTP POST que a função 'f' vai receber.
// Retorna o status code e o corpo da mensagem de resposta. Utiliza HTTP Auth
// Basic com os argumentos 'usuario' e 'senha'.
func testPostAuth(t *testing.T, form url.Values, f func(w http.ResponseWriter, r *http.Request), usuario string, senha string) (int, string) {
	request, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(usuario, senha)

	// Armazena a resposta da rota
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(recorder, request)

	result := recorder.Result()
	statusCode := result.StatusCode
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return statusCode, string(body)
}

// testGetAuth realiza um request HTTP GET que a função 'f' vai receber.
// Retorna o status code e o corpo da mensagem de resposta. Utiliza HTTP Auth
// Basic com os argumentos 'usuario' e 'senha'.
func testGetAuth(t *testing.T, dados map[string]string,
	f func(w http.ResponseWriter, r *http.Request),
	usuario string, senha string) (int, string) {
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetBasicAuth(usuario, senha)
	q := request.URL.Query()
	for k, v := range dados {
		q.Add(k, v)
	}
	request.URL.RawQuery = q.Encode()

	// Armazena a resposta da rota
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(recorder, request)

	result := recorder.Result()
	statusCode := result.StatusCode
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return statusCode, string(body)
}

// testPopulaDatabase insere dados no banco de dados necessários para a
// realização dos testes. Essa função cria no total:
// - 4 usuários
// - 2 compras no mesmo dia para o primeiro usuário
// - 1 compra em um outro dia para o primeiro usuário
// - 1 compra no mesmo dia que a anterior para o segundo usuário
// - 1 venda no mesmo dia que as duas compras anteriores para o primeiro usuário
// - 1 compra em um dia "irrelevante" para o terceiro usuário de 1 BTC
func testPopulaDatabase() error {
	// Usuário 1
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
	if err := database.InsereUsuario(&usr); err != nil {
		return err
	}

	// Usuário 2
	senha, err = passenc.GeraHashed([]byte("senhavalido2"))
	if err != nil {
		return err
	}
	usr = database.Usuario{
		Email:      "valido2@gmail.com",
		Senha:      senha,
		Nome:       "Conta Válida 2",
		Nascimento: "1994-03-08",
	}
	if err := database.InsereUsuario(&usr); err != nil {
		return err
	}

	// Usuário 3
	senha, err = passenc.GeraHashed([]byte("senhavalido3"))
	if err != nil {
		return err
	}
	usr = database.Usuario{
		Email:      "valido3@gmail.com",
		Senha:      senha,
		Nome:       "Conta Válida 3",
		Nascimento: "1994-03-09",
	}
	if err := database.InsereUsuario(&usr); err != nil {
		return err
	}

	// Usuário 4
	senha, err = passenc.GeraHashed([]byte("senhavalido4"))
	if err != nil {
		return err
	}
	usr = database.Usuario{
		Email:      "valido4@gmail.com",
		Senha:      senha,
		Nome:       "Conta Válida 4",
		Nascimento: "1994-03-10",
	}
	if err := database.InsereUsuario(&usr); err != nil {
		return err
	}

	// Compras
	if err := database.InsereTransacao("valido1@gmail.com", true, 0.001, 10, "2018-01-01"); err != nil {
		return err
	}
	if err := database.InsereTransacao("valido1@gmail.com", true, 0.002, 20, "2018-01-01"); err != nil {
		return err
	}
	if err := database.InsereTransacao("valido1@gmail.com", true, 0.003, 30, "2018-01-02"); err != nil {
		return err
	}
	if err := database.InsereTransacao("valido2@gmail.com", true, 0.004, 40, "2018-01-02"); err != nil {
		return err
	}
	if err := database.InsereTransacao("valido3@gmail.com", true, 1, 400, "2012-01-02"); err != nil {
		return err
	}
	// Venda
	if err := database.InsereTransacao("valido2@gmail.com", false, 0.0005, 15, "2018-01-02"); err != nil {
		return err
	}

	return nil
}

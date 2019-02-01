package transacao

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

func TestCompraHTTP(t *testing.T) {
	// Resultado esperado dos testes
	esperado := database.Transacao{
		Usuario:  "valido4@gmail.com",
		Compra:   true,
		Bitcoins: 0.03,
		Dia:      "2015-01-01",
	}

	// Formulário válido
	form := url.Values{}
	form.Set("qtd", "0.03")
	form.Set("data", esperado.Dia)
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP := func(w http.ResponseWriter, r *http.Request) {
		err := CompraHTTP(r, esperado.Usuario)
		if !erros.Vazio(err) {
			t.Fatalf("Erro inesperado na transação: %v", err)
		}
		trs, err2 := database.AdquireTransacoesDeUsuario(esperado.Usuario)
		if err2 != nil {
			t.Fatalf("Erro inesperado na transação: %v", err2)
		}
		if len(trs) != 1 {
			t.Errorf("Mais de uma transação encontrada quando deveria haver uma: %v", trs)
		}
		tr := trs[0]
		if tr.Usuario != esperado.Usuario ||
			tr.Compra != esperado.Compra ||
			tr.Bitcoins != esperado.Bitcoins ||
			tr.Dia != esperado.Dia {
			t.Errorf("Dados da transação incorretos: %v", tr)
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)

	// Quantidade inválida
	form.Set("qtd", "-2")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		err := CompraHTTP(r, esperado.Usuario)
		if err.Error() != ErrQtdInvalida.Error() {
			t.Fatalf("Erro inesperado na transação: %v", err)
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)
}

func TestVendaHTTP(t *testing.T) {
	// Resultado esperado dos testes
	esperado := database.Transacao{
		Usuario:  "valido3@gmail.com",
		Compra:   false,
		Bitcoins: 0.0001,
		Dia:      "2012-01-01",
	}

	// Formulário válido
	form := url.Values{}
	form.Set("qtd", "0.0001")
	form.Set("data", esperado.Dia)
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP := func(w http.ResponseWriter, r *http.Request) {
		err := VendaHTTP(r, esperado.Usuario)
		if !erros.Vazio(err) {
			t.Fatalf("Erro inesperado na transação: %v", err)
		}
		trs, err2 := database.AdquireTransacoesDeUsuario(esperado.Usuario)
		if err2 != nil {
			t.Fatalf("Erro inesperado na transação: %v", err2)
		}
		if len(trs) != 2 {
			t.Errorf("Mais de duas transação encontrada quando deveria haver duas: %v", trs)
		}
		tr := trs[1]
		if tr.Usuario != esperado.Usuario ||
			tr.Compra != esperado.Compra ||
			tr.Bitcoins != esperado.Bitcoins ||
			tr.Dia != esperado.Dia {
			t.Errorf("Dados da transação incorretos: %v", tr)
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)

	// Quantidade inválida
	form.Set("qtd", "-2")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		err := VendaHTTP(r, esperado.Usuario)
		if err.Error() != ErrQtdInvalida.Error() {
			t.Fatalf("Erro inesperado na transação: %v", err)
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)

	// Saldo insuficiente
	form.Set("qtd", "1000")
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP = func(w http.ResponseWriter, r *http.Request) {
		err := VendaHTTP(r, esperado.Usuario)
		if err.Error() != ErrSaldoInsuficiente.Error() {
			t.Fatalf("Erro inesperado na transação: %v", err)
		}
	}
	testRealizaRequestHTTPPostForm(t, form, rotaHTTP)
}

func TestTransacoesDiaHTTP(t *testing.T) {
	// Usamos o dia 2018-01-02 para verificar as transações que conhecemos
	dados := map[string]string{"data": "2018-01-02"}
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP := func(w http.ResponseWriter, r *http.Request) {
		resp, err := TransacoesDiaHTTP(r)
		if !erros.Vazio(err) {
			t.Fatalf("Erro inesperado na transação: %v", err)
		}
		respostaEsperada := `{"transacoes":[{"usuario":"valido1@gmail.com","compra":true,"creditos":30,"bitcoins":0.003,"dia":"2018-01-02"},{"usuario":"valido2@gmail.com","compra":true,"creditos":40,"bitcoins":0.004,"dia":"2018-01-02"},{"usuario":"valido2@gmail.com","compra":false,"creditos":15,"bitcoins":0.0005,"dia":"2018-01-02"}]}`
		if string(resp) != respostaEsperada {
			t.Errorf("Lista de transações incorreta: %v", string(resp))
		}
	}
	testRealizaRequestHTTPGetForm(t, dados, rotaHTTP)
}

func TestTransacoesUsuarioHTTP(t *testing.T) {
	// Usamos o usuário valido1@gmail.com para verificar as transações que conhecemos
	dados := map[string]string{"email": "valido1@gmail.com"}
	// Função que vai chamar a função a ser testada e tratar seu retorno
	rotaHTTP := func(w http.ResponseWriter, r *http.Request) {
		resp, err := TransacoesUsuarioHTTP(r)
		if !erros.Vazio(err) {
			t.Fatalf("Erro inesperado na transação: %v", err)
		}
		respostaEsperada := `{"transacoes":[{"usuario":"valido1@gmail.com","compra":true,"creditos":10,"bitcoins":0.001,"dia":"2018-01-01"},{"usuario":"valido1@gmail.com","compra":true,"creditos":20,"bitcoins":0.002,"dia":"2018-01-01"},{"usuario":"valido1@gmail.com","compra":true,"creditos":30,"bitcoins":0.003,"dia":"2018-01-02"}]}`
		if string(resp) != respostaEsperada {
			t.Errorf("Lista de transações incorreta: %v", string(resp))
		}
	}
	testRealizaRequestHTTPGetForm(t, dados, rotaHTTP)
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

// testRealizaRequestHTTPGetForm é uma função auxiliar para geração de requests
// HTTP com formulário GET
func testRealizaRequestHTTPGetForm(t *testing.T, dados map[string]string,
	f func(w http.ResponseWriter, r *http.Request)) {
	// Criação do request HTTP
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := request.URL.Query()
	for k, v := range dados {
		q.Add(k, v)
	}
	request.URL.RawQuery = q.Encode()

	// Recorder para armazenar a resposta
	recorder := httptest.NewRecorder()

	// Upload do servidor de teste
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(recorder, request)
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

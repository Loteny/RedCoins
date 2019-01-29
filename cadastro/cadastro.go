// Package cadastro trata de todas as funções para cadastramento do cliente,
// desde a recepção do request HTTPS até a inserção dos dados no banco de dados
// e verificação dos credenciais para autenticação.
// Esse package usa exclusivamente erros.Erros como estrutura de erros.
package cadastro

import (
	"net/http"

	"github.com/loteny/redcoins/comunicacao"
	"github.com/loteny/redcoins/erros"
)

// Lista de possíveis erros do módulo
var (
	ErrMetodoPost         = erros.Cria(false, 405, "")
	ErrEmailInvalido      = erros.Cria(false, 400, "email invalido")
	ErrSenhaInvalida      = erros.Cria(false, 400, "senha invalida")
	ErrSenhaMuitoLonga    = erros.Cria(false, 400, "senha muito longa")
	ErrNomeInvalido       = erros.Cria(false, 400, "nome invalido")
	ErrNascimentoInvalido = erros.Cria(false, 400, "nascimento invalido")
)

// Estrutura que contém todos os dados cadastrais de um usuário
type dadosCadastrais struct {
	email      string
	senha      string
	nome       string
	nascimento string
}

// CadastraHTTPS realiza o cadastro de um usuário a partir de um request HTTPS
func CadastraHTTPS(w http.ResponseWriter, r *http.Request) {
	if _, err := validaDadosCadastroRequestHTTP(r); err != nil {
		interno, status, erroNativo := err.(erros.Erros).Abre()
		if !interno {
			comunicacao.RespondeErro(w, status, erroNativo)
			return
		}
		comunicacao.Responde(w, status, []byte{})
		return
	}
	comunicacao.RespondeSucesso(w, []byte{})
}

// validaDadosCadastroRequestHTTP valida os dados cadastrais apropriados
// recebidos no request HTTP. O request deve ser do tipo POST, caso contrário,
// ocorrerá o erro HTTP de status code 405. Após adquirir os dados do request, a
// função chama validaDadosCadastro para validar e limpar os dados
func validaDadosCadastroRequestHTTP(r *http.Request) (dados dadosCadastrais, err error) {
	// Verifica o método do request
	if r.Method != "POST" {
		return dadosCadastrais{}, ErrMetodoPost
	}
	// Adquire os dados do request
	comunicacao.RealizaParseForm(r)
	dados.email = r.PostFormValue("email")
	dados.senha = r.PostFormValue("senha")
	dados.nome = r.PostFormValue("nome")
	dados.nascimento = r.PostFormValue("nascimento")
	// Validação e retorno
	err = validaDadosCadastro(&dados)
	return
}

// validaDadosCadastro verifica todos os dados recebidos em 'dados'. Se for
// possível formatar corretamente e tornar os dados válidos, isso será feito. Se
// houver algum erro incorrigível com os dados, a função retorna o status code e
// mensagem de erros apropriados.
func validaDadosCadastro(dados *dadosCadastrais) (err error) {
	if err := email(dados.email); err != nil {
		return err
	}
	if err := senha(dados.senha); err != nil {
		return err
	}
	if err := nome(dados.nome); err != nil {
		return err
	}
	if err := nascimento(dados.nascimento); err != nil {
		return err
	}
	return
}

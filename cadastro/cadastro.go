// Package cadastro trata de todas as funções para cadastramento do cliente,
// desde a recepção do request HTTPS até a inserção dos dados no banco de dados
// e verificação dos credenciais para autenticação
package cadastro

import (
	"errors"
	"net/http"

	"github.com/loteny/redcoins/comunicacao"
)

// Lista de possíveis erros do módulo
var (
	ErrMetodoPost = errors.New("405 - O método deve ser POST")
)

type dadosCadastrais struct {
	email      string
	senha      string
	nome       string
	nascimento string
}

// CadastraHTTPS realiza o cadastro de um usuário a partir de um request HTTPS
func CadastraHTTPS(w http.ResponseWriter, r *http.Request) {
	if _, status, err := validaDadosCadastroRequestHTTP(r); err != nil {
		comunicacao.Responde(w, status, []byte{})
		return
	}
	comunicacao.RespondeSucesso(w, []byte{})
}

// validaDadosCadastroRequestHTTP valida os dados cadastrais apropriados
// recebidos no request HTTP. O request deve ser do tipo POST, caso contrário,
// ocorrerá o erro HTTP de status code 405. Após adquirir os dados do request, a
// função chama validaDadosCadastro para validar e limpar os dados
func validaDadosCadastroRequestHTTP(r *http.Request) (dados dadosCadastrais, status int, err error) {
	// Verifica o método do request
	if r.Method != "POST" {
		return dadosCadastrais{}, 405, ErrMetodoPost
	}
	// Adquire os dados do request
	comunicacao.RealizaParseForm(r)
	dados.email = r.PostFormValue("email")
	dados.senha = r.PostFormValue("senha")
	dados.nome = r.PostFormValue("nome")
	dados.nascimento = r.PostFormValue("nascimento")
	// Validação e retorno
	status, err = validaDadosCadastro(&dados)
	return
}

// validaDadosCadastro verifica todos os dados recebidos em 'dados'. Se for
// possível formatar corretamente e tornar os dados válidos, isso será feito. Se
// houver algum erro incorrigível com os dados, a função retorna o status code e
// mensagem de erros apropriados.
func validaDadosCadastro(dados *dadosCadastrais) (status int, err error) {
	return
}

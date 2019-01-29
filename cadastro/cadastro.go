// Package cadastro trata de todas as funções para cadastramento do cliente,
// desde a recepção do request HTTPS até a inserção dos dados no banco de dados
// e verificação dos credenciais para autenticação.
// Esse package usa exclusivamente erros.Erros como estrutura de erros.
package cadastro

import (
	"net/http"

	"github.com/loteny/redcoins/comunicacao"
	"github.com/loteny/redcoins/database"
	"github.com/loteny/redcoins/erros"
	"github.com/loteny/redcoins/passenc"
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

// RealizaCadastroRequestHTTP realiza o cadastro de um usuário a partir de um
// request HTTP. O request e os dados do usuário serão validados com a função
// ValidaDadosCadastroRequestHTTP.
func RealizaCadastroRequestHTTP(r *http.Request) error {
	// Valida o usuário
	dados, err := validaDadosCadastroRequestHTTP(r)
	if err != nil {
		return err
	}
	// Gera a senha hashed
	senhaHashed, err := passenc.GeraHashed([]byte(dados.senha))
	if err != nil {
		return err
	}
	// Insere o usuário no banco de dados
	usr := database.Usuario{
		Email:      dados.email,
		Senha:      senhaHashed,
		Nome:       dados.nome,
		Nascimento: dados.nascimento,
	}
	if err := database.InsereUsuario(&usr); err != nil {
		return err
	}
	return nil
}

// VerificaLoginRequestHTTP verifica se o usuário existe e a senha está correta
// a partir de um request HTTP. O request deve ser do tipo POST.
func VerificaLoginRequestHTTP(r *http.Request) (bool, error) {
	// Verifica o método do request
	if r.Method != "POST" {
		return false, ErrMetodoPost
	}
	// Adquire os dados do request
	if err := comunicacao.RealizaParseForm(r); err != nil {
		return false, err
	}
	email := r.PostFormValue("email")
	senha := r.PostFormValue("senha")

	return verificaLogin(email, senha)
}

// ValidaDadosCadastroRequestHTTP valida os dados cadastrais apropriados
// recebidos no request HTTP. O request deve ser do tipo POST, caso contrário,
// ocorrerá o erro HTTP de status code 405. Após adquirir os dados do request, a
// função chama validaDadosCadastro para validar e limpar os dados
func validaDadosCadastroRequestHTTP(r *http.Request) (dados dadosCadastrais, err error) {
	// Verifica o método do request
	if r.Method != "POST" {
		return dadosCadastrais{}, ErrMetodoPost
	}
	// Adquire os dados do request
	if err := comunicacao.RealizaParseForm(r); err != nil {
		return dadosCadastrais{}, err
	}
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

// verificaLogin verifica se os credenciais existem e estão corretos no banco
// de dados utilizando encriptação de senhas
func verificaLogin(email string, senha string) (bool, error) {
	// Adquire a senha hashed do banco de dados
	senhaDB, err := database.AdquireSenhaHashed(email)
	if err == database.ErrUsuarioNaoExiste {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Verifica se as senhas são as mesmas
	sucesso, err := passenc.VerificaSenha([]byte(senha), senhaDB)
	if err != nil {
		return false, err
	}
	return sucesso, nil
}

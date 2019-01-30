package main

import (
	"net/http"

	"github.com/loteny/redcoins/cadastro"
	"github.com/loteny/redcoins/comunicacao"
	"github.com/loteny/redcoins/erros"
	"github.com/loteny/redcoins/transacao"
)

// RotaCadastro realiza o cadastro de um usuário a partir de um request HTTPS.
// O pedido deve ser feito com o método POST e ter os campos "nome", "senha",
// "nascimento" e "email" preenchidos. A validação de cada um desses campos
// está definida no módulo 'cadastro', no arquivo 'campos.go'.
func RotaCadastro(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		comunicacao.Responde(w, http.StatusMethodNotAllowed, []byte{})
		return
	}
	respostaPadrao(w, r, http.StatusCreated, cadastro.RealizaCadastroRequestHTTP)
}

// RotaCompra realiza a compra de Bitcoins para um usuário a partir de um
// request HTTPS. O pedido deve ser feito com o método POST e ter os campos
// "email", "senha" e "qtd" preenchidos, sendo o último a quantidade de Bitcoins
// a ser comprada, apenas dígitos e com o separado decimal sendo ponto.
func RotaCompra(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		comunicacao.Responde(w, http.StatusMethodNotAllowed, []byte{})
		return
	}
	respostaPadraoAutenticada(w, r, http.StatusCreated, transacao.CompraHTTP)
}

// RotaVenda realiza a venda de Bitcoins para um usuário a partir de um request
// HTTPS. O pedido deve ser feito com o método POST e ter os campos "email",
// "senha" e "qtd" preenchidos, sendo o último a quantidade de Bitcoins a ser
// comprada, apenas dígitos e com o separado decimal sendo ponto.
func RotaVenda(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		comunicacao.Responde(w, http.StatusMethodNotAllowed, []byte{})
		return
	}
	respostaPadraoAutenticada(w, r, http.StatusCreated, transacao.VendaHTTP)
}

// respostaPadrao chamada a função 'f' e envia a resposta HTTP adequada
// dependendo do resultado da função, considerando que a função retorna um
// erros.Erros.
// Se o erro for interno, apenas o statusCode deve ser enviado para o cliente,
// indicando um erro no servidor sem dar detalhes de seu funcionamento. Se não
// for, os erros devem ser enviados para cliente.
// statusSucesso indica o status code HTTP para caso de sucesso.
func respostaPadrao(w http.ResponseWriter, r *http.Request, statusSucesso int, f func(*http.Request) erros.Erros) {
	if err := f(r); !erros.Vazio(err) {
		interno, status, _ := erros.Abre(err)
		if !interno {
			comunicacao.RespondeErro(w, status, err)
			return
		}
		comunicacao.Responde(w, status, []byte{})
		return
	}
	comunicacao.Responde(w, http.StatusCreated, []byte{})
}

// respostaPadraoAutenticada é idêntica à respostaPadrao, mas autentica o
// usuário com autenticaUsuarioPost antes de proceder às operações
func respostaPadraoAutenticada(w http.ResponseWriter, r *http.Request, statusSucesso int, f func(*http.Request) erros.Erros) {
	if autenticado, err := autenticaUsuarioPost(r); !erros.Vazio(err) {
		interno, status, _ := erros.Abre(err)
		if !interno {
			comunicacao.RespondeErro(w, status, err)
			return
		}
		comunicacao.Responde(w, status, []byte{})
		return
	} else if !autenticado {
		comunicacao.Responde(w, http.StatusForbidden, []byte{})
		return
	}
	respostaPadrao(w, r, statusSucesso, f)
}

// autenticaUsuarioPost verifica se o usuário é cadastrado e se a senha está
// correta a partir dos campos POST "email" e "senha"
func autenticaUsuarioPost(r *http.Request) (bool, erros.Erros) {
	return cadastro.VerificaLoginRequestHTTP(r)
}

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
	if err := cadastro.RealizaCadastroRequestHTTP(r); err != nil {
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

// RotaCompra realiza a compra de Bitcoins para um usuário a partir de um
// request HTTPS. O pedido deve ser feito com o método POST e ter os campos
// "email", "senha" e "qtd" preenchidos, sendo o último a quantidade de Bitcoins
// a ser comprada, apenas dígitos e com o separado decimal sendo ponto.
func RotaCompra(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		comunicacao.Responde(w, http.StatusMethodNotAllowed, []byte{})
		return
	}
	if err := transacao.CompraHTTP(r); err != nil {
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

// RotaVenda realiza a venda de Bitcoins para um usuário a partir de um request
// HTTPS. O pedido deve ser feito com o método POST e ter os campos "email",
// "senha" e "qtd" preenchidos, sendo o último a quantidade de Bitcoins a ser
// comprada, apenas dígitos e com o separado decimal sendo ponto.
func RotaVenda(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		comunicacao.Responde(w, http.StatusMethodNotAllowed, []byte{})
		return
	}
	if err := transacao.VendaHTTP(r); err != nil {
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

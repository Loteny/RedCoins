package main

import (
	"net/http"

	"github.com/loteny/redcoins/cadastro"
	"github.com/loteny/redcoins/comunicacao"
	"github.com/loteny/redcoins/erros"
)

// RotaCadastro realiza o cadastro de um usuário a partir de um request HTTPS.
// O pedido deve ser feito com o método POST e ter os campos "nome", "senha",
// "nascimento" e "email" preenchidos. A validação de cada um desses campos
// está definida no módulo 'cadastro', no arquivo 'campos.go'.
func RotaCadastro(w http.ResponseWriter, r *http.Request) {
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

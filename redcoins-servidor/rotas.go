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

	// Se o erro gerado for externo, envia os erros para o usuário. Se não,
	// envia apenas um status code indicando erro interno. Se não houve erro
	// gerado, envia status code de sucesso.
	if err := cadastro.RealizaCadastroRequestHTTP(r); !erros.Vazio(err) {
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

// RotaCompra realiza a compra de Bitcoins para um usuário a partir de um
// request HTTPS. O pedido deve ser feito com o método POST e ter os campos
// "email", "senha", "qtd" e "data" preenchidos, sendo "qtd" a quantidade de
// Bitcoins a ser comprada, apenas dígitos e com o separado decimal sendo ponto
// e "data" no formato "YYYY-MM-DD".
func RotaCompra(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		comunicacao.Responde(w, http.StatusMethodNotAllowed, []byte{})
		return
	}

	autenticado, email := autenticaUsuario(w, r)
	if !autenticado {
		return
	}

	// Se o erro gerado for externo, envia os erros para o usuário. Se não,
	// envia apenas um status code indicando erro interno. Se não houve erro
	// gerado, envia status code de sucesso.
	if err := transacao.CompraHTTP(r, email); !erros.Vazio(err) {
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

// RotaVenda realiza a venda de Bitcoins para um usuário a partir de um request
// HTTPS. O pedido deve ser feito com o método POST e ter os campos "email",
// "senha", "qtd" e "data" preenchidos, sendo "qtd" a quantidade de Bitcoins a
// ser vendida, apenas dígitos e com o separado decimal sendo ponto e "data" no
// formato "YYYY-MM-DD".
func RotaVenda(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		comunicacao.Responde(w, http.StatusMethodNotAllowed, []byte{})
		return
	}

	autenticado, email := autenticaUsuario(w, r)
	if !autenticado {
		return
	}

	// Se o erro gerado for externo, envia os erros para o usuário. Se não,
	// envia apenas um status code indicando erro interno. Se não houve erro
	// gerado, envia status code de sucesso.
	if err := transacao.VendaHTTP(r, email); !erros.Vazio(err) {
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

// RotaRelatorioDia retorna todas as transações feitas em um determinado dia
// no campo "data" (YYYY-MM-DD)
func RotaRelatorioDia(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		comunicacao.Responde(w, http.StatusMethodNotAllowed, []byte{})
		return
	}

	autenticado, _ := autenticaUsuario(w, r)
	if !autenticado {
		return
	}

	// Se o erro gerado for externo, envia os erros para o usuário. Se não,
	// envia apenas um status code indicando erro interno. Se não houve erro
	// gerado, envia status code de sucesso.
	resposta, err := transacao.TransacoesDiaHTTP(r)
	if !erros.Vazio(err) {
		interno, status, _ := erros.Abre(err)
		if !interno {
			comunicacao.RespondeErro(w, status, err)
			return
		}
		comunicacao.Responde(w, status, []byte{})
		return
	}
	comunicacao.Responde(w, http.StatusOK, resposta)
}

// RotaRelatorioUsuario retorna todas as transações feitas em um determinado
// usuário a partir de seu e-mail no campo "email"
func RotaRelatorioUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		comunicacao.Responde(w, http.StatusMethodNotAllowed, []byte{})
		return
	}

	autenticado, _ := autenticaUsuario(w, r)
	if !autenticado {
		return
	}

	// Se o erro gerado for externo, envia os erros para o usuário. Se não,
	// envia apenas um status code indicando erro interno. Se não houve erro
	// gerado, envia status code de sucesso.
	resposta, err := transacao.TransacoesUsuarioHTTP(r)
	if !erros.Vazio(err) {
		interno, status, _ := erros.Abre(err)
		if !interno {
			comunicacao.RespondeErro(w, status, err)
			return
		}
		comunicacao.Responde(w, status, []byte{})
		return
	}
	comunicacao.Responde(w, http.StatusOK, resposta)
}

// autenticaUsuario verifica se o usuário é cadastrado e se a senha está
// correta utilizando Basic Auth e retorna, também, seu e-mail. Em caso de erro
// de autenticação, a função responde devidamente ao cliente que o pedido foi
// proibido (ou um erro ocorreu).
func autenticaUsuario(w http.ResponseWriter, r *http.Request) (bool, string) {
	autenticado, email, err := cadastro.VerificaLoginRequestHTTP(r)
	if !erros.Vazio(err) {
		interno, status, _ := erros.Abre(err)
		if !interno {
			comunicacao.RespondeErro(w, status, err)
			return false, ""
		}
		comunicacao.Responde(w, status, []byte{})
		return false, ""
	} else if !autenticado {
		comunicacao.Responde(w, http.StatusForbidden, []byte{})
		return false, ""
	}
	return true, email
}

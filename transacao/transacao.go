// Package transacao abstrai a realização de transações no servidor
package transacao

import (
	"net/http"
	"strconv"

	"github.com/loteny/redcoins/comunicacao"
	"github.com/loteny/redcoins/database"
	"github.com/loteny/redcoins/erros"
	"github.com/loteny/redcoins/precobtc"
)

// Lista de possíveis erros do módulo
var (
	ErrMetodoPost  = erros.Cria(false, 405, "")
	ErrQtdInvalida = erros.Cria(false, 400, "qtd de bitcoins inválida")
)

// CompraHTTP realiza uma compra a partir de um request HTTP
func CompraHTTP(r *http.Request) error {
	// Adquire os dados da compra
	email, qtd, err := validaDadosTransacao(r)
	if err != nil {
		return err
	}
	preco, err := precobtc.Preco(qtd)
	if err != nil {
		return err
	}
	// Insere no banco de dados
	if err := database.InsereTransacao(email, true, qtd, preco); err != nil {
		return err
	}
	return nil
}

// validaDadosTransacao verifica se o pedido foi feito em POST e verifica se a
// quantidade de Bitcoins a ser comprada ou vendida é válida
func validaDadosTransacao(r *http.Request) (string, float64, error) {
	// Verifica o método do request
	if r.Method != "POST" {
		return "", 0, ErrMetodoPost
	}
	// Adquire os dados do request
	if err := comunicacao.RealizaParseForm(r); err != nil {
		return "", 0, err
	}
	email := r.PostFormValue("email")
	qtd := r.PostFormValue("qtd")
	// O e-mail já deve estar validado graças à autenticação. Mesmo que não
	// esteja, se o e-mail foi inválido de alguma forma, o ID do usuário não
	// será encontrado no banco de dados e nem o servidor nem o banco de dados
	// passaram a malfuncionar.
	fQtd, err := strconv.ParseFloat(qtd, 64)
	if err == strconv.ErrSyntax || fQtd < 0 {
		return "", 0, ErrQtdInvalida
	} else if err != nil {
		return "", 0, err
	}

	return email, fQtd, nil
}

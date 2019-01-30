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
	ErrQtdInvalida       = erros.Cria(false, 400, "qtd de bitcoins inválida")
	ErrSaldoInsuficiente = erros.Cria(false, 400, database.ErrSaldoInsuficiente.Error())
)

// CompraHTTP realiza uma compra de Bitcoins a partir de um request HTTP
func CompraHTTP(r *http.Request) error {
	return transacaoHTTP(r, true)
}

// VendaHTTP realiza uma venda de Bitcoins a partir de um request HTTP
func VendaHTTP(r *http.Request) error {
	return transacaoHTTP(r, false)
}

func transacaoHTTP(r *http.Request, compra bool) error {
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
	if err := database.InsereTransacao(email, compra, qtd, preco); err == database.ErrSaldoInsuficiente {
		return ErrSaldoInsuficiente
	} else if err != nil {
		return err
	}
	return nil
}

// validaDadosTransacao verifica se a quantidade de Bitcoins a ser comprada ou
// vendida é válida e retorna o e-mail do usuário e a quantidade de Bitcoins
// para transação.
func validaDadosTransacao(r *http.Request) (string, float64, error) {
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

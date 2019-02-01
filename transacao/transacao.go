// Package transacao abstrai a realização de transações no servidor.
// Esse package usa exclusivamente a estrutura de erros 'erros.Erros'.
package transacao

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/loteny/redcoins/comunicacao"
	"github.com/loteny/redcoins/database"
	"github.com/loteny/redcoins/erros"
	"github.com/loteny/redcoins/precobtc"
)

// Lista de possíveis erros do módulo
var (
	ErrQtdInvalida       = erros.Cria(false, 400, "qtd_invalida")
	ErrDataInvalida      = erros.Cria(false, 400, "data_invalida")
	ErrSaldoInsuficiente = erros.Cria(false, 400, "saldo_insuficiente")
)

// CompraHTTP realiza uma compra de Bitcoins a partir de um request HTTP
func CompraHTTP(r *http.Request, email string) erros.Erros {
	return transacaoHTTP(r, email, true)
}

// VendaHTTP realiza uma venda de Bitcoins a partir de um request HTTP
func VendaHTTP(r *http.Request, email string) erros.Erros {
	return transacaoHTTP(r, email, false)
}

// TransacoesDiaHTTP adquire todas as transações em um dia "YYYY-MM-DD" no campo
// "data", retornando os bytes da string JSON com as transações para o cliente
func TransacoesDiaHTTP(r *http.Request) ([]byte, erros.Erros) {
	// Adquire o e-mail do request
	if err := comunicacao.RealizaParseForm(r); err != nil {
		return nil, erros.CriaInternoPadrao(err)
	}
	data := r.FormValue("data")

	// Adquire as transações
	transacoes, err := database.AdquireTransacoesEmDia(data)
	if err != nil {
		return nil, erros.CriaInternoPadrao(err)
	}
	trBytes, err := json.Marshal(map[string][]database.Transacao{"transacoes": transacoes})
	if err != nil {
		return nil, erros.CriaInternoPadrao(err)
	}
	return trBytes, erros.CriaVazio()
}

// TransacoesUsuarioHTTP adquire todas as transações de um usuário a partir de
// seu e-mail no campo "email", retornando os bytes da string JSON com as
// transações para o cliente
func TransacoesUsuarioHTTP(r *http.Request) ([]byte, erros.Erros) {
	// Adquire o e-mail do request
	if err := comunicacao.RealizaParseForm(r); err != nil {
		return nil, erros.CriaInternoPadrao(err)
	}
	email := r.FormValue("email")

	// Adquire as transações
	transacoes, err := database.AdquireTransacoesDeUsuario(email)
	if err != nil {
		return nil, erros.CriaInternoPadrao(err)
	}
	trBytes, err := json.Marshal(map[string][]database.Transacao{"transacoes": transacoes})
	if err != nil {
		return nil, erros.CriaInternoPadrao(err)
	}
	return trBytes, erros.CriaVazio()
}

func transacaoHTTP(r *http.Request, email string, compra bool) erros.Erros {
	// Adquire os dados da compra
	qtd, data, err := validaDadosTransacao(r)
	if !erros.Vazio(err) {
		return err
	}
	preco, err2 := precobtc.Preco(qtd)
	if err2 != nil {
		return erros.CriaInternoPadrao(err2)
	}
	// Insere no banco de dados
	if err := database.InsereTransacao(email, compra, qtd, preco, data); err == database.ErrSaldoInsuficiente {
		return ErrSaldoInsuficiente
	} else if err != nil {
		return erros.CriaInternoPadrao(err)
	}
	return erros.CriaVazio()
}

// validaDadosTransacao verifica se a quantidade de Bitcoins a ser comprada ou
// vendida é válida e retorna a quantidade de Bitcoins para transação e a data
// da transação.
func validaDadosTransacao(r *http.Request) (float64, string, erros.Erros) {
	// Adquire os dados do request
	if err := comunicacao.RealizaParseForm(r); err != nil {
		return 0, "", erros.CriaInternoPadrao(err)
	}
	qtd := r.PostFormValue("qtd")
	data := r.PostFormValue("data")
	// O e-mail já deve estar validado graças à autenticação. Mesmo que não
	// esteja, se o e-mail foi inválido de alguma forma, o ID do usuário não
	// será encontrado no banco de dados e nem o servidor nem o banco de dados
	// passaram a malfuncionar.
	fQtd, err := strconv.ParseFloat(qtd, 64)
	if err != nil || fQtd < 0 {
		return 0, "", ErrQtdInvalida
	}
	// Data da transação
	if _, err := time.Parse("2006-01-02", data); err != nil {
		return 0, "", ErrDataInvalida
	}

	return fQtd, data, erros.CriaVazio()
}

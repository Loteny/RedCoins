// Package database serve para abstrair operações com o banco de dados.
// Devido à natureza de alteração de estruturas e dados do banco de dados, os
// testes devem ser executados sequencialmente para maior confiabilidade.
package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// Inicializa as configurações da package com as variáveis de ambiente
	usuarioDb = os.Getenv("REDCOINS_DB_USR")
	senhaDb = os.Getenv("REDCOINS_DB_SENHA")
	dbNome = os.Getenv("REDCOINS_DB_TESTEDBNOME")
	testDbNome = os.Getenv("REDCOINS_DB_TESTEDBNOME")
	enderecoDb = os.Getenv("REDCOINS_DB_DBADDR")
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome

	// Inicializa o banco de dados com alguns valores padrões
	testPopulaDatabase()
}

func TestInsereUsuario(t *testing.T) {
	// Usuário a ser criado
	usr := Usuario{
		Email:      "testeinsereusuario@gmail.com",
		Senha:      []byte("123456"),
		Nascimento: "1942-07-10",
		Nome:       "Ronnie James Dio",
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	if err := InsereUsuario(&usr); err != nil {
		t.Fatalf("Erro ao inserir dado no banco de dados: %v", err)
	}

	// Verifica se o usuário foi inserido corretamente
	sqlCode := `SELECT email, senha, nome, nascimento
		FROM usuario
		WHERE email=?;`
	usrResposta := Usuario{}
	if err := db.QueryRow(sqlCode, usr.Email).Scan(
		&usrResposta.Email,
		&usrResposta.Senha,
		&usrResposta.Nome,
		&usrResposta.Nascimento); err != nil {
		t.Fatalf("Erro ao adquirir a linha de usuário: %v", err)
	}
	if !(usrResposta.Email == usr.Email &&
		bytes.Equal(usrResposta.Senha, usr.Senha) &&
		usrResposta.Nome == usr.Nome &&
		usrResposta.Nascimento == usr.Nascimento) {
		t.Fatalf("Usuário inserido incorretamente.\nOriginal: %v\nAdquirido: %v", usr, usrResposta)
	}

	// Verifica se o código corretamente retorna o erro adequado ao cadastrar um usuário repetido
	if err := InsereUsuario(&usr); err != ErrUsuarioDuplicado {
		t.Fatalf("Erro inesperado ao inserir usuário duplicado: %v", err)
	}
}

func TestAdquireSenhaHashed(t *testing.T) {
	// Usuário existente
	senha, err := AdquireSenhaHashed("valido1@gmail.com")
	if err != nil {
		t.Fatalf("Erro inesperado ao adquirir senha: %v", err)
	}
	if string(senha) != "senhavalido1" {
		t.Errorf("Senha retornada incorretamente: %v", senha)
	}

	// Usuário não existente
	senha, err = AdquireSenhaHashed("naoexistente@gmail.com")
	if err != ErrUsuarioNaoExiste {
		t.Errorf("Retorno inesperado para usuário inexistente.\nSenha: %v\nErro: %v", senha, err)
	}
}

func TestInsereTransacao(t *testing.T) {
	// Compra inicial que não deve dar erros
	err := InsereTransacao("valido3@gmail.com", true, 0.00001, 0.00001, "2012-01-01")
	if err != nil {
		t.Fatalf("Erro inesperado na transação: %v", err)
	}

	// Venda que deve ocorrer corretamente
	err = InsereTransacao("valido3@gmail.com", false, 0.000005, 0.00001, "2012-01-01")
	if err != nil {
		t.Fatalf("Erro inesperado na transação: %v", err)
	}

	// Venda que deve acarretar em saldo insuficiente
	err = InsereTransacao("valido3@gmail.com", false, 0.00000501, 0.00001, "2012-01-01")
	if err != ErrSaldoInsuficiente {
		t.Fatalf("Erro inesperado na transação: %v", err)
	}
}

func TestAdquireTransacoesDeUsuario(t *testing.T) {
	// Utiliza a conta valido1@gmail.com para checar as transações
	transacoes, err := AdquireTransacoesDeUsuario("valido1@gmail.com")
	if err != nil {
		t.Errorf("Erro inesperado ao adquirir transações: %v", err)
	}

	valorEsperado := `[{valido1@gmail.com true 10 0.004 2018-01-01} ` +
		`{valido1@gmail.com false 30 0.002 2018-01-02}]`
	if valorEsperado != fmt.Sprintf("%v", transacoes) {
		t.Errorf("Lista de transações possui valor inesperado: %v", transacoes)
	}
}

func TestAdquireTransacoesEmDia(t *testing.T) {
	// Utiliza a data 2018-01-02 para checar as transações
	transacoes, err := AdquireTransacoesEmDia("2018-01-02")
	if err != nil {
		t.Errorf("Erro inesperado ao adquirir transações: %v", err)
	}

	valorEsperado := `[{valido2@gmail.com true 20 0.003 2018-01-02} ` +
		`{valido1@gmail.com false 30 0.002 2018-01-02}]`
	if valorEsperado != fmt.Sprintf("%v", transacoes) {
		t.Errorf("Lista de transações possui valor inesperado: %v", transacoes)
	}
}

// testPopulaDatabase deleta o banco de dados de testes, cria novamente e cria:
// - 3 usuários, sendo o último sem transação
// - 2 compras em dias diferentes, uma parada cada usuário
// - 2 vendas em dias diferentes, uma para cada usuário
// A venda do usuário 1 ocorre no mesmo dia que a compra do usuário 2.
func testPopulaDatabase() {
	// Recria o banco de dados
	testResetaDatabase()
	db, err := sql.Open("mysql", dsn)
	defer db.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Usuários
	sqlCode := `INSERT INTO usuario
		(email, senha, nome, nascimento)
		VALUES
			("valido1@gmail.com", "senhavalido1", "Conta Válida 1", "1994-03-07"),
			("valido2@gmail.com", "senhavalido2", "Conta Válida 2", "1994-03-08"),
			("valido3@gmail.com", "senhavalido3", "Conta Válida 3", "1994-03-08");`
	if _, err := db.Exec(sqlCode); err != nil {
		log.Fatalf("%v", err)
	}
	// Transações
	sqlCode = `INSERT INTO transacao
		(usuario_id, compra, creditos, bitcoins, dia)
		VALUES
			(1, 1, 10, 0.004, "2018-01-01"),
			(2, 1, 20, 0.003, "2018-01-02"),
			(1, 0, 30, 0.002, "2018-01-02"),
			(2, 0, 40, 0.001, "2018-01-03");`
	if _, err := db.Exec(sqlCode); err != nil {
		log.Fatalf("%v", err)
	}
}

// testResetaDatabase deleta o banco de dados e cria novamente
func testResetaDatabase() {
	// Deleta o banco de dados
	tempDSN := usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/"
	db, err := sql.Open("mysql", tempDSN)
	defer db.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}
	sqlCode := `DROP DATABASE IF EXISTS ` + testDbNome + `;`
	if _, err = db.Exec(sqlCode); err != nil {
		log.Fatalf("%v", err)
	}

	// Cria o banco de dados
	sqlCode = `CREATE DATABASE ` + testDbNome + ` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`
	if _, err := db.Exec(sqlCode); err != nil {
		log.Fatalf("%v", err)
	}

	// Cria as tabelas
	if err := criaTabelas(); err != nil {
		log.Fatalf("%v", err)
	}
}

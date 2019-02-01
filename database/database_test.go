// Package database serve para abstrair operações com o banco de dados.
// Devido à natureza de alteração de estruturas e dados do banco de dados, os
// testes devem ser executados sequencialmente para maior confiabilidade.
package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// Inicializa as configurações do módulo com o arquivo config.json
	arquivoConfig, err := os.Open("./config.json")
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo de configurações da database: %s", err)
	}
	var c config
	if err := json.NewDecoder(arquivoConfig).Decode(&c); err != nil {
		log.Fatalf("Erro ao ler configurações da database: %s", err)
	}
	usuarioDb = c.Database.UsuarioDb
	senhaDb = c.Database.SenhaDb
	dbNome = c.Database.TestDbNome
	testDbNome = c.Database.TestDbNome
	enderecoDb = c.Database.EnderecoDb
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome

	// Inicializa o banco de dados com alguns valores padrões
}

func TestInsereUsuario(t *testing.T) {
	usr := Usuario{
		Email:      "teste@gmail.com",
		Senha:      []byte("123456"),
		Nascimento: "1942-07-10",
		Nome:       "Ronnie James Dio",
	}

	// Altera o banco de dados usado pelo módulo para usar o de testes
	backupDsn := dsn
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome
	defer func() { dsn = backupDsn }()
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
	row := db.QueryRow(sqlCode, usr.Email)
	usrResposta := Usuario{}
	if err := row.Scan(
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
	// Altera o banco de dados usado pelo módulo para usar o de testes
	backupDsn := dsn
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome
	defer func() { dsn = backupDsn }()

	// Usuário existente
	senha, err := AdquireSenhaHashed("teste@gmail.com")
	if err != nil {
		t.Fatalf("Erro inesperado ao adquirir senha: %v", err)
	}
	if string(senha) != "123456" {
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
	err := InsereTransacao("teste@gmail.com", true, 0.00001, 0.00001, "2018-01-01")
	if err != nil {
		t.Fatalf("Erro inesperado na transação: %v", err)
	}

	// Venda que deve ocorrer corretamente
	err = InsereTransacao("teste@gmail.com", false, 0.000005, 0.00001, "2018-01-01")
	if err != nil {
		t.Fatalf("Erro inesperado na transação: %v", err)
	}

	// Venda que deve acarretar em saldo insuficiente
	err = InsereTransacao("teste@gmail.com", false, 0.00000501, 0.00001, "2018-01-01")
	if err != ErrSaldoInsuficiente {
		t.Fatalf("Erro inesperado na transação: %v", err)
	}
}

func TestAdquireTransacoesDeUsuario(t *testing.T) {
	transacoes, err := AdquireTransacoesDeUsuario("teste@gmail.com")
	if err != nil {
		t.Errorf("Erro inesperado ao adquirir transações: %v", err)
	}

	valorEsperado := `[{teste@gmail.com true 1350 0.001 2018-03-07} ` +
		`{teste@gmail.com false 253 0.00029 2018-03-07} ` +
		`{teste@gmail.com false 563 0.00045 2018-08-27}]`
	if valorEsperado != fmt.Sprintf("%v", transacoes) {
		t.Errorf("Lista de transações possui valor inesperado: %v", transacoes)
	}
}

func TestAdquireTransacoesEmDia(t *testing.T) {
	transacoes, err := AdquireTransacoesEmDia("2018-03-07")
	if err != nil {
		t.Errorf("Erro inesperado ao adquirir transações: %v", err)
	}

	valorEsperado := `[{teste@gmail.com true 1350 0.001 2018-03-07} ` +
		`{teste@gmail.com false 253 0.00029 2018-03-07}]`
	if valorEsperado != fmt.Sprintf("%v", transacoes) {
		t.Errorf("Lista de transações possui valor inesperado: %v", transacoes)
	}
}

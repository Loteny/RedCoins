// Package database serve para abstrair operações com o banco de dados.
// Devido à natureza de alteração de estruturas e dados do banco de dados, os
// testes devem ser executados sequencialmente para maior confiabilidade.
package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCriaTabelas(t *testing.T) {
	// Abre o banco de dados de testes
	backupDsn := testAlteraDsn()
	defer testRetornaDsn(backupDsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	// Lista todas as tabelas a serem criadas
	// Atenção para a ordem da lista: as tabelas serão deletadas nessa ordem,
	// com verificações de foreign keys. Isso pode fazer com que uma tabela não
	// possa ser deletada porque é referenciada por uma outra.
	tabelas := [2]string{"transacao", "usuario"}

	// Deleta as tabelas se existem
	for _, tabela := range tabelas {
		testDeletaTabela(t, db, tabela)
	}

	// Cria as tabelas
	if err := CriaTabelas(); err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Verifica se as tabelas existem
	for _, tabela := range tabelas {
		sqlCode := `SELECT COUNT(*)
			FROM information_schema.tables
			WHERE
				TABLE_SCHEMA = ? AND
				TABLE_NAME = ?;`
		var qtdTabelas int
		rows := db.QueryRow(sqlCode, testDbNome, tabela)
		if err := rows.Scan(&qtdTabelas); err != nil {
			if err == sql.ErrNoRows {
				t.Errorf("Tabela %v não foi criada.", tabela)
			} else {
				t.Fatalf("Erro inesperado na query: %v", err)
			}
		}
		if qtdTabelas != 1 {
			t.Errorf("Quantidade inesperada de tabelas: %v (deveria ser 1)",
				qtdTabelas)
		}
	}
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
	backupDsn := testAlteraDsn()
	defer testRetornaDsn(backupDsn)

	// Limpa a tabela de transações
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir banco de dados: %v", err)
	}
	testLimpaTabela(t, db, "transacao")

	// Compra inicial que não deve dar erros
	err = InsereTransacao("teste@gmail.com", true, 0.00001, 0.00001, "2018-01-01")
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
	backupDsn := testAlteraDsn()
	defer testRetornaDsn(backupDsn)

	// Popula as tabelas do banco de dados
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir banco de dados: %v", err)
	}
	testPopulaTabelas(t, db)

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
	backupDsn := testAlteraDsn()
	defer testRetornaDsn(backupDsn)

	// Popula as tabelas do banco de dados
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir banco de dados: %v", err)
	}
	testPopulaTabelas(t, db)

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

// testAlteraDsn faz com que o módulo use o banco de dados de testes
func testAlteraDsn() string {
	backupDsn := dsn
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome
	return backupDsn
}

// testRetornaDsn desfaz a função testAlteraDsn
func testRetornaDsn(backupDsn string) {
	dsn = backupDsn
}

// testDeletaTabela deleta a tabela se existir
func testDeletaTabela(t *testing.T, db *sql.DB, tabela string) {
	sqlCode := `DROP TABLE IF EXISTS ` + tabela + `;`
	if _, err := db.Exec(sqlCode); err != nil {
		t.Fatalf("Erro inesperado ao deletar tabela %v: %v", tabela, err)
	}
}

// testLimpaTabela realiza a função 'truncate' na tabela selecionada
func testLimpaTabela(t *testing.T, db *sql.DB, tabela string) {
	sqlCode := `SET FOREIGN_KEY_CHECKS = 0;`
	if _, err := db.Exec(sqlCode); err != nil {
		t.Fatalf("Erro inesperado ao limpar tabela %v: %v", tabela, err)
	}
	sqlCode = `TRUNCATE TABLE ` + tabela + `;`
	if _, err := db.Exec(sqlCode); err != nil {
		t.Fatalf("Erro inesperado ao limpar tabela %v: %v", tabela, err)
	}
	sqlCode = `SET FOREIGN_KEY_CHECKS = 1;`
	if _, err := db.Exec(sqlCode); err != nil {
		t.Fatalf("Erro inesperado ao limpar tabela %v: %v", tabela, err)
	}
}

// testPopulaTabelas popula todas as tabelas do banco de dados com uma variação
// de dados. Essa função limpa as tabelas antes de popular.
func testPopulaTabelas(t *testing.T, db *sql.DB) {
	testLimpaTabela(t, db, "transacao")
	testLimpaTabela(t, db, "usuario")
	testPopulaUsuario(t, db)
	testPopulaTransacao(t, db)
}

// testPopulaUsuario popula a tabela de usuários
func testPopulaUsuario(t *testing.T, db *sql.DB) {
	usr := Usuario{
		Email:      "teste@gmail.com",
		Senha:      []byte("123456"),
		Nascimento: "1942-07-10",
		Nome:       "Ronnie James Dio",
	}
	if err := InsereUsuario(&usr); err != nil {
		t.Fatalf("Erro ao popular tabela de usuários: %v", err)
	}

	usr = Usuario{
		Email:      "segundo@hotmail.com",
		Senha:      []byte("password"),
		Nascimento: "1946-09-05",
		Nome:       "Freddie Mercury",
	}
	if err := InsereUsuario(&usr); err != nil {
		t.Fatalf("Erro ao popular tabela de usuários: %v", err)
	}
}

// testPopulaTransacao popula a tabela de transações (dependente de dados
// inseridos com a função testPopulaUsuario)
func testPopulaTransacao(t *testing.T, db *sql.DB) {
	// As duas primeiras transações do primeiro usuário serão no mesmo dia (07/03/2018)

	// Primeiro usuário
	// Compra de algumas Bitcoins
	if err := InsereTransacao("teste@gmail.com", true, 0.001, 1350, "2018-03-07"); err != nil {
		t.Fatalf("Erro ao popular tabela de transações: %v", err)
	}
	// Vendas de algumas Bitcoins
	if err := InsereTransacao("teste@gmail.com", false, 0.00029, 253, "2018-03-07"); err != nil {
		t.Fatalf("Erro ao popular tabela de transações: %v", err)
	}
	if err := InsereTransacao("teste@gmail.com", false, 0.00045, 563, "2018-08-27"); err != nil {
		t.Fatalf("Erro ao popular tabela de transações: %v", err)
	}

	// Uma simples transação para um segundo usuário
	if err := InsereTransacao("segundo@hotmail.com", true, 0.023, 5826, "2019-01-02"); err != nil {
		t.Fatalf("Erro ao popular tabela de transações: %v", err)
	}
}

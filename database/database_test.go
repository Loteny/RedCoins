// Package database serve para abstrair operações com o banco de dados
package database

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// Nome do banco de dados para testes
var testDbNome = "redcoins_teste"

func TestCriaTabelas(t *testing.T) {
	// Altera o banco de dados usado pelo módulo para usar o de testes
	backupDsn := dsn
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome
	defer func() { dsn = backupDsn }()

	// Lista todas as tabelas a serem criadas
	tabelas := [1]string{"usuario"}

	// Deleta as tabelas se existem
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}
	for _, tabela := range tabelas {
		sqlCode := `DROP TABLE IF EXISTS ` + tabela + `;`
		if _, err := db.Exec(sqlCode); err != nil {
			t.Fatalf("Erro inesperado ao deletar tabelas: %v", err)
		}
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

func TestCriaTabelaUsuario(t *testing.T) {
	// Altera o banco de dados usado pelo módulo para usar o de testes
	backupDsn := dsn
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome
	defer func() { dsn = backupDsn }()

	// Deleta a tabela se existir
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}
	sqlCode := `DROP TABLE IF EXISTS usuario;`
	if _, err := db.Exec(sqlCode); err != nil {
		t.Fatalf("Erro inesperado ao deletar tabela: %v", err)
	}

	// Cria a tabela
	if err := criaTabelaUsuario(db); err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Verifica se as tabelas existem
	sqlCode = `SELECT COUNT(*)
		FROM information_schema.tables
		WHERE
			TABLE_SCHEMA = ? AND
			TABLE_NAME = ?;`
	var qtdTabelas int
	rows := db.QueryRow(sqlCode, testDbNome, "usuario")
	if err := rows.Scan(&qtdTabelas); err != nil {
		if err == sql.ErrNoRows {
			t.Errorf("Tabela não foi criada.")
		} else {
			t.Fatalf("Erro inesperado na query: %v", err)
		}
	}
	if qtdTabelas != 1 {
		t.Errorf("Quantidade inesperada de tabelas: %v (deveria ser 1)",
			qtdTabelas)
	}
}

func TestInsereUsuario(t *testing.T) {
	usr := Usuario{
		email:      "teste@gmail.com",
		senha:      "123456",
		senhaHash:  "hash_teste",
		nascimento: "1942-07-10",
		nome:       "Ronnie James Dio",
	}

	// Altera o banco de dados usado pelo módulo para usar o de testes
	backupDsn := dsn
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome
	defer func() { dsn = backupDsn }()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	// Limpa a tabela de usuários se existir
	sqlCode := `TRUNCATE usuario;`
	if _, err := db.Exec(sqlCode); err != nil {
		t.Fatalf("Erro inesperado ao limpar tabela: %v", err)
	}

	if err := InsereUsuario(&usr); err != nil {
		t.Fatalf("Erro ao inserir dado no banco de dados: %v", err)
	}

	// Verifica se o usuário foi inserido corretamente
	sqlCode = `SELECT email, senha, senha_hash, nome, nascimento
		FROM usuario
		WHERE email=?;`
	row := db.QueryRow(sqlCode, usr.email)
	usrResposta := Usuario{}
	if err := row.Scan(
		&usrResposta.email,
		&usrResposta.senha,
		&usrResposta.senhaHash,
		&usrResposta.nome,
		&usrResposta.nascimento); err != nil {
		t.Fatalf("Erro ao adquirir a linha de usuário: %v", err)
	}
	if usrResposta != usr {
		t.Fatalf("Usuário inserido incorretamente.\nOriginal: %v\nAdquirido: %v", usr, usrResposta)
	}

	// Verifica se o código corretamente retorna o erro adequado ao cadastrar um usuário repetido
	if err := InsereUsuario(&usr); err != ErrUsuarioDuplicado {
		t.Fatalf("Erro inesperado ao inserir usuário duplicado: %v", err)
	}
}

func TestAdquireSenhaEHash(t *testing.T) {
	// Altera o banco de dados usado pelo módulo para usar o de testes
	backupDsn := dsn
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome
	defer func() { dsn = backupDsn }()

	// Usuário existente
	senha, hash, err := AdquireSenhaEHash("teste@gmail.com")
	if err != nil {
		t.Fatalf("Erro inesperado ao adquirir senha/hash: %v", err)
	}
	if senha != "123456" || hash != "hash_teste" {
		t.Errorf("Dados retornados incorretamente.\nSenha: %v\nHash: %v", senha, hash)
	}

	// Usuário não existente
	senha, hash, err = AdquireSenhaEHash("naoexistente@gmail.com")
	if err != ErrUsuarioNaoExiste {
		t.Errorf("Retorno inesperado para usuário inexistente.\nSenha: %v\nHash: %v\nErro: %v", senha, hash, err)
	}
}

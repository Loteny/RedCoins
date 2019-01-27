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
	dsn = usuario + ":" + senha + "@tcp(" + endereco + ")/" + testDbNome
	defer func() { dsn = backupDsn }()

	// Lista todas as tabelas a serem criadas
	tabelas := [1]string{"usuarios"}

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
	dsn = usuario + ":" + senha + "@tcp(" + endereco + ")/" + testDbNome
	defer func() { dsn = backupDsn }()

	// Deleta a tabela se existir
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}
	sqlCode := `DROP TABLE IF EXISTS usuarios;`
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
	rows := db.QueryRow(sqlCode, testDbNome, "usuarios")
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

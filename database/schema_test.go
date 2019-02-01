package database

// Os testes desse arquivo devem ser rodados um a um, sem ficar em paralelo
// com quaisquer outros testes do banco de dados.

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCriaDatabase(t *testing.T) {
	// Primeiro, excluímos o banco de dados para ter certeza de que ele não
	// existe
	tempDSN := usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/"
	db, err := sql.Open("mysql", tempDSN)
	defer db.Close()
	if err != nil {
		t.Fatalf("%v", err)
	}
	sqlCode := `DROP DATABASE IF EXISTS ` + testDbNome + `;`
	if _, err := db.Exec(sqlCode); err != nil {
		t.Fatalf("%v", err)
	}

	// Cria o banco de dados e verifica de que ele e suas tabelas existem
	if err := CriaDatabase(); err != nil {
		t.Fatalf("%v", err)
	}
	var qtd uint
	sqlCode = `SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name=?;`
	if err := db.QueryRow(sqlCode, dbNome).Scan(&qtd); err != nil {
		t.Fatalf("%v", err)
	} else if qtd != 1 {
		t.Fatalf("Banco de dados não criado. Qtd: %v", qtd)
	}
	sqlCode = `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema=?;`
	if err := db.QueryRow(sqlCode, dbNome).Scan(&qtd); err != nil {
		t.Fatalf("%v", err)
	} else if qtd != 2 {
		t.Errorf("Quantidade inesperada de tabelas: %v", qtd)
	}

	// Testamos para verificar se o banco de dados não é recriado se já existe.
	// Inserimos um valor qualquer em uma tabela para verificar se ele persiste.
	db2, err := sql.Open("mysql", dsn)
	defer db2.Close()
	if err != nil {
		t.Fatalf("%v", err)
	}
	sqlCode = `INSERT INTO usuario (email, senha, nome, nascimento) VALUES ("teste@gmail.com", "123456", "Usuário Teste", "2018-01-01");`
	if _, err := db2.Exec(sqlCode); err != nil {
		t.Fatalf("%v", err)
	}
	if err := CriaDatabase(); err != nil {
		t.Fatalf("%v", err)
	}
	sqlCode = `SELECT COUNT(*) FROM usuario;`
	if err := db2.QueryRow(sqlCode).Scan(&qtd); err != nil {
		t.Fatalf("%v", err)
	} else if qtd != 1 {
		t.Errorf("Quantidade inesperada de usuários: %v", qtd)
	}
}

func TestDeletaDatabaseTeste(t *testing.T) {
	// Verifica se o banco de dados existe. Se não existe, cria.
	tempDSN := usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/"
	db, err := sql.Open("mysql", tempDSN)
	defer db.Close()
	if err != nil {
		t.Fatalf("%v", err)
	}
	// Verifica a quantidade de banco de dados com o nome do banco de dados do
	// servidor (deve ser 0 ou 1). Cria se não existe.
	var qtd uint
	sqlCode := `SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name=?;`
	if err := db.QueryRow(sqlCode, dbNome).Scan(&qtd); err != nil {
		t.Fatalf("%v", err)
	} else if qtd == 0 {
		sqlCode := `CREATE DATABASE ` + dbNome + ` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`
		if _, err := db.Exec(sqlCode); err != nil {
			t.Fatalf("%v", err)
		}
	}

	// Deleta o banco de dados
	if err := DeletaDatabaseTeste(); err != nil {
		t.Fatalf("%v", err)
	}
	// Verifica se foi realmente deletado
	sqlCode = `SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name=?;`
	if err := db.QueryRow(sqlCode, dbNome).Scan(&qtd); err != nil {
		t.Fatalf("%v", err)
	} else if qtd != 0 {
		t.Fatalf("Banco de dados não foi deletado (qtd: %v)", qtd)
	}

	// Tentar deletar um banco de dados inexistente não deve gerar erros
	if err := DeletaDatabaseTeste(); err != nil {
		t.Fatalf("%v", err)
	}
}

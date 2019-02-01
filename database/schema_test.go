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
	db, err := sql.Open("mysql", dsn)
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
	}
	if qtd != 1 {
		t.Fatalf("Banco de dados não criado. Qtd: %v", qtd)
	}
	sqlCode = `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema=?;`
	if err := db.QueryRow(sqlCode, dbNome).Scan(&qtd); err != nil {
		t.Fatalf("%v", err)
	}
	if qtd != 2 {
		t.Errorf("Quantidade inesperada de tabelas: %v", qtd)
	}
}

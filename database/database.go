// Package database serve para abstrair operações com o banco de dados
package database

import (
	"database/sql"

	// Driver MySQL
	_ "github.com/go-sql-driver/mysql"
)

// Configurações para o banco de dados
var (
	// Nome de usuário
	usuario = "root"
	// Senha do usuário
	//senha = "tvM@v:2gj@A')cH5"
	senha = ""
	// Nome do banco de dados
	dbNome = "redcoins"
	// Endereço do bando de dados com port
	endereco = "localhost:55555"
)

// Data Source Name: string completa para conexão com o banco de dados
var dsn = usuario + ":" + senha + "@tcp(" + endereco + ")/" + dbNome

// CriaTabelas cria as tabelas do banco de dados do servidor
func CriaTabelas() error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err := criaTabelaUsuario(db); err != nil {
		return err
	}

	return nil
}

// criaTabelaUsuario cria a tabela 'usuarios' no banco de dados que armazena
// os dados cadastrais dos usuários
func criaTabelaUsuario(db *sql.DB) error {
	sqlCode := `CREATE TABLE usuarios (
		email VARCHAR(128),
		senha VARCHAR(64) NOT NULL,
		senha_hash CHAR(32) NOT NULL,
		nome VARCHAR(255) NOT NULL,
		nascimento DATE NOT NULL,
		CONSTRAINT pk_usuarios_email PRIMARY KEY (email)
	) ENGINE=InnoDB;`
	if _, err := db.Exec(sqlCode); err != nil {
		return err
	}
	return nil
}

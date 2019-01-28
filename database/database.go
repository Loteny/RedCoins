// Package database serve para abstrair operações com o banco de dados
package database

import (
	"database/sql"
	"errors"

	// Driver MySQL
	_ "github.com/go-sql-driver/mysql"
)

// Configurações para o banco de dados
var (
	// Nome de usuário
	usuarioDb = "root"
	// Senha do usuário
	//senha = "tvM@v:2gj@A')cH5"
	senhaDb = ""
	// Nome do banco de dados
	dbNome = "redcoins"
	// Endereço do bando de dados com port
	enderecoDb = "localhost:55555"
	// Data Source Name: string completa para conexão com o banco de dados
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + dbNome
)

// Erros possíveis do módulo
var (
	ErrUsuarioDuplicado = errors.New("E-mail já cadastrado")
	ErrUsuarioNaoExiste = errors.New("Usuário não existente")
)

// Usuario é a estrutura para a tabela 'usuario'
type Usuario struct {
	email      string
	senha      string
	senhaHash  string
	nome       string
	nascimento string
}

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

// InsereUsuario cria uma nova linha na tabela 'usuario'. Retorna
// ErrUsuarioDuplicado se o usuário é repetido (mesmo e-mail)
func InsereUsuario(usr *Usuario) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Verifica se o usuário existe
	if err := verificaUsuarioDuplicado(db, usr.email); err != nil {
		return err
	}

	// Insere usuário no banco de dados
	sqlCode := `INSERT INTO usuario
		(email, senha, senha_hash, nome, nascimento)
		VALUES (?, ?, ?, ?, ?);`
	if _, err := db.Exec(
		sqlCode,
		usr.email,
		usr.senha,
		usr.senhaHash,
		usr.nome,
		usr.nascimento); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// AdquireSenhaEHash retorna a senha e o hash usado na senha de um usuário a
// partir de seu email. Se o usuário não existe, retorna ErrUsuarioNaoExiste.
func AdquireSenhaEHash(email string) (string, string, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", "", err
	}

	var senha, hash string
	// Adquire os dados do banco de dados
	sqlCode := `SELECT senha, senha_hash
		FROM usuario
		WHERE email=?;`
	row := db.QueryRow(sqlCode, email)
	err = row.Scan(&senha, &hash)
	// Se o usuário não existe, retorna ErrUsuarioNaoExiste
	if err == sql.ErrNoRows {
		return "", "", ErrUsuarioNaoExiste
	} else if err != nil {
		return "", "", err
	}

	return senha, hash, nil
}

// criaTabelaUsuario cria a tabela 'usuario' no banco de dados que armazena
// os dados cadastrais dos usuários
func criaTabelaUsuario(db *sql.DB) error {
	sqlCode := `CREATE TABLE usuario (
		id INT(11) UNSIGNED AUTO_INCREMENT,
		email VARCHAR(128) UNIQUE NOT NULL,
		senha VARCHAR(64) NOT NULL,
		senha_hash CHAR(32) NOT NULL,
		nome VARCHAR(255) NOT NULL,
		nascimento DATE NOT NULL,
		CONSTRAINT pk_usuario_id PRIMARY KEY (id)
	) ENGINE=InnoDB;`
	if _, err := db.Exec(sqlCode); err != nil {
		return err
	}
	return nil
}

// verificaUsuarioDuplicado verifica se existe um usuário na tabela 'usuario'
// com o e-mail passado. Se existe, retorna ErrUsuarioDuplicado. Se não existe,
// retorna nil.
func verificaUsuarioDuplicado(db *sql.DB, email string) error {
	sqlCode := `SELECT COUNT(*) FROM usuario WHERE email=?;`
	row := db.QueryRow(sqlCode, email)
	var rowQtd int
	if err := row.Scan(&rowQtd); err != nil {
		return err
	}
	if rowQtd > 0 {
		return ErrUsuarioDuplicado
	}
	return nil
}

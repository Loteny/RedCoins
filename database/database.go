// Package database serve para abstrair operações com o banco de dados
package database

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"strconv"

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
	// Nome do banco de dados de testes
	testDbNome = "redcoins_teste"
	// Endereço do bando de dados com port
	enderecoDb = "localhost:55555"
	// Data Source Name: string completa para conexão com o banco de dados
	dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + dbNome
)

// Erros possíveis do módulo
var (
	ErrUsuarioDuplicado  = errors.New("email_ja_cadastrado")
	ErrUsuarioNaoExiste  = errors.New("usuario_nao_existente")
	ErrSaldoInsuficiente = errors.New("saldo_insuficiente")
)

// Usuario é a estrutura para a tabela 'usuario'.
// O campo 'senha' deve conter até 60 bytes
type Usuario struct {
	Email      string
	Senha      []byte
	Nome       string
	Nascimento string
}

// Transacao é a estrutura com dados de uma transação
type Transacao struct {
	Usuario  string  `json:"usuario"`
	Compra   bool    `json:"compra"`
	Creditos float64 `json:"creditos"`
	Bitcoins float64 `json:"bitcoins"`
	Dia      string  `json:"dia"`
}

// CriaDatabase verifica se o banco de dados do servidor está criado. Se não
// estiver, cria. Se estiver, verifica se o banco de dados possui alguma tabela.
// Se possui, não faz nada e retorna nenhum erro. Se não possui, cria as tabelas
// necessárias para operação do servidor.
func CriaDatabase() error {
	// Entra no MySQL sem banco de dados e verifica se o banco de dados do
	// servidor existe
	tempDSN := usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/"
	db, err := sql.Open("mysql", tempDSN)
	defer db.Close()
	if err != nil {
		return err
	}
	// Verifica a quantidade de banco de dados com o nome do banco de dados do
	// servidor (deve ser 0 ou 1)
	var qtd uint
	sqlCode := `SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name=?;`
	if err := db.QueryRow(sqlCode, dbNome).Scan(&qtd); err != nil {
		return err
	}
	// Cria o banco de dados se não existe
	if qtd == 0 {
		sqlCode := `CREATE DATABASE ` + dbNome + ` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`
		if _, err := db.Exec(sqlCode); err != nil {
			return err
		}
	}
	return criaTabelas()
}

// InsereUsuario cria uma nova linha na tabela 'usuario'. Retorna
// ErrUsuarioDuplicado se o usuário é repetido (mesmo e-mail)
func InsereUsuario(usr *Usuario) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Verifica se o usuário existe
	if err := verificaUsuarioDuplicado(db, usr.Email); err != nil {
		return err
	}

	// Insere usuário no banco de dados
	sqlCode := `INSERT INTO usuario
		(email, senha, nome, nascimento)
		VALUES (?, ?, ?, ?);`
	if _, err := db.Exec(
		sqlCode,
		usr.Email,
		usr.Senha,
		usr.Nome,
		usr.Nascimento); err != nil {
		return err
	}

	return nil
}

// AdquireSenhaHashed retorna o campo 'senha' do usuário (salvada hashed) a
// partir de seu email. Se o usuário não existe, retorna ErrUsuarioNaoExiste.
func AdquireSenhaHashed(email string) ([]byte, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return []byte{}, err
	}

	senha := make([]byte, 60)
	// Adquire os dados do banco de dados
	sqlCode := `SELECT senha FROM usuario WHERE email=?;`
	row := db.QueryRow(sqlCode, email)
	err = row.Scan(&senha)
	// Se o usuário não existe, retorna ErrUsuarioNaoExiste
	if err == sql.ErrNoRows {
		return []byte{}, ErrUsuarioNaoExiste
	} else if err != nil {
		return []byte{}, err
	}

	return senha, nil
}

// InsereTransacao cria uma nova transação no banco de dados a partir do e-mail
// de um usuário, do tipo da transação (compra ou venda), a quantidade de
// BitCoins a ser comprada ou vendida, o valor pago ou recebido em reais pela
// transação e a data da transação. A quantidade de BitCoins e o valor em reais
// devem ser números inteiros: os 10 primeiros dígitos do valor em reais e os 8
// primeiro dígitos. da quantidade de BitCoins formam as partes decimais de seus
// valores reais A data deve estar no formato "YYYY-MM-DD".
func InsereTransacao(email string, compra bool, bitcoins float64, preco float64, data string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Adquire o ID do usuário para buscas mais simples
	usrID, err := adquireUsuarioIDDeEmail(tx, email)
	if err != nil {
		return err
	}
	// Trava as linhas associadas às transações do usuário para verificação
	// de saldo e verifica se possui saldo suficiente para a transação
	saldoBitcoins, err := adquireSaldosUsuario(tx, usrID)
	if err != nil {
		return err
	} else if !compra && saldoBitcoins < bitcoins {
		return ErrSaldoInsuficiente
	}

	// Insere a transação no banco de dados
	if err := insereLinhaTransacao(tx, usrID, compra, preco, bitcoins, data); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// AdquireTransacoesDeUsuario adquire todas as transações feitas por um usuário
// identificado pelo seu e-mail na forma []Transacao
func AdquireTransacoesDeUsuario(email string) ([]Transacao, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	sqlCode := `SELECT
		t.compra, t.creditos, t.bitcoins, t.dia
		FROM usuario AS u
		INNER JOIN transacao AS t ON t.usuario_id = u.id
		WHERE u.email=?;`
	rows, err := db.Query(sqlCode, email)
	if err != nil {
		return nil, err
	}

	// Armazena os resultados na variável de retorno 'transacoes'
	defer rows.Close()
	transacoes := make([]Transacao, 0)
	for rows.Next() {
		tr := Transacao{Usuario: email}
		compra := make([]uint8, 1)
		if err := rows.Scan(&compra, &tr.Creditos, &tr.Bitcoins, &tr.Dia); err != nil {
			return nil, err
		}
		tr.Compra = compra[0] == 1
		transacoes = append(transacoes, tr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transacoes, nil
}

// AdquireTransacoesEmDia adquire todas as transações feitas em um determinado
// dia no formato "YYYY-MM-DD" e retorna uma lista de Transacao
func AdquireTransacoesEmDia(dia string) ([]Transacao, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	sqlCode := `SELECT
		u.email, t.compra, t.creditos, t.bitcoins, t.dia
		FROM transacao AS t
		INNER JOIN usuario AS u ON u.id = t.usuario_id
		WHERE t.dia=?;`
	rows, err := db.Query(sqlCode, dia)
	if err != nil {
		return nil, err
	}

	// Armazena os resultados na variável de retorno 'transacoes'
	defer rows.Close()
	transacoes := make([]Transacao, 0)
	for rows.Next() {
		tr := Transacao{}
		compra := make([]uint8, 1)
		if err := rows.Scan(&tr.Usuario, &compra, &tr.Creditos, &tr.Bitcoins, &tr.Dia); err != nil {
			return nil, err
		}
		tr.Compra = compra[0] == 1
		transacoes = append(transacoes, tr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transacoes, nil
}

// verificaSemTabelas verifica se o banco de dados está sem nenhum tabela.
// Retorna 'true' caso esteja ou 'false' se existe alguma tabela no banco
// de dados.
func verificaSemTabelas(db *sql.DB) (bool, error) {
	var qtd uint
	sqlCode := `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema=?;`
	if err := db.QueryRow(sqlCode, dbNome).Scan(&qtd); err != nil {
		return false, err
	}
	return qtd == 0, nil
}

// criaTabelas cria as tabelas do banco de dados do servidor se o banco de dados
// não tiver nenhuma tabela presente. Se tiver, retorna sem erros (assume-se que
// as tabelas foram criadas corretamente). Se for necessário atualizar o banco
// de dados, é necessário fazê-lo manualmente.
func criaTabelas() error {
	db, err := sql.Open("mysql", dsn)
	defer db.Close()
	if err != nil {
		return err
	}

	// Verifica se o banco de dados está vazio
	if vazio, err := verificaSemTabelas(db); err != nil {
		return err
	} else if !vazio {
		return nil
	}

	// Todo o banco de dados deve ser gerado em uma única transação
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Criação das tabelas individualmente
	if err := criaTabelaUsuario(tx); err != nil {
		return err
	}
	if err := criaTabelaTransacao(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// criaTabelaUsuario cria a tabela 'usuario' no banco de dados que armazena
// os dados cadastrais dos usuários
func criaTabelaUsuario(tx *sql.Tx) error {
	sqlCode := `CREATE TABLE usuario (
		id INT(11) UNSIGNED AUTO_INCREMENT,
		email VARCHAR(128) UNIQUE NOT NULL,
		senha CHAR(60) NOT NULL,
		nome VARCHAR(255) NOT NULL,
		nascimento DATE NOT NULL,
		CONSTRAINT pk_usuario_id PRIMARY KEY (id)
	) ENGINE=InnoDB;`
	if _, err := tx.Exec(sqlCode); err != nil {
		return err
	}
	// Adiciona uma index no e-mail do usuário para otimizar pesquisas
	sqlCode = `ALTER TABLE usuario
		ADD INDEX idx_usuario_email (email);`
	if _, err := tx.Exec(sqlCode); err != nil {
		return err
	}
	return nil
}

// criaTabelaTransacao cria a tabela 'transacao' no banco de dados que armazena
// os dados de transações efetuadas pelos usuários.
// O valor 'creditos' indica qual foi o valor em reais adquirido ou concedido
// pelo usuário na transação, 'bitcoins' indica o mesmo para seu crédito de
// BitCoins, e 'compra' indica se a transação foi uma compra ou venda de
// BitCoins (0 = venda; 1 = compra). 'dia' indica quando a transação foi
// realizada (YYYY-MM-DD).
func criaTabelaTransacao(tx *sql.Tx) error {
	sqlCode := `CREATE TABLE transacao (
		id INT(11) UNSIGNED AUTO_INCREMENT,
		usuario_id INT(11) UNSIGNED NOT NULL,
		compra BIT(1) NOT NULL,
		creditos DECIMAL(18,9) NOT NULL,
		bitcoins DECIMAL(15,8) NOT NULL,
		dia DATE NOT NULL,
		CONSTRAINT pk_transacao_id PRIMARY KEY (id),
		CONSTRAINT fk_transacao_usuario_id
			FOREIGN KEY (usuario_id)
			REFERENCES usuario(id)
	) ENGINE=InnoDB;`
	if _, err := tx.Exec(sqlCode); err != nil {
		return err
	}
	// Adiciona uma index no ID do usuário e uma index na coluna de data para
	// otimizar pesquisas
	sqlCode = `ALTER TABLE transacao
		ADD INDEX idx_transacao_usuario_id (usuario_id);`
	if _, err := tx.Exec(sqlCode); err != nil {
		return err
	}
	sqlCode = `ALTER TABLE transacao
		ADD INDEX idx_transacao_dia (dia);`
	if _, err := tx.Exec(sqlCode); err != nil {
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

// adquireUsuarioIDDeEmail adquire o ID da tabela 'usuario' a partir de seu
// e-mail. Pode retornar ErrUsuarioNaoExiste.
func adquireUsuarioIDDeEmail(tx *sql.Tx, email string) (uint, error) {
	sqlCode := `SELECT id FROM usuario WHERE email=?;`
	var id sql.NullInt64
	err := tx.QueryRow(sqlCode, email).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, ErrUsuarioNaoExiste
	} else if err != nil {
		return 0, err
	}
	return uint(id.Int64), nil
}

// adquireSaldosUsuario adquire o saldo de BitCoins que um usuário possui pelo
// seu ID. Essa função também trava todas as linhas do usuário na tabela
// 'transacao'.
func adquireSaldosUsuario(tx *sql.Tx, usrID uint) (float64, error) {
	sqlCode := `SELECT
		IFNULL(SUM(IF(t.compra=1, t.bitcoins, -1 * t.bitcoins)), "0") AS bitcoins
		FROM transacao AS t
		WHERE t.usuario_id=?
		FOR UPDATE;`
	var bitcoins sql.NullString
	err := tx.QueryRow(sqlCode, usrID).Scan(&bitcoins)
	// Caso não haja resultados, o usuário não fez nenhuma transação e seu
	// crédito de Bitcoins é zero
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	fBitcoins, err := strconv.ParseFloat(bitcoins.String, 64)
	if err != nil {
		return 0, err
	}
	return fBitcoins, nil
}

// insereLinhaTransacao insere diretamente uma nova linha de transação no banco
// de dados.
func insereLinhaTransacao(tx *sql.Tx, usuario uint, compra bool, preco float64, bitcoins float64, data string) error {
	sqlCode := `INSERT INTO
	transacao (usuario_id, compra, creditos, bitcoins, dia)
	VALUES (?, ?, ?, ?, ?);`
	var intCompra uint8
	if compra {
		intCompra = 1
	} else {
		intCompra = 0
	}
	if _, err := tx.Exec(sqlCode, usuario, intCompra, fmt.Sprintf("%18.9f", preco), fmt.Sprintf("%15.8f", bitcoins), data); err != nil {
		return err
	}
	return nil
}

// init altera o DSN para usar o banco de dados de teste se a execução estiver
// em modo de teste
func init() {
	if flag.Lookup("test.v") != nil {
		dbNome = testDbNome
		dsn = usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/" + testDbNome
	}
}

package database

// Esse arquivo define funções que configuram o banco de dados como um todo
// (criação e remoção de tabelas, por exemplo)

import "database/sql"

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

// DeletaDatabaseTeste deleta totalmente o banco de dados de teste do servidor.
// Essa função só existe especificamente para o banco de dados de teste para
// diminuir as chances de ocorrer um erro e o banco de dados oficial ser
// deletado.
func DeletaDatabaseTeste() error {
	tempDSN := usuarioDb + ":" + senhaDb + "@tcp(" + enderecoDb + ")/"
	db, err := sql.Open("mysql", tempDSN)
	defer db.Close()
	if err != nil {
		return err
	}
	sqlCode := `DROP DATABASE IF EXISTS ` + testDbNome + `;`
	_, err = db.Exec(sqlCode)
	return err
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

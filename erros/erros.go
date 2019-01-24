// Package erros disponibiliza uma nova estrutura de erros que implementa a
// interface "error". A vantagem dessa estrutura é a presença de uma flag que
// indica se o erro é interno do servidor e deve ser registrado em logs de erros
// para análise e debug, ou se o erro é externo por parte do cliente, caso no
// qual o erro não deve ser registrado em logs, e um código apropriado deve ser
// enviado para o cliente tratar e apresentar o erro apropriadamente ao usuário
package erros

import (
	"errors"
	"log"
)

// Erros é a estrutura que deve ser usada como objeto que implementa a interface
// 'erro' nos retornos de funções
type Erros struct {
	// interno se o erro é do servidor (true) ou se é problema de operação
	// (false)
	interno bool
	// msg é a mensagem que erros têm por padrão
	// Se o erro for interno, recomenda-se que se utilize essa variável do mesmo
	// jeito que se utilizaria a mensagem de erros com as funções padrões do
	// Golang. Se o erro for externo, recomenda-se que se utilize um código ou
	// identificador pequeno para ser passado para o cliente, que identificará o
	// código e tratará o erro adequadamente, construindo por si mesmo a
	// mensagem de erro apropriada
	msg string
	// statusCode armazena o status code segundo o protocolo HTTP que a
	// resposta HTTP deve ter graças à presença do erro gerado
	statusCode int
}

// Cria gera uma nova struct Erros
func Cria(interno bool, statusCode int, msg string) Erros {
	return Erros{interno: interno, statusCode: statusCode, msg: msg}
}

// Error é a função para a estrutura 'Erros' implementar a interface 'error'
func (e Erros) Error() string {
	return e.msg
}

// Abre é a função a ser chamada na hora de tratar o erro. Essa função registra
// o erro no log de erros (no caso, utilizando o pacote 'log', então a mensagem
// vai para o stderr) e retorna, além de seus dados extras (variáveis 'interno'
// e 'statusCode'), um objeto 'error' gerado pelo package padrão do Golang com
// a mensagem de erro idêntica à do objeto.
func (e Erros) Abre() (bool, int, error) {
	if e.interno {
		log.Print(e)
	}
	return e.interno, e.statusCode, errors.New(e.msg)
}

// Padrao retorna o erro gerado pela função 'errors.New' com a mesma mensagem
// que a sua própria
func (e Erros) Padrao() error {
	return errors.New(e.msg)
}

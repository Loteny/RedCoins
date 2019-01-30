// Package erros disponibiliza uma nova struct de erros que implementa a
// interface 'error' e funções para tratar tanto essa nova estrutura quanto um
// objeto que implementa a interface 'error' comum. Assim, pode-se diferenciar
// entre erros externos e internos ao servidor, armazenar múltiplos erros em um
// só objeto de erros, etc. Funções "estáticas" desse módulo podem agir tanto em
// objetos da struct própria desse módulo quanto em quaisquer outros objetos que
// implementem a interface 'error'.
package erros

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Erros é a estrutura principal desse módulo. Pode armazenar dados extras do
// erro como o HTTP Status Code a ser retornado, se o erro é interno ou não
// (para fins de logging) e uma lista de erros ao invés de uma única mensagem de
// erro.
type Erros struct {
	// interno se o erro é do servidor (true) ou se é problema de operação
	// (false). Em geral, essa flag é útil para fins de logging: erros internos
	// são registrados no servidor, enquanto erros de operação não são.
	// Em caso de múltiplos erros, erros internos têm prioridade. Ou seja, se
	// um erro interno ocorrer e for unido com um erro de operação, o resultado
	// será o erro interno.
	interno bool
	// msg é a mensagem que erros têm por padrão
	// Se o erro for interno, recomenda-se que se utilize essa variável do mesmo
	// jeito que se utilizaria a mensagem de erros com as funções padrões do
	// Golang. Se o erro for externo, recomenda-se que se utilize um código ou
	// identificador pequeno para ser passado para o cliente, que identificará o
	// código e tratará o erro adequadamente, construindo por si mesmo a
	// mensagem de erro apropriada.
	// É tratada como uma string de um único elemento em caso de erros internos
	// e uma lista de strings em caso de erros de operação.
	msg []string
	// statusCode armazena o status code segundo o protocolo HTTP que a
	// resposta HTTP deve ter graças à presença do erro gerado
	statusCode int
}

// Cria gera uma nova estrutura 'Erros'
func Cria(interno bool, statusCode int, msg string) Erros {
	msgs := make([]string, 1)
	msgs[0] = msg
	return Erros{interno: interno, statusCode: statusCode, msg: msgs}
}

// CriaInternoPadrao cria uma estrutura 'Erros' com interno = true, statusCode =
// 500 e a mesma mensagem de erro que o erro passado
func CriaInternoPadrao(err error) Erros {
	msgs := make([]string, 1)
	msgs[0] = err.Error()
	return Erros{interno: true, statusCode: http.StatusInternalServerError, msg: msgs}
}

// Error é a função para a estrutura 'Erros' implementar a interface 'error'.
// Se o erro for interno, retorna uma string com o único erro da array. Caso
// contrário, retorna uma string com uma lista JSON de erros
func (e Erros) Error() string {
	if e.interno {
		if s := len(e.msg); s < 1 {
			return ""
		}
		return e.msg[0]
	}
	msg, err := json.Marshal(e.msg)
	if err != nil {
		return "erro interno não identificado"
	}
	return string(msg)
}

// Abre é a função a ser chamada na hora de tratar o erro. Essa função registra
// o erro no log de erros (no caso, utilizando o pacote 'log', então a mensagem
// vai para o stderr) e retorna, além de seus dados extras (variáveis 'interno'
// e 'statusCode'), um objeto 'error' gerado pelo package padrão do Golang com
// a mensagem de erro idêntica à do objeto.
// Se o erro não for uma struct 'Erros', retorna os valores padrões: erro
// interno, status code 500 e a mensagem do próprio erro.
func Abre(e error) (bool, int, error) {
	objErros, sucesso := e.(Erros)
	if !sucesso {
		log.Print(e)
		return true, http.StatusInternalServerError, e
	}
	if objErros.interno {
		log.Print(e)
	}
	return objErros.interno, objErros.statusCode, errors.New(objErros.Error())
}

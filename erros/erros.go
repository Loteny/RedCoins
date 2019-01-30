// Package erros disponibiliza uma nova struct de erros que implementa a
// interface 'error' e funções para tratar tanto essa nova estrutura quanto um
// objeto que implementa a interface 'error' comum. Assim, pode-se diferenciar
// entre erros externos e internos ao servidor, armazenar múltiplos erros em um
// só objeto de erros, etc.
// Todas as funções resultam em um objeto do tipo 'Erros'. Esse objeto possui
// uma slice de strings para armazenar múltiplas mensagens de erros, portanto,
// não é possível realizar comparações com os objetos em si.
package erros

import (
	"encoding/json"
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

// CriaVazio gera uma nova estrutura 'Erros' sem uma mensagens de erros. Útil
// para preparar para uma possível lista de erros. O erro criado não é interno
// e possui statusCode 0.
func CriaVazio() Erros {
	return Erros{interno: false, statusCode: 0, msg: []string{}}
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
func Abre(e error) (bool, int, Erros) {
	objErros, sucesso := e.(Erros)
	if !sucesso {
		log.Print(e)
		return true, http.StatusInternalServerError, CriaInternoPadrao(e)
	}
	if objErros.interno {
		log.Print(e)
	}
	return objErros.interno, objErros.statusCode, objErros
}

// Adiciona insere um erro na lista de erros da struct. Se o erro não for do
// tipo Erros, não adiciona e simplesmente retorna o erro original.
func Adiciona(e error, msg string) Erros {
	if e == nil {
		return CriaVazio()
	}
	objErros, sucesso := e.(Erros)
	if !sucesso {
		return CriaInternoPadrao(e)
	}
	objErros.msg = append(objErros.msg, msg)
	return objErros
}

// JuntaErros une dois erros. Se um dos erros for interno, este prevalecerá
// completamente e os erros não serão juntados. Se os dois forem internos, o
// primeiro na ordem da lista de argumentos prevalece. Se os dois não forem
// internos, o statusCode do primeiro prevalece e as mensagens são unidas.
func JuntaErros(e1, e2 error) Erros {
	// Retorna o único erro não-nil ou nil se um dos dois ou os dois forem nil
	if e1 == nil {
		return CriaInternoPadrao(e2)
	} else if e2 == nil {
		return CriaInternoPadrao(e1)
	}
	// Verifica se os objetos são do tipo 'Erros'
	objErros1, sucesso1 := e1.(Erros)
	if !sucesso1 {
		return CriaInternoPadrao(e1)
	}
	objErros2, sucesso2 := e2.(Erros)
	if !sucesso2 {
		return CriaInternoPadrao(e2)
	}
	// Verifica se um dos dois é interno
	if objErros1.interno {
		return objErros1
	} else if objErros2.interno {
		return objErros2
	}
	// Verifica se um dos dois está vazio
	if Vazio(objErros1) {
		return objErros2
	} else if Vazio(objErros2) {
		return objErros1
	}
	// Une os dois erros
	objErros1.msg = append(objErros1.msg, objErros2.msg...)
	return objErros1
}

// Vazio verifica se o erro existe de fato ou se é apenas uma estrutura vazia
// e nenhum erro ocorreu (identificado por statusCode = 0 e lista de mensagens
// vazia).
func Vazio(e error) bool {
	if e == nil {
		return true
	}
	objErros, sucesso := e.(Erros)
	if !sucesso {
		return false
	}
	return objErros.statusCode == 0 && len(objErros.msg) == 0
}

// Lista retorna uma lista de strings que são as mensagens de erros do objeto
func Lista(e error) []string {
	if e == nil {
		return []string{}
	}
	objErros, sucesso := e.(Erros)
	if !sucesso {
		return []string{e.Error()}
	}
	return objErros.msg
}

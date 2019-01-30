package erros

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"
)

func TestCria(t *testing.T) {
	msg := make([]string, 1)
	msg[0] = "mensagem de teste de erro"
	original := Erros{interno: true, statusCode: 200, msg: msg}
	gerado := Cria(true, 200, msg[0])

	if original.interno != gerado.interno ||
		original.statusCode != gerado.statusCode ||
		original.msg[0] != gerado.msg[0] {
		t.Errorf("Estruturas diferentes.\nGerado: %#v\nOriginal: %#v",
			gerado, original)
	}
}

func TestCriaInternoPadrao(t *testing.T) {
	err := errors.New("mensagem de teste de erro")
	gerado := CriaInternoPadrao(err)

	if gerado.Error() != err.Error() {
		t.Errorf("Mensagens de erros diferentes.\nGerado: %v\nOriginal: %v",
			gerado, err)
	}
	if !gerado.interno {
		t.Error("Função criou erro externo ao invés de interno.")
	}
	if gerado.statusCode != 500 {
		t.Errorf("Função gerou status code diferente de 500. Gerado: %v",
			gerado.statusCode)
	}
}

func TestError(t *testing.T) {
	msg := "mensagem de teste de erro"
	e := Cria(true, 500, msg)
	msgRecebida := e.Error()

	if msgRecebida != msg {
		t.Errorf("Mensagem esperada: %v\nObtida: %v", msg, msgRecebida)
	}
}

func TestAbreInterno(t *testing.T) {
	buf := testAbre(t, true)
	if buf.String() == "" {
		t.Errorf("Logging de erro incorreto. Log adquirido: %v", buf.String())
	}
}

func TestAbreExterno(t *testing.T) {
	buf := testAbre(t, false)
	if buf.String() != "" {
		t.Errorf("Logging não deveria ocorrer. Log adquirido: %v", buf.String())
	}
}

// testAbre é a função base para os testes da função Abre. Outras funções de
// teste da função Abre podem derivar dessa função
func testAbre(t *testing.T, interno bool) bytes.Buffer {
	// Código para lermos o log de erros gerado pela função
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() { log.SetOutput(os.Stderr) }()

	erroNormal := errors.New("mensagem de teste de erro")
	e := Cria(interno, 200, erroNormal.Error())
	internoRecebido, statusCode, erroRecebido := Abre(e)

	if interno {
		if erroRecebido.Error() != erroNormal.Error() {
			t.Errorf("Mensagem esperada: %v\nObtida: %v", erroNormal, erroRecebido)
		}
	} else {
		if erroRecebido.Error() != "[\""+erroNormal.Error()+"\"]" {
			t.Errorf("Mensagem esperada: %v\nObtida: %v", erroNormal, erroRecebido)
		}
	}
	if internoRecebido != interno {
		t.Errorf("Valor de 'interno' incorreto (%v, deveria ser %v)",
			internoRecebido, interno)
	}
	if statusCode != 200 {
		t.Errorf("Valor de 'statusCode' incorreto (deveria ser %v, foi %v)",
			200, statusCode)
	}

	return buf
}

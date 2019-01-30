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
	g := Cria(true, 200, msg[0])
	gerado := g.(Erros)

	if original.interno != gerado.interno ||
		original.statusCode != gerado.statusCode ||
		original.msg[0] != gerado.msg[0] {
		t.Errorf("Estruturas diferentes.\nGerado: %#v\nOriginal: %#v",
			gerado, original)
	}
}

func TestCriaVazio(t *testing.T) {
	err := CriaVazio()
	e := err.(Erros)
	if e.interno != false ||
		len(e.msg) != 0 ||
		e.statusCode != 0 {
		t.Errorf("Erro vazio criado incorretamente: %v", e)
	}
}

func TestCriaInternoPadrao(t *testing.T) {
	err := errors.New("mensagem de teste de erro")
	g := CriaInternoPadrao(err)
	gerado := g.(Erros)

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

func TestAbre(t *testing.T) {
	// Erro interno (com logging)
	buf := testAbre(t, true)
	if buf.String() == "" {
		t.Errorf("Logging de erro incorreto. Log adquirido: %v", buf.String())
	}
	// Erro externo (sem logging)
	buf = testAbre(t, false)
	if buf.String() != "" {
		t.Errorf("Logging não deveria ocorrer. Log adquirido: %v", buf.String())
	}
	// Erro que não é struct 'Erros' (com logging)
	e := errors.New("erro teste")
	interno, statusCode, err := Abre(e)
	if e != err ||
		interno != true ||
		statusCode != 500 {
		t.Errorf("Valores gerados inválidos: %v / %v / %v", interno, statusCode, err)
	}
}

func TestAdiciona(t *testing.T) {
	err := Cria(true, 500, "erro 1")
	err = Adiciona(err, "erro 2")
	e := err.(Erros)
	if e.msg[0] != "erro 1" || e.msg[1] != "erro 2" {
		t.Errorf("Mensagens de erros inesperadas.\n1: %v\n2: %v", e.msg[0], e.msg[1])
	}
}

func TestJuntaErros(t *testing.T) {
	// Testa união de dois erros não-internos
	e1 := Cria(false, 500, "e1").(Erros)
	e2 := Cria(false, 500, "e2").(Erros)
	eResultado := JuntaErros(e1, e2)
	if r := eResultado.Error(); r != `["e1","e2"]` {
		t.Errorf("União de erros teve resultado incorreto: %v", r)
	}
	// Se um dos erros for interno, o resultado deve ser ele mesmo
	e3 := Cria(true, 400, "e3").(Erros)
	eResultado = JuntaErros(e1, e3)
	if r := eResultado.Error(); r != `e3` {
		t.Errorf("União de erros teve resultado incorreto: %v", r)
	}
}

func TestTransforma(t *testing.T) {
	// Erro vazio
	err := CriaVazio()
	e := err.(Erros)
	if err := e.Transforma(); err != nil {
		t.Errorf("Erro inesperado ao transformar erro: %v", err)
	}
	// Erro com item
	err = Adiciona(e, "erro 1")
	e = err.(Erros)
	if err := e.Transforma(); err == nil {
		t.Errorf("Erro 'nil' quando não deveria ser: %v", e)
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

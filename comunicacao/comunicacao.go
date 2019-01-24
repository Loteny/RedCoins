// Package comunicacao abstrai a construção das mensagens HTTP e deixa no
// formato padronizado e correto para toda a API RedCoins
package comunicacao

import (
	"encoding/json"
	"log"
	"net/http"
)

// Responde envia uma resposta HTTP com status code 's' no formato JSON com o
// conteúdo 'r'
func Responde(w http.ResponseWriter, s int, r []byte) error {
	w.WriteHeader(s)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(r)
	if err != nil {
		log.Printf("comunicacao: Responde: %s", err)
	}
	return err
}

// RespondeSucesso considera o status code da resposta HTTP como 200 (OK) e
// invoca a função Responde
func RespondeSucesso(w http.ResponseWriter, r []byte) error {
	return Responde(w, http.StatusOK, r)
}

// RespondeErro envia a mensagem de erro passada para função em JSON
func RespondeErro(w http.ResponseWriter, statusCode int, e error) error {
	m := make(map[string]string)
	m["erro"] = e.Error()
	msg, _ := json.Marshal(m)
	return Responde(w, statusCode, msg)
}

// RealizaParseForm abstrai a realização da operação http.ParseForm
func RealizaParseForm(r *http.Request) (err error) {
	if err = r.ParseForm(); err != nil {
		log.Printf("comunicacao: RealizaParseForm: %s", err)
	}
	return
}

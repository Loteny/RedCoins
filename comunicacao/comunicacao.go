// Package comunicacao abstrai a construção das mensagens HTTP e deixa no
// formato padronizado e correto para toda a API RedCoins
package comunicacao

import (
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

// RealizaParseForm abstrai a realização da operação http.ParseForm
func RealizaParseForm(r *http.Request) (err error) {
	if err = r.ParseForm(); err != nil {
		log.Printf("comunicacao: RealizaParseForm: %s", err)
	}
	return
}

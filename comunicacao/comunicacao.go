// Package comunicacao abstrai a construção das mensagens HTTP e deixa no
// formato padronizado e correto para toda a API RedCoins
package comunicacao

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/loteny/redcoins/erros"
)

// Responde envia uma resposta HTTP com status code 's' no formato JSON com o
// conteúdo 'r'
func Responde(w http.ResponseWriter, s int, r []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	_, err := w.Write(r)
	if err != nil {
		log.Printf("comunicacao: Responde: %s", err)
	}
	return err
}

// RespondeErro envia a mensagem de erro passada para função em JSON
func RespondeErro(w http.ResponseWriter, statusCode int, e erros.Erros) error {
	m := make(map[string][]string)
	m["erros"] = erros.Lista(e)
	msg, err := json.Marshal(m)
	if err != nil {
		log.Printf("comunicacao: RespondeErro: %s", err)
	}
	return Responde(w, statusCode, msg)
}

// RealizaParseForm abstrai a realização da operação http.ParseForm
func RealizaParseForm(r *http.Request) (err error) {
	if err = r.ParseForm(); err != nil {
		log.Printf("comunicacao: RealizaParseForm: %s", err)
	}
	return
}

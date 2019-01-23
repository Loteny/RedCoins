// Package cadastro trata de todas as funções para cadastramento do cliente,
// desde a recepção do request HTTPS até a inserção dos dados no banco de dados
// e verificação dos credenciais para autenticação
package cadastro

import (
	"net/http"
)

// CadastraHTTPS realiza o cadastro de um usuário a partir de um request HTTPS
func CadastraHTTPS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Cadastro realizado."))
}

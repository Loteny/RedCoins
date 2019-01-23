package cadastro

import "net/http"

// CadastraHTTPS realiza o cadastro de um usuário a partir de um request HTTPS
func CadastraHTTPS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Cadastro realizado."))
}

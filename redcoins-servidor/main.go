package main

import (
	"log"
	"net/http"
)

// addr especifica o endereço (incluindo porta) do servidor
var addr = "0.0.0.0:8080"

// pem e key especificam a localização em disco
var pem = "../src/github.com/loteny/redcoins/certs/server.pem"
var key = "../src/github.com/loteny/redcoins/certs/server.key"

// escutaConexoes estabelece as rotas do servidor, coloca o servidor em modo de
// escuta por novos clientes e os envia para a função apropriada para processar
// o pedido do cliente e respondê-lo
func escutaConexoes() {
	estabeleceRotas()
	err := http.ListenAndServeTLS(addr, pem, key, nil)
	if err != nil {
		log.Fatalf("Erro na função ListenAndServeTLS: %s", err)
	}
}

// estabeleceRotas define as funções a serem chamadas para cada rota que for
// pedida do servidor
func estabeleceRotas() {
	http.HandleFunc("/cadastro", RotaCadastro)
	// http.HandleFunc("/compra", nil)
	// http.HandleFunc("/venda", nil)
}

func main() {
	escutaConexoes()
}

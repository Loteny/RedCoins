package main

import (
	"log"
	"net/http"
	"os"

	"github.com/loteny/redcoins/database"
)

// addrHTTPS e addrHTTP especificam o endereço (incluindo porta) do servidor
// utilizando HTTPS e HTTP, respectivamente
var addrHTTPS string
var addrHTTP string

// pem e key especificam a localização em disco dos arquivos para TLS (se o path
// for relativo, é em relação ao executável do servidor)
var pem string
var key string

// config é a estrutura com as configurações do servidor
type config struct {
	RedcoinsServidor struct {
		AddrHTTPS string `json:"addrHttps"`
		AddrHTTP  string `json:"addrHttp"`
		Pem       string `json:"pem"`
		Key       string `json:"key"`
	} `json:"redcoins-servidor"`
}

func init() {
	// Inicializa as configurações da package com as variáveis de ambiente
	addrHTTPS = os.Getenv("REDCOINS_SV_ADDRHTTPS")
	addrHTTP = os.Getenv("REDCOINS_SV_ADDRHTTP")
	pem = os.Getenv("REDCOINS_SV_PEM")
	key = os.Getenv("REDCOINS_SV_KEY")
}

// escutaConexoes estabelece as rotas do servidor, coloca o servidor em modo de
// escuta por novos clientes e os envia para a função apropriada para processar
// o pedido do cliente e respondê-lo.
// Se o programa foi chamado com --sem-tls, o servidor se comunica com HTTP ao
// invés de HTTPS
func escutaConexoes() {
	estabeleceRotas()
	if len(os.Args) > 1 && os.Args[1] == "--sem-tls" {
		err := http.ListenAndServe(addrHTTP, nil)
		if err != nil {
			log.Fatalf("Erro na função ListenAndServeTLS: %s", err)
		}
		return
	}
	err := http.ListenAndServeTLS(addrHTTPS, pem, key, nil)
	if err != nil {
		log.Fatalf("Erro na função ListenAndServeTLS: %s", err)
	}
}

// estabeleceRotas define as funções a serem chamadas para cada rota que for
// pedida do servidor
func estabeleceRotas() {
	http.HandleFunc("/cadastro", RotaCadastro)
	http.HandleFunc("/transacoes/compra", RotaCompra)
	http.HandleFunc("/transacoes/venda", RotaVenda)
	http.HandleFunc("/relatorios/data", RotaRelatorioDia)
	http.HandleFunc("/relatorios/usuario", RotaRelatorioUsuario)
}

func main() {
	if err := database.CriaDatabase(); err != nil {
		log.Fatalf("Erro ao tentar banco de dados: %s", err)
	}
	escutaConexoes()
}

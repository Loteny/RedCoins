package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/loteny/redcoins/database"
)

// addr especifica o endereço (incluindo porta) do servidor
var addr string

// pem e key especificam a localização em disco dos arquivos para TLS (se o path
// for relativo, é em relação ao executável do servidor)
var pem string
var key string

// config é a estrutura com as configurações do servidor
type config struct {
	RedcoinsServidor struct {
		Addr string `json:"addr"`
		Pem  string `json:"pem"`
		Key  string `json:"key"`
	} `json:"redcoins-servidor"`
}

func init() {
	// Inicializa as configurações do módulo com o arquivo config.json
	arquivoConfig, err := os.Open("./config.json")
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo de configurações do servidor: %s", err)
	}
	var c config
	if err := json.NewDecoder(arquivoConfig).Decode(&c); err != nil {
		log.Fatalf("Erro ao ler configurações do servidor: %s", err)
	}
	addr = c.RedcoinsServidor.Addr
	pem = c.RedcoinsServidor.Pem
	key = c.RedcoinsServidor.Key
}

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
	http.HandleFunc("/compra", RotaCompra)
	http.HandleFunc("/venda", RotaVenda)
	http.HandleFunc("/relatorio/data", RotaRelatorioDia)
	http.HandleFunc("/relatorio/usr", RotaRelatorioUsuario)
}

func main() {
	if err := database.CriaDatabase(); err != nil {
		log.Fatalf("Erro ao tentar banco de dados: %s", err)
	}
	escutaConexoes()
}

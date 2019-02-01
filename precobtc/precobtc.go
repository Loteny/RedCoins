// Package precobtc serve para adquirir o preço de Bitcoins atualizado
package precobtc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// Erros possíveis internos do package
var (
	ErrStatusCodeInesperado = errors.New("resposta inesperada de CoinMarketCap")
)

// Variáveis para o cache do preço da Bitcoin. Armazenar essas variáveis em
// cache evita pedidos constantes HTTP para um outro servidor. A data do cache
// é armazenada no formato "YYYY-MM-DD-HH"
var (
	precoCache float64
	dataCache  string
)

// respostaJSON segue os padrões JSON da API para obter o preço de Bitcoins
// utilizada para obter o preço da Bitcoin em BRL
type respostaJSON struct {
	Data struct {
		Quotes struct {
			BRL struct {
				Price float64 `json:"price"`
			} `json:"BRL"`
		} `json:"quotes"`
	} `json:"data"`
}

// adquirePreco retorna o preço da Bitcoin de uma estrutura 'respostaJSON'
func (r respostaJSON) adquirePreco() float64 {
	return r.Data.Quotes.BRL.Price
}

// PrecoUnidade retorna o preço de uma unidade de Bitcoin em BRL
func PrecoUnidade() (float64, error) {
	// Primeiro, checa se o cache está atualizado.
	// O cache é considerado atualizado quando a data de atualização no formato
	// "YYYY-MM-DD-HH" é igual à data atual, ou seja, qualquer mudança no ano,
	// mês, dia ou hora acarreta em um cache desatualizado.
	dataAtual := time.Now().Format("2006-01-02-15")
	if dataCache == dataAtual {
		return precoCache, nil
	}

	// Request para o site externo
	rHTTP, err := http.Get("https://api.coinmarketcap.com/v2/ticker/1/?convert=BRL")
	if err != nil {
		return 0, err
	} else if rHTTP.StatusCode != http.StatusOK {
		return 0, ErrStatusCodeInesperado
	}

	// Tenta encaixar a resposta na estrutura JSON conhecida
	rBytes, err := ioutil.ReadAll(rHTTP.Body)
	if err != nil {
		return 0, err
	}
	rJSON := respostaJSON{}
	if err := json.Unmarshal(rBytes, &rJSON); err != nil {
		return 0, err
	}

	// Atualiza o cache
	precoCache = rJSON.adquirePreco()
	dataCache = dataAtual

	return precoCache, nil
}

// Preco retorna o preço de uma quantidade de Bitcoins em BRL
func Preco(qtd float64) (float64, error) {
	preco, err := PrecoUnidade()
	if err != nil {
		return 0, err
	}
	return preco * qtd, nil
}

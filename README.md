# RedCoins

RedCoins é uma API para exchange de RedCoins.

## Instruções

A API funciona apenas com HTTPS e possui certificado TLS auto-assinado.

Para utilizar, pode-se utilizar o Docker do projeto. No diretório root deste repositório, o arquivo Dockerfile pode gerar a imagem Docker para rodar o projeto:

```bash
docker build -t leoschsenna/redcoins-sv .
```

Com a imagem gerada, basta executar o projeto com o arquivo docs/docker-compose.yml:

```bash
docker swarm init
docker stack -c docs/docker-compose.yml rds
```

Para parar o servidor:

```bash
docker stack rm rds
docker swarm leave --force
```

Alternativamente, dado que exista um servidor MySQL em execução para ser utilizado e as ferramentas de Go devidamente instaladas, é possível executar o projeto começando pela instalação das dependências:

```bash
go get github.com/go-sql-driver/mysql
go get golang.org/x/crypto/bcrypt
```

E a instalação do servidor:

```bash
go install github.com/loteny/redcoins/redcoins-servidor
```

Também é necessário configurar o servidor adequadamente. Para isso, deve-se criar uma cópia do arquivo config_sample.json, renomeá-la para config.json e colocá-la no mesmo diretório que o executável do servidor. Depois, basta executar o servidor normalmente. O servidor é capaz de criar o banco de dados e suas tabelas durante sua inicialização.

Para enviar requests para o servidor, pode-se utilizar os comandos de cURL gerados pelo Swagger a partir da documentação, porém, é necessário o acréscimo do parâmetro ```-k``` para aceitar conexões inseguras, já que o servidor possui certificado TLS auto-assinado.

A API está documentada no arquivo docs/documentacao_api.yaml. Um diagrama do banco de dados está presente em docs/diagrama_db.png.

## Comandos cURL

Aqui estão listados alguns comandos de cURL para testes. Parâmetros em {chaves} devem ser substituídos pelos valores reais.

### Cadastro de usuário

```bash
curl -X POST "https://{link do servidor}/cadastro" -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "email={e-mail do usuário}&senha={senha do usuário}&nome={nome do usuário}&nascimento={data de nascimento do usuário} -k"
```

### Compra de bitcoins

```bash
curl -X POST "https://{link do servidor}/transacoes/compra" -H "accept: application/json" -H "authorization: Basic {autenticação do usuário}" -H "Content-Type: application/x-www-form-urlencoded" -d "qtd={quantidade a ser comprada}&data={data da transação} -k"
```

### Venda de bitcoins

```bash
curl -X POST "https://{link do servidor}/transacoes/venda" -H "accept: application/json" -H "authorization: Basic {autenticação do usuário}" -H "Content-Type: application/x-www-form-urlencoded" -d "qtd={quantidade a ser vendida}&data={data da transação} -k"
```

### Relatório de transações por usuário

```bash
curl -X GET "https://{link do servidor}/relatorios/usuario?email={e-mail do usuário}" -H "accept: application/json" -H "authorization: Basic {autenticação do usuário} -k"
```

### Relatório de transações por data

```bash
curl -X GET "https://{link do servidor}/relatorios/data?data={data das transações}" -H "accept: application/json" -H "authorization: Basic {autenticação do usuário} -k"
```
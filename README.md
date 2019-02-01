# RedCoins

RedCoins é uma API para exchange de RedCoins.

## Instruções

A API funciona com HTTPS e possui certificado TLS auto-assinado. Também é possível iniciar o servidor em HTTP com o primeiro argumento para execução ```--sem-tls```.

Para utilizar, pode-se utilizar o Docker do projeto. No diretório root deste repositório, o arquivo Dockerfile pode gerar a imagem Docker para rodar o projeto:

```bash
docker build -t leoschsenna/redcoins-sv:https . -f docs/Dockerfile_HTTPS
docker build -t leoschsenna/redcoins-sv:http . -f docs/Dockerfile_HTTP
```

Ou pode-se adquirir o projeto em cloud:

```bash
docker pull leoschsenna/redcoins-sv:https
docker pull leoschsenna/redcoins-sv:http
```

Também é necessário adquirir o repositório de MySQL:

```bash
docker pull mysql
```

Com a imagem gerada, basta executar o projeto com o arquivo docs/docker-compose.yml:

```bash
docker swarm init
docker stack deploy -c docs/docker-compose.yml rds
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

Também é necessário configurar o servidor adequadamente. Para isso, podem ser utilizadas variáveis do ambiente. As variáveis de ambiente do servidor começam com REDCOINS_, e são acompanhadas de um outro prefixo indicando seu package (SV_ para o servidor e DB_ para o banco de dados). Abaixo estão listadas todas as variáveis de ambiente que o projeto usa em formato de um exemplo de como configurá-las utilizando Windows:

```bash
SET REDCOINS_SV_ADDRHTTPS=0.0.0.0:443
SET REDCOINS_SV_ADDRHTTP=0.0.0.0:80
SET REDCOINS_SV_PEM=../src/github.com/loteny/redcoins/certs/server.pem
SET REDCOINS_SV_KEY=../src/github.com/loteny/redcoins/certs/server.key
SET REDCOINS_DB_USR=root
SET REDCOINS_DB_SENHA=tvMv2gjAcH5a
SET REDCOINS_DB_DBNOME=redcoins
SET REDCOINS_DB_TESTEDBNOME=redcoins_teste
SET REDCOINS_DB_DBADDR=host.docker.internal:3306
```

O servidor é capaz de criar o banco de dados e suas tabelas durante sua inicialização. Portanto, é necessário apenas que o servidor seja configurado para utilizar um usuário com permissões para criar e gerenciar banco de dados.

Para enviar requests para o servidor, pode-se utilizar os comandos de cURL gerados pelo Swagger a partir da documentação, porém, é necessário o acréscimo do parâmetro ```-k``` para aceitar conexões inseguras, já que o servidor possui certificado TLS auto-assinado.

A API está documentada no arquivo docs/documentacao_api.yaml. Um diagrama do banco de dados está presente em docs/diagrama_db.png.

## Instruções para testes

Para executar os testes de packages que dependem de banco de dados, é necessário adicionar o arquivo config.json dentro de sua pasta para realizar os testes da package. Por exemplo, para executar os testes da package cadastro, é necessário que haja um arquivo config.json com as configurações corretas no diretório do módulo (ficando lado a lado com o arquivo cadastro.go, por exemplo). As packages que dependem do config.json no mesmo diretório para realizar os testes são:

- cadastro
- database
- redcoins-servidor (main)
- transacao

A package database não pode executar todos os testes da package simultaneamente. Os testes no arquivo schema_test.go devem ser todos executados individualmente. Os outros testes podem ser executados simultaneamente.

## Comandos cURL

Aqui estão listados alguns comandos de cURL para testes. Parâmetros em {chaves} devem ser substituídos pelos valores reais. Cada comando possui dois exemplos: por link, onde os dados da Basic Auth vão no path do pedido onde caracteres especiais devem estar encodados com percent encode (por exemplo, @ se torna %40), e por parâmetro, onde os credenciais devem estar em base64 (exceto o cadastro de usuário, que não requer autenticação).

### Cadastro de usuário

```bash
curl -X POST "https://{link do servidor}/cadastro" -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "email={e-mail do usuário}&senha={senha do usuário}&nome={nome do usuário}&nascimento={data de nascimento do usuário}" -k -v
```

Exemplo:

```bash
curl -X POST "https://127.0.0.1/cadastro" -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "email=teste%40gmail%2Ecom&senha=password123&nome=Usu%C3%A1rio%20Teste&nascimento=1990%2D03%2D08" -k -v
```

### Compra de bitcoins

```bash
curl -X POST "https://{e-mail}:{senha}@{link do servidor}/transacoes/compra" -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "qtd={quantidade a ser comprada}&data={data da transação}" -k -v
curl -X POST "https://{link do servidor}/transacoes/compra" -H "accept: application/json" -H "authorization: Basic {autenticação do usuário}" -H "Content-Type: application/x-www-form-urlencoded" -d "qtd={quantidade a ser comprada}&data={data da transação}" -k -v
```

Exemplo:

```bash
curl -X POST "https://teste%40gmail%2Ecom:password123@127.0.0.1/transacoes/compra" -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "qtd=0%2E002&data=2019%2D01%2D22" -k -v
```

### Venda de bitcoins

```bash
curl -X POST "https://{e-mail}:{senha}@{link do servidor}/transacoes/venda" -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "qtd={quantidade a ser vendida}&data={data da transação}" -k -v
curl -X POST "https://{link do servidor}/transacoes/venda" -H "accept: application/json" -H "authorization: Basic {autenticação do usuário}" -H "Content-Type: application/x-www-form-urlencoded" -d "qtd={quantidade a ser vendida}&data={data da transação}" -k -v
```

Exemplo:

```bash
curl -X POST "https://teste%40gmail%2Ecom:password123@127.0.0.1/transacoes/venda" -H "accept: application/json" -H "Content-Type: application/x-www-form-urlencoded" -d "qtd=0%2E002&data=2019%2D01%2D22" -k -v
```

### Relatório de transações por usuário

```bash
curl -X GET "https://{e-mail}:{senha}@{link do servidor}/relatorios/usuario?email={e-mail do usuário}" -H "accept: application/json" -k -v
curl -X GET "https://{link do servidor}/relatorios/usuario?email={e-mail do usuário}" -H "accept: application/json" -H "authorization: Basic {autenticação do usuário}" -k -v
```

Exemplo:

```bash
curl -X GET "https://teste%40gmail%2Ecom:password123@127.0.0.1/relatorios/usuario?email=teste%40gmail%2Ecom" -H "accept: application/json" -k -v
```

### Relatório de transações por data

```bash
curl -X GET "https://{e-mail}:{senha}@{link do servidor}/relatorios/data?data={data das transações}" -H "accept: application/json" -k -v
curl -X GET "https://{link do servidor}/relatorios/data?data={data das transações}" -H "accept: application/json" -H "authorization: Basic {autenticação do usuário}" -k -v
```

Exemplo:

```bash
curl -X GET "https://teste%40gmail%2Ecom:password123@127.0.0.1/relatorios/data?data=2019%2D01%2D22" -H "accept: application/json" -k -v
```

## Certificados TLS

Os certificados TLS foram gerados utilizando o seguinte formato de comando:

```bash
openssl req -new -nodes -x509 -out certs/server.pem -keyout certs/server.key -days 3650 -subj "//C=BR\ST=ES\L=Cidade\O=Organização\OU=IT\emailAddress=email@gmail.com"
```

Caso seja necessário ou desejado gerar outros certificados TLS, esse comando pode ser utilizado.
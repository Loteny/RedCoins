FROM golang

# Dependências
RUN go get github.com/go-sql-driver/mysql
RUN go get golang.org/x/crypto/bcrypt

ADD . /go/src/github.com/loteny/redcoins
ADD config_sample.json /go/bin/config.json

RUN go install github.com/loteny/redcoins/redcoins-servidor
WORKDIR /go/bin
ENTRYPOINT /go/bin/redcoins-servidor
EXPOSE 443
package cadastro

import (
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/loteny/redcoins/erros"
)

// email verifica se o e-mail é válido (formato regex /.+@.+/)
func email(email string) error {
	return validacaoMatchSimples(email, "^.+@.+$", ErrEmailInvalido)
}

// senha verifica se a senha possui pelo menos 6 caracteres e no máximo 64
func senha(senha string) error {
	return validacaoMatchSimples(senha, "^.{6,64}$", ErrSenhaInvalida)
}

// nome verifica se o campo não está vazio e se não excede 128 caracteres
func nome(nome string) error {
	if utf8.RuneCountInString(nome) <= 0 || utf8.RuneCountInString(nome) > 128 {
		return ErrNomeInvalido
	}
	return nil
}

// nascimento verifica se a data está no formato válido e se a data é passada.
// Problemas com fuso horário não são importantes, visto que só seriam
// possivelmente bloqueados datas de nascimentos de recém-nascidos por problemas
// de fuso horário. Além disso, o 'Time' resultante da data de entrada estará no
// início do dia (00h00m00...), portanto, há uma "margem de erro" nessa função,
// mas essa margem é pequena (algumas horas, possivelmente alguns dias,
// dependendo de mudanças específicas de fusos horários) e pode ser ignorada
func nascimento(data string) error {
	dataTime, err := time.Parse("2006-01-02", data)
	if err != nil {
		return ErrNascimentoInvalido
	}
	agora := time.Now()
	if dataTime.After(agora) {
		return ErrNascimentoInvalido
	}
	return nil
}

// validacaoMatchSimples executa uma validação básica com um regex passado como
// argumento. Retorna o erro gerado pela função do regex, caso houve algum, ou o
// erro passado por argumento para essa função no caso de o regex não bater com
// o dado passado.
func validacaoMatchSimples(s string, reg string, e error) error {
	matched, err := regexp.MatchString(reg, s)
	if err != nil {
		return erros.CriaInternoPadrao(err)
	}
	if !matched {
		return e
	}
	return nil
}

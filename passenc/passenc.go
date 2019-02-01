package passenc

import (
	"golang.org/x/crypto/bcrypt"
)

// GeraHashed gera uma array de bytes contendo as informações do algoritmo de
// encriptação, o salt utilizado e a senha hashed.
// O resultado possui entre 59 e 60 bytes.
// Até 50 bytes, a senha é garantida de funcionar como esperado. Além dessa
// quantidade, problemas indesejados como truncamento de senha podem ocorrer.
func GeraHashed(senha []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(senha, bcrypt.MinCost)
	if err != nil {
		return []byte{}, err
	}
	return hash, nil
}

// VerificaSenha checa se a senha e a senha hashed passada são equivalentes. O
// retorno da função indica se são.
func VerificaSenha(senha, senhaHashed []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(senhaHashed, senha)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

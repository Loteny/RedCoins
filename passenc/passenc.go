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

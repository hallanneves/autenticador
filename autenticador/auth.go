package autenticador

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	models "github.com/hallanneves/autenticador/models"
)

type cache struct {
	tempo    time.Time
	password string
}

var contadorConexao = 0

var listaCache = make(map[string]*cache)
var mutexLista sync.Mutex

func senhaHash(senha string) string {
	//Converte para Sha256
	hash := sha256.New()
	hash.Write([]byte(senha))
	passwordHash := hash.Sum(nil)
	return hex.EncodeToString(passwordHash)
}

//ValidaAutenticacao retorna se as credencias são validas
func ValidaAutenticacao(credencias *models.Credenciais) (int, error) {

	mutexLista.Lock()
	if credencialCache, existeNoCache := listaCache[*credencias.Usuario]; existeNoCache {
		mutexLista.Unlock()

		if time.Now().Sub(credencialCache.tempo).Minutes() < tempoCache {
			mutexLista.Lock()
			delete(listaCache, *credencias.Usuario)
			mutexLista.Unlock()
		} else {
			if senhaHash(*credencias.Senha) == credencialCache.password {
				return 200, nil
			}
			return 401, nil
		}
	} else {
		mutexLista.Unlock()
	}

	//busca o usuário e a senha no banco de dados
	var password string
	db := distribuidorDeConexao()
	err := db.conexao.QueryRow("SELECT password FROM credenciais WHERE username= ?", credencias.Usuario).Scan(&password)
	if err == sql.ErrNoRows {
		return 401, nil
	} else if err != nil {
		return 500, errors.New("erro no mysql: " + err.Error() + " db: " + db.conf.Host)
	}
	//coloca na cache
	mutexLista.Lock()
	listaCache[*credencias.Usuario] = &cache{tempo: time.Now(), password: password}
	mutexLista.Unlock()

	//verifica a senha
	if senhaHash(*credencias.Senha) == password {
		return 200, nil
	}
	return 401, nil

}

//distribuidorDeConexao distribui as requisições entre as basses de dados cadastradas
//Como não é critico o envio consecutivo de mais de uma requisição ao mesmo banco não foi tratada a concorencia pela variável contadorConexao
func distribuidorDeConexao() sqlConf {
	db := mysqlpool[contadorConexao]
	if contadorConexao < len(mysqlpool)-1 {
		contadorConexao = contadorConexao + 1
	} else {
		contadorConexao = 0
	}
	return db
}

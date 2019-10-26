package autenticador

import (
	"errors"
	"sync"
	"time"

	conf "github.com/hallanneves/autenticador/conf"

	"database/sql"
	//driver mysql do sql
	_ "github.com/go-sql-driver/mysql"
)

type sqlConf struct {
	conexao *sql.DB
	conf    conf.Mysql
}

var mysqlpool []sqlConf
var globalMutex sync.RWMutex

// InicializaMysql inicializa conexao com o cluster Mysql
func InicializaMysql() error {
	globalMutex.Lock()
	if mysqlpool != nil {
		globalMutex.Unlock()
		return nil
	}
	for _, element := range conf.ConfigConecta.MySQLPool {
		conn, err := sql.Open("mysql", element.User+":"+element.Pass+"@("+element.Host+":"+element.Port+")/autenticador")

		if err != nil {
			globalMutex.Unlock()
			return err
		}
		conn.SetConnMaxLifetime(time.Minute * 5)
		var novoElemento = sqlConf{conexao: conn, conf: element}

		mysqlpool = append(mysqlpool, novoElemento)
	}

	globalMutex.Unlock()

	return nil
}

// VerificaMysql verifica se todos os bancos de dados estao online
func VerificaMysql() error {
	var err error

	//verifica se as bases do mysql estao ativas
	for _, db := range mysqlpool {
		err = db.conexao.Ping()
		if err != nil {
			return errors.New("mysql PING: " + err.Error() + " db: " + db.conf.Host)
		}
		_, err := db.conexao.Query("SELECT 1")
		if err != nil {
			return errors.New("mysql SELECT 1: " + err.Error() + " db: " + db.conf.Host)
		}
	}
	return nil

}

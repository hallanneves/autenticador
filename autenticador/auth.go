package autenticador

import (
	"database/sql"

	models "github.com/hallanneves/autenticador/models"
)

//ValidaAutenticacao retorna se as credencias são validas
func ValidaAutenticacao(credencias *models.Credenciais) (int, error) {

	var status int = 401
	var password string
	err := mysqlstatus.QueryRow("SELECT password FROM credencias WHERE username= ?", credencias.Usuario).Scan(&password)
	if err == sql.ErrNoRows {
		status = 401
	}
	return status, err

}

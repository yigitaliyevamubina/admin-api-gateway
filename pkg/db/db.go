package db

import (
	"database/sql"
	"fmt"
	"myproject/admin-api-gateway/config"
	_ "github.com/lib/pq" 
)

func ConnectToDB(cfg config.Config) (*sql.DB, func(), error) {
	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	connDB, err := sql.Open("postgres", psqlString)
	if err != nil {
		return nil, nil, err
	}

	cleanUpFunc := func() {
		connDB.Close()
	}

	return connDB, cleanUpFunc, nil
}

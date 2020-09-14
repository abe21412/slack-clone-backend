package providers

import (
	"context"
	"log"

	"github.com/abe21412/slack-clone-backend/src/models"
	"github.com/abe21412/slack-clone-backend/src/util/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool = db.GetPool()

func CreateUser(user *models.User) (string, error) {
	sql := `insert into users(id, display_name, first_name, last_name, email) values($1, $2, $3, $4, $5) 
			on conflict on constraint users_pkey do update 
			set (id, display_name, first_name, last_name, email) = (EXCLUDED.id, EXCLUDED.display_name, EXCLUDED.first_name, EXCLUDED.last_name, EXCLUDED.email)
			returning id;`
	log.Println(user)
	row := pool.QueryRow(context.Background(), sql, user.ID, user.DisplayName, user.FirstName, user.LastName, user.Email)
	var userID string
	if err := row.Scan(&userID); err != nil {
		log.Println(err.Error())
		return "", err
	}
	return userID, nil
}

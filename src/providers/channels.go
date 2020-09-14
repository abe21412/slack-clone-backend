package providers

import (
	"context"
	"log"

	"github.com/abe21412/slack-clone-backend/src/models"
	"github.com/delaemon/go-uuidv4"
)

func CreateChannel(channel *models.Channel, userID string) (*models.Channel, error) {
	channelID, err := uuidv4.Generate()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	createChannelSQL := "insert into channels (id, name, workspace_id, public, description) values ($1, $2, $3, $4, $5);"
	var subscriptionsSQL string
	var subscriptionParams []interface{}
	if channel.Public {
		subscriptionsSQL = `insert into subscriptions (user_id, workspace_id, channel_id) 
							select user_id, workspace_id, $1 as channel_id from members where workspace_id = $2;`
		subscriptionParams = []interface{}{channelID, channel.WorkspaceID}
	} else {
		subscriptionsSQL = "insert into subscriptions (user_id, workspace_id, channel_id) values ($1, $2, $3);"
		subscriptionParams = []interface{}{userID, channel.WorkspaceID, channelID}
	}
	tx, err := pool.Begin(context.Background())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	_, err = tx.Exec(context.Background(), createChannelSQL, channelID, channel.Name, channel.WorkspaceID, channel.Public, channel.Description)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	_, err = tx.Exec(context.Background(), subscriptionsSQL, subscriptionParams...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	channel.ID = channelID
	log.Println(channel)
	return channel, nil
}

func GetChannels(userID, workspaceID string) ([]models.Channel, error) {
	var channels []models.Channel
	sql := `select c.id, c.name, c.workspace_id, c.public, c.dm, coalesce(c.description, '') from channels c
			 join subscriptions s on s.workspace_id = c.workspace_id and s.channel_id = c.id 
			where s.user_id = $1 and s.workspace_id = $2;`
	rows, err := pool.Query(context.Background(), sql, userID, workspaceID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var c models.Channel
		err := rows.Scan(&c.ID, &c.Name, &c.WorkspaceID, &c.Public, &c.DM, &c.Description)
		if err != nil {
			return nil, err
		}
		channels = append(channels, c)
	}
	return channels, nil
}

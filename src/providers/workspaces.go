package providers

import (
	"context"
	"log"

	"github.com/abe21412/slack-clone-backend/src/models"
	"github.com/delaemon/go-uuidv4"
)

func CreateWorkspace(workspace *models.Workspace) (string, error) {
	workspaceID, err := uuidv4.Generate()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	workspaceSQL := `insert into workspaces(id, name, owner) values($1, $2, $3);`
	membersSQL := `insert into members(user_id, workspace_id) values ($1, $2);`
	log.Println(workspace)
	tx, err := pool.Begin(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	_, err = tx.Exec(context.Background(), workspaceSQL, workspaceID, workspace.Name, workspace.OwnerID)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	_, err = tx.Exec(context.Background(), membersSQL, workspace.OwnerID, workspaceID)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	generalChannel := &models.Channel{Name: "general", WorkspaceID: workspaceID, Public: true}
	CreateChannel(generalChannel, workspace.OwnerID)
	return workspaceID, nil
}

func GetWorkspaces(userID string) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	sql := `select w.id, w.name, w.owner from workspaces w join members m on w.id = m.workspace_id where m.user_id = $1;`
	rows, err := pool.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var w models.Workspace
		err := rows.Scan(&w.ID, &w.Name, &w.OwnerID)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, w)
	}
	return workspaces, nil
}

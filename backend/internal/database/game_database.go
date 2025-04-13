package database

import (
	"context"
	"fmt"
)

func (pd *PostgresDriver) GetGameConfig(ctx context.Context) (map[string]interface{}, error) {

	gameConfig := make(map[string]interface{})

	rows, err := pd.pool.Query(ctx, "SELECT property, value FROM gameconfig")

	if err != nil {
		return nil, fmt.Errorf("error fetching gameconfig: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var property string
		var value string
		if err := rows.Scan(&property, &value); err != nil {
			return nil, fmt.Errorf("error scanning gameconfig: %v", err)
		}

		gameConfig[property] = value
	}

	return gameConfig, nil

}

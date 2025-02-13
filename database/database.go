package database

import "github.com/brightside-dev/ronin-fitness-be/database/client"

func New() client.DatabaseService {
	return client.NewMySQL()
}

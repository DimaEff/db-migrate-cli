package presets

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Connection struct {
	ID        int
	Name      string
	SourceURL string
	TargetURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ConnectToCliDb(dataSourceFilePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceFilePath)
	if err != nil {
		return nil, fmt.Errorf("connect to sqlite3: %w", err)
	}
	return db, nil
}

// GetConnections retrieves all saved sources and targets urls
// of connections from the database
func GetConnections(db *sql.DB) ([]Connection, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := db.Query(`
  select id,
  name,
        source_connection_url,
        target_connection_url,
        created_at,
        updated_at
        from connections;
        `)

	if err != nil {
		return nil, fmt.Errorf("get connections query: %w", err)
	}
	defer rows.Close()

	connections := make([]Connection, 0, 10)

	for rows.Next() {
		var connection Connection
		err := rows.Scan(
			&connection.ID,
			&connection.Name,
			&connection.SourceURL,
			&connection.TargetURL,
			&connection.CreatedAt,
			&connection.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("get connections scan: %w", err)
		}

		connections = append(connections, connection)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get connections query: %w", err)
	}

	return connections, nil
}

// SaveConnectionsURLs saves a new connection`s source and target urls
// in the database and returns the saved connection
func SaveConnectionsURLs(db *sql.DB, name string, sourceURL string, targetURL string) (*Connection, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var connection Connection
	err := db.QueryRow(`
  insert into connections (name, source_connection_url, target_connection_url)
  values (?, ?, ?)
  returning id, name, source_connection_url, target_connection_url, created_at, updated_at;
  `, name, sourceURL, targetURL).Scan(
		&connection.ID,
		&connection.Name,
		&connection.SourceURL,
		&connection.TargetURL,
		&connection.CreatedAt,
		&connection.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("CreateConnection query: %w", err)
	}

	return &connection, nil
}

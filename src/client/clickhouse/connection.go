// repository/clickhouse.go
package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/ClickHouse/clickhouse-go" // Import ClickHouse driver for SQL package
)

type ClickHouseRepo struct {
	conn *sql.DB
}

// NewClickHouseClient initializes a ClickHouse client
func NewClickHouseClient() (*sql.DB, error) {
	uri := os.Getenv("CLICKHOUSE_URL")
	client, err := sql.Open("clickhouse", uri)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(); err != nil {
		fmt.Println("Failed to connect to ClickHouse:", err)
		client.Close()
		return nil, err
	}
	fmt.Println("Connected to ClickHouse")
	return client, nil
}

// NewClickHouseRepo creates a new instance of ClickHouseRepo
func NewClickHouseRepo() (*ClickHouseRepo, error) {
	client, err := NewClickHouseClient()
	if err != nil {
		return nil, err
	}
	return &ClickHouseRepo{conn: client}, nil
}

// Close closes the ClickHouse connection
func (repo *ClickHouseRepo) Close() {
	repo.conn.Close()
}

func (repo *ClickHouseRepo) GetTodoMetrics(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := repo.conn.QueryContext(ctx, `
		SELECT
			type,
			COUNT(*) AS total_tasks,
			SUM(CASE WHEN done = 1 THEN 1 ELSE 0 END) AS completed_tasks,
			SUM(CASE WHEN done = 0 THEN 1 ELSE 0 END) AS not_completed_tasks,
			SUM(effortHr) AS total_effort,
			IF(COUNT(*) = 0, 0, (SUM(CASE WHEN done = 1 THEN 1 ELSE 0 END) / COUNT(*)) * 100) AS completion_percentage
		FROM todo_app.todos
		GROUP BY type;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []map[string]interface{}
	for rows.Next() {
		var (
			taskType             string
			totalTasks           uint64
			completedTasks       uint64
			notCompletedTasks    uint64
			totalEffort          float64
			completionPercentage float64
		)
		if err := rows.Scan(&taskType, &totalTasks, &completedTasks, &notCompletedTasks, &totalEffort, &completionPercentage); err != nil {
			return nil, err
		}
		metrics = append(metrics, map[string]interface{}{
			"type":                  taskType,
			"total_tasks":           totalTasks,
			"completed_tasks":       completedTasks,
			"not_completed_tasks":   notCompletedTasks,
			"total_effort":          totalEffort,
			"completion_percentage": completionPercentage,
		})
	}
	return metrics, nil
}

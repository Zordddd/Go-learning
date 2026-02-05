package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/Zordddd/learning/taskAPI/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func ConnectDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (repo *TaskRepository) GetTasks(ctx context.Context) ([]models.Task, error) {
	query := `SELECT * FROM tasks`

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Text,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *TaskRepository) CreateTask(ctx context.Context, task models.TaskCreated) error {
	query := `INSERT INTO tasks (title, text, completed, created_at, updated_at)
		      VALUES ($1, $2, $3, $4, $5)
		      RETURNING id`
	_, err := repo.DB.ExecContext(ctx, query, task.Title, task.Text, false, time.Now(), time.Now())
	return err
}

func (repo *TaskRepository) UpdateTask(ctx context.Context, task models.TaskUpdated) error {
	query := `UPDATE tasks SET title=$2, text=$3, completed=$4, updated_at=$5 WHERE id=$1 RETURNING id`
	err := repo.DB.QueryRowContext(ctx, query, task.ID, task.Title, task.Text, task.Completed, time.Now()).Scan(&task.ID)
	return err
}

func (repo *TaskRepository) DeleteTask(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := repo.DB.ExecContext(ctx, query, id)
	return err
}

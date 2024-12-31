package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {

	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, insert)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, subject, description)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var todo model.TODO
	todo.ID = id
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	var (
		stmt *sql.Stmt
		err  error
		rows *sql.Rows
	)

	fmt.Println("prevID:", prevID, "size:", size)

	if prevID == 0 {
		stmt, err = s.db.PrepareContext(ctx, read)
	} else {
		stmt, err = s.db.PrepareContext(ctx, readWithID)
	}

	if err != nil {
		fmt.Println("PrepareContext Error", err)
		return nil, err
	}

	defer stmt.Close()

	if prevID == 0 {
		rows, err = stmt.QueryContext(ctx, size)
	} else {
		rows, err = stmt.QueryContext(ctx, prevID, size)
	}

	if err != nil {
		fmt.Println("QueryContext Error", err)
		return nil, err
	}

	defer rows.Close()

	todos := []*model.TODO{}

	for rows.Next() {
		var todo model.TODO

		err = rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

		if err != nil {
			fmt.Println("rows.Scan Error", err)
			return nil, err

		}

		fmt.Println("todo:", todo)

		todos = append(todos, &todo)

	}

	if err = rows.Err(); err != nil {
		fmt.Println("rows.Err()", err)
		return nil, err
	}

	fmt.Println("todos:", todos)

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, update)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, subject, description, id)

	if err != nil {
		return nil, err
	}

	affectedRows, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if affectedRows == 0 {
		return nil, &model.ErrNotFound{Message: "TODO not found"}
	}

	var todo model.TODO

	todo.ID = id
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	fmt.Println("input ids:", ids)

	if len(ids) == 0 {
		fmt.Println("There are no ids")
		return nil
		// return errors.New("threre are no ids")
	}

	placeholders := strings.Repeat(",?", len(ids)-1)

	fmt.Println("placeholders:", placeholders)

	query := fmt.Sprintf(deleteFmt, placeholders)

	fmt.Println("query:", query)

	stmt, err := s.db.PrepareContext(ctx, query)

	if err != nil {
		fmt.Println("PrepareContext Error:", err)
		return err
	}

	defer stmt.Close()

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	res, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		fmt.Println("ExecContext Error:", err)
		return err
	}

	affectedRows, err := res.RowsAffected()

	if err != nil {
		fmt.Println("RowsAffected Error:", err)
		return err
	}

	if affectedRows == 0 {
		return &model.ErrNotFound{Message: "there are no affectedRows"}
	}

	return nil
}

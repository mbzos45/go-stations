package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/TechBowl-japan/go-stations/model"
	"log"
	"strings"
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
	todo := model.TODO{
		Subject:     subject,
		Description: description,
	}
	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}(stmt)
	res, err := stmt.ExecContext(ctx, todo.Subject, todo.Description)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	todo.ID = id
	row := s.db.QueryRowContext(ctx, confirm, id)
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
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
	var todos []*model.TODO
	var rows *sql.Rows
	if prevID == 0 {
		stmt, err := s.db.PrepareContext(ctx, read)
		if err != nil {
			return nil, err
		}
		defer func(stmt *sql.Stmt) {
			if err := stmt.Close(); err != nil {
				log.Println(err)
			}
		}(stmt)
		rows, err = stmt.QueryContext(ctx, size)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	} else {
		stmt, err := s.db.PrepareContext(ctx, readWithID)
		if err != nil {
			return nil, err
		}
		defer func(stmt *sql.Stmt) {
			if err := stmt.Close(); err != nil {
				log.Println(err)
			}
		}(stmt)
		rows, err = stmt.QueryContext(ctx, prevID, size)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}(rows)
	for rows.Next() {
		todo := model.TODO{}
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			log.Println(err)
			return nil, err
		}
		todos = append(todos, &todo)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	if len(todos) == 0 {
		return []*model.TODO{}, nil
	}
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
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}(stmt)
	res, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, &model.ErrNotFound{}
	}
	todo := &model.TODO{
		ID: id,
	}
	err = s.db.QueryRowContext(ctx, confirm, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`
	if len(ids) != 0 {
		stmt, err := s.db.PrepareContext(ctx, fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1)))
		if err != nil {
			return err
		}
		defer func(stmt *sql.Stmt) {
			err := stmt.Close()
			if err != nil {
				log.Println(err)
			}
		}(stmt)
		var arg []interface{}
		for _, v := range ids {
			arg = append(arg, v)
		}
		res, err := stmt.ExecContext(ctx, arg...)
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return &model.ErrNotFound{}
		}
	}
	return nil
}

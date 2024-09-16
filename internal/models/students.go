package models

import (
	"database/sql"
	"errors"
)

type StudentModelInterface interface {
	Insert(id, name, major, startDate, endDate string) error
	Get(id string) (*Student, error)
}

type Student struct {
	ID        string
	Name      string
	Major     string
	StartDate string
	EndDate   string
}

type StudentModel struct {
	DB *sql.DB
}

func (m *StudentModel) Insert(id, name, major, startDate, endDate string) error {
	stmt := `INSERT INTO students (id, name, major, start_date, end_date)
	VALUES($1, $2, $3, $4, $5)`

	_, err := m.DB.Exec(stmt, id, name, major, startDate, endDate)
	if err != nil {
		return err
	}

	return nil
}

func (m *StudentModel) Get(id string) (*Student, error) {
	stmt := `SELECT id, name, major, start_date, end_date FROM students
	WHERE id = $1
	LIMIT 1;`

	row := m.DB.QueryRow(stmt, id)

	s := &Student{}

	err := row.Scan(&s.ID, &s.Name, &s.Major, &s.StartDate, &s.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

package storage

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"go_final_project/task"
)

const taskLimit = 10

func (s *Storage) Add(t *task.Task) (int, error) {
	ins, err := s.db.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)
	if err != nil {
		log.Println("can't add task:", err)
		return 0, err
	}

	id, err := ins.LastInsertId()
	if err != nil {
		log.Println("can't get task id:", err)
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetList() ([]task.Task, error) {
	rows, err := s.db.Query(
		`SELECT * 
		 FROM scheduler
		 ORDER BY date ASC
		 LIMIT :limit`,
		sql.Named("limit", taskLimit),
	)
	if err != nil {
		log.Println("can't get tasks by GetList:", err)
		return nil, err
	}

	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		t := task.Task{}

		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Println("can't get tasks by GetList:", err)
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if rows.Err() != nil {
		log.Println("can't get tasks by GetList:", rows.Err())
		return nil, rows.Err()
	}

	return tasks, nil
}

func (s *Storage) SearchTasks(search string) ([]task.Task, error) {
	log.Printf("looking for tasks with search parameter %s", search)

	var rows *sql.Rows
	var err error

	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		log.Println("Search by text")

		rows, err = s.db.Query(
			`SELECT * 
			FROM scheduler
			WHERE title LIKE :search OR comment LIKE :search
			ORDER BY date LIMIT :limit`,
			sql.Named("search", "%"+search+"%"),
			sql.Named("limit", taskLimit),
		)
	} else {
		log.Println("Search by date")

		target := date.Format("20060102")

		rows, err = s.db.Query(
			`SELECT * 
			   FROM scheduler
			   WHERE date LIKE :target
			   ORDER BY date LIMIT :limit`,
			sql.Named("target", "%"+target+"%"),
			sql.Named("limit", taskLimit),
		)
	}

	if err != nil {
		log.Println("can't find tasks by SearchTasks:", err)
		return nil, err
	}

	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		t := task.Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Println("can't find tasks by SearchTasks.", err)
			return nil, err
		}

		tasks = append(tasks, t)
	}
	log.Printf("Found %d tasks", len(tasks))

	return tasks, nil
}

func (s *Storage) GetTask(id string) (task.Task, error) {
	log.Println("Search task by ID:", id)

	row := s.db.QueryRow(
		`SELECT * 
		FROM scheduler 
		WHERE id = :id`,
		sql.Named("id", id),
	)

	var t task.Task

	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		log.Println("can't get task by id:", id, err)

		if errors.Is(err, sql.ErrNoRows) {
			return task.Task{}, errors.New(" ")
		}
		return task.Task{}, err
	}

	return t, nil
}

func (s *Storage) Update(t task.Task) error {
	log.Printf("Update task by ID:%s", t.ID)

	_, err := s.db.Exec(
		`UPDATE scheduler 
		SET date=:date, title= :title, comment= :comment, repeat= :repeat
		WHERE id= :id`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID),
	)
	if err != nil {
		log.Println("can't update task:", err)
		return err
	}

	log.Println("update successful")

	return nil
}

func (s *Storage) DeleteTask(id string) error {
	log.Println("Delete task ID:", id)

	_, err := s.db.Exec(
		`DELETE
		FROM scheduler
		WHERE id = :id`,
		sql.Named("id", id),
	)
	if err != nil {
		log.Println("can't delete task:", err)
		return errors.New("task not found")
	}

	log.Println("delete task successful")

	return nil
}

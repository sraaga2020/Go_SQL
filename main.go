package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

type Task struct {
	TaskName string
	Urgent   bool
}

func main() {
	connStr := "postgres://pqgotest:secret@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTaskTable(db)

	data := []Task{}
	rows, err := db.Query("SELECT taskname, urgent FROM task")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var taskname string
	var urgent bool

	for rows.Next() {
		rows.Scan(&taskname, &urgent)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Task{taskname, urgent})
	}

	fmt.Println(data)

	// var taskname string
	fmt.Print("Enter task name: ")
	fmt.Scanln(&taskname)

	var bool_input string
	fmt.Print("Is it an urgent task (enter true/false): ")
	fmt.Scanln(&bool_input)

	urgent, err = strconv.ParseBool(bool_input)
	if err != nil {
		fmt.Println("Invalid input.")
		return
	}

	task := Task{taskname, urgent}
	pk := insertTask(db, task)

	fmt.Printf("ID = %d\n", pk)

	// var taskname string
	// var urgent bool

	query := "SELECT taskname, urgent FROM task WHERE id = $1"
	err = db.QueryRow(query, pk).Scan(&taskname, &urgent)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows found with ID %d", pk)
		}
		log.Fatal(err)
	}

	fmt.Printf("Task: %s\n", taskname)
	fmt.Printf("Urgent: %t\n", urgent)

}

func createTaskTable(db *sql.DB) {

	query := `CREATE TABLE IF NOT EXISTS task (
                id SERIAL PRIMARY KEY,
                taskname VARCHAR(100) NOT NULL,
                urgent BOOLEAN
            )`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

}

func insertTask(db *sql.DB, task Task) int {

	query := `INSERT INTO task (taskname, urgent)
        VALUES ($1, $2) RETURNING id`

	var pk int

	err := db.QueryRow(query, task.TaskName, task.Urgent).Scan(&pk)

	if err != nil {
		log.Fatal(err)
	}

	return pk
}

package db

import (
	"fmt"
	"time"
)

const formatDate = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

/*
AddTask inserts a new task into the scheduler table.

Parameters:
- task: pointer to a Task struct containing the task details (date, title, comment, repeat).

Behavior:
- Executes an INSERT query to add the task to the database.
- Retrieves the auto-generated ID of the inserted row using LastInsertId().

Returns:
- The ID of the newly created task as int64.
- An error if the insertion or ID retrieval fails.
*/
func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return id, fmt.Errorf("error add task: %v", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error receipt id: %v", err)
	}
	return id, err
}

/*
GetTask retrieves a single task from the scheduler table by its ID.

Parameters:
- id: string representation of the task ID to look up.

Behavior:
- Executes a SELECT query to fetch the task with the matching ID.
- Scans the returned row into a local Task variable.

Returns:
- A pointer to the Task struct if found.
- nil and an error if the task is not found or the query fails.
*/
func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	var t Task
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, fmt.Errorf("task not found: %v", err)
	}
	return &t, nil
}

/*
UpdateTask updates an existing task in the scheduler table with new values.

Parameters:
- task: pointer to a Task struct containing updated details, including the ID.

Behavior:
- Executes an UPDATE query to modify the task fields based on the ID.
- Checks the number of affected rows to ensure the task actually existed.

Returns:
  - No value (error only) on success.
  - An error if the update fails, if no rows were affected (task not found),
    or if retrieving the affected row count fails.
*/
func UpdateTask(task *Task) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("error update task: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

/*
UpdateTaskDate updates only the date field of an existing task.

Parameters:
- id: string ID of the task to update.
- nextDate: new date string to set for the task.

Behavior:
- Executes an UPDATE query that changes only the date column.
- Verifies that at least one row was affected to confirm the task existed.

Returns:
- No value (error only) on success.
- An error if the update fails, no rows were affected, or row count retrieval fails.
*/
func UpdateTaskDate(id string, nextDate string) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	res, err := db.Exec(query, nextDate, id)
	if err != nil {
		return fmt.Errorf("error update task date")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

/*
DeleteTask removes a task from the scheduler table by its ID.

Parameters:
- id: string ID of the task to delete.

Behavior:
- Executes a DELETE query targeting the row with the given ID.
- Ensures that at least one row was deleted to confirm the task existed.

Returns:
- No value (error only) on success.
- An error if the deletion fails, no rows were affected, or row count retrieval fails.
*/
func DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error delete task: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

/*
Tasks retrieves a list of tasks from the scheduler table with optional filtering.

Parameters:
- limit: maximum number of tasks to return.
- search: optional search term that can be:
  - empty: returns the most recent tasks sorted by date (ascending).
  - a valid date string (in DD.MM.YYYY format): returns tasks for that specific date.
  - any other string: returns tasks where the title or comment contains the search term.

Behavior:
- Builds a dynamic SQL query based on the search type.
- Uses parameterized queries to prevent SQL injection.
- Iterates over the result rows and populates a slice of Task pointers.
- Handles errors during row scanning and final row iteration.

Returns:
- A slice of pointers to Task structs matching the criteria.
- An error if query execution, row scanning, or iteration fails.
*/
func Tasks(limit int, search string) ([]*Task, error) {
	var query string
	var args []any

	switch {
	case search == "":
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
		args = append(args, limit)
	case checkIsDate(search):
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`
		t, _ := time.Parse("02.01.2006", search)
		dbDate := t.Format(formatDate)

		args = append(args, dbDate, limit)
	default:
		query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?`

		searchParam := "%" + search + "%"
		args = append(args, searchParam, searchParam, limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return []*Task{}, fmt.Errorf("error receipt task: %v", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scaning rows: %v", err)
		}
		tasks = append(tasks, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return tasks, nil
}

/*
checkIsDate checks whether a given string can be parsed as a date in DD.MM.YYYY format.

Parameters:
- str: the string to validate.

Returns:
- true if the string is a valid date in the expected format.
- false otherwise.
*/
func checkIsDate(str string) bool {
	_, err := time.Parse("02.01.2006", str)
	return err == nil
}

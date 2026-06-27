package api

import (
	"Graduate-work/pkg/db"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
SignInHandler handles user sign-in by validating the provided password against
the environment variable TODO_PASSWORD. On success, it generates a JWT token
with a hash of the password as a claim and returns it to the client.

Steps:
  - Decodes the JSON request body containing the password.
  - Compares the received password with the configured password.
  - If valid, computes a SHA-256 hash of the password, creates JWT claims with
    the hash, and signs the token using the password as the secret.
  - Returns the token in JSON format on success or an appropriate error
    response on failure.
*/
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	type AuthReq struct {
		Password string `json:"password"`
	}

	var (
		req         AuthReq
		tokenString string
		errSigned   error
	)

	pass := os.Getenv("TODO_PASSWORD")

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "issue filling error"})
		return
	}

	if req.Password != pass {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid password"})
		return
	} else {
		h := sha256.Sum256([]byte(req.Password))
		passHash := hex.EncodeToString(h[:])

		claims := jwt.MapClaims{
			"hash": passHash,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, errSigned = token.SignedString([]byte(pass))
		if errSigned != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error hash password"})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

/*
NextDayHandler calculates the next occurrence date based on the provided
repetition rule. It accepts optional parameters for the base date and the
current time.

Parameters:
- date: optional base date string.
- repeat: required repetition rule (e.g., "daily", "weekly").
- now: optional current time; if omitted, uses the server's current time.

Returns:
  - The calculated next date string on success.
  - An error response if the 'repeat' parameter is missing or if the 'now'
    parameter has an invalid date format.
*/
func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")
	nowParam := r.FormValue("now")

	if repeatParam == "" {
		http.Error(w, "missing repeat parameter", http.StatusBadRequest)
		return
	}

	var nowTime time.Time
	var err error

	if nowParam == "" {
		nowTime = time.Now()
	} else {
		nowTime, err = time.Parse(FormatDate, nowParam)
		if err != nil {
			http.Error(w, "Invalid 'now' date format", http.StatusBadRequest)
			return
		}
	}

	nextDateStr, err := NextDate(nowTime, dateParam, repeatParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	fmt.Fprint(w, nextDateStr)
}

/*
AddTaskHandler adds a new task to the database. It expects a JSON payload
with task details.

Validation steps:
- Ensures the request body is valid JSON.
- Checks that the task title is provided.
- Validates and adjusts the task date using checkDate().

Returns:
- The ID of the newly created task on success.
- An appropriate error response if validation or database insertion fails.
*/
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "issue filling error"})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "the issue title is not specified"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error check date"})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := strconv.FormatInt(id, 10)

	writeJSON(w, http.StatusOK, map[string]string{"id": res})
}

/*
GetTaskHandler retrieves a task from the database by its ID.

Expects:
- A query parameter "id" specifying the task ID.

Returns:
- The task details in JSON format if found.
- An error response if the ID is missing or the task is not found.
*/
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idTask := r.URL.Query().Get("id")
	if idTask == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID not specified"})
		return
	}
	task, err := db.GetTask(idTask)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "task not found"})
		return
	}
	writeJSON(w, http.StatusOK, task)
}

/*
UpdateTaskHandler updates an existing task in the database.

Expects:
- A JSON payload with updated task details, including the task ID and title.

Validation steps:
- Ensures the request body is valid JSON.
- Verifies that both the task ID and title are provided.
- Validates and adjusts the task date using checkDate().

Returns:
- A success response if the task is updated.
- An error response if validation fails or the task cannot be updated.
*/
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "issue filling error"})
		return
	}

	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID not specified"})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title not specified"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error check date"})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "task not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}

/*
DeleteTaskHandler deletes a task from the database by its ID.

Expects:
- A query parameter "id" specifying the task ID to delete.

Returns:
- A success response if the task is deleted.
- An error response if the ID is missing or deletion fails.
*/
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID not specified"})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error delete task"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{})
}

/*
TaskDoneHandler marks a task as done. Depending on whether the task has a
repetition rule, it either deletes the task or updates its next occurrence date.

Logic:
  - If the task has no repetition rule (Repeat is empty), the task is deleted.
  - If the task has a repetition rule, the next occurrence date is calculated
    using NextDate() and stored in the database.

Expects:
- A query parameter "id" specifying the task ID.

Returns:
  - A success response if the operation completes.
  - An error response if the ID is missing, the task is not found, or the
    operation fails.
*/
func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID not specified"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "task not found"})
		return
	}

	if task.Repeat == "" {
		if err = db.DeleteTask(task.ID); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error delete task"})
			return
		}
	} else {
		nextDateWork, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error update next date"})
			return
		}
		err = db.UpdateTaskDate(task.ID, nextDateWork)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error update next date"})
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{})
}

/*
checkDate validates and normalizes the date field of a task.

Behavior:
  - If the task date is not provided, sets it to today's date.
  - Parses the date string according to FormatDate.
  - If a repetition rule is defined, calculates the next occurrence date.
  - If the specified date is in the past and there is no repetition rule,
    sets the date to today; if there is a repetition rule, uses the calculated
    next occurrence date instead.

Returns:
  - nil on success.
  - An error describing the validation issue if the date format or repetition
    rule is invalid.
*/
func checkDate(task *db.Task) error {
	nowTime := time.Now()
	now := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)

	if task.Date == "" {
		task.Date = now.Format(FormatDate)
	}

	t, err := time.Parse(FormatDate, task.Date)
	if err != nil {
		return fmt.Errorf("incorrect date format: %v", err)
	}

	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("incorrect repetition rule: %v", err)
		}
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(FormatDate)
		} else {
			task.Date = next
		}
	}

	return nil
}

/*
stubHandler returns a "Not Implemented" response for unfinished or
placeholder endpoints.

Sets:
- Content-Type header to application/json.
- HTTP status code to 501 Not Implemented.
- JSON body with an error message indicating the handler is not implemented.
*/
func stubHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"error": "not implemented"})
}

/*
writeJSON writes a JSON response with the specified HTTP status code.

Steps:
- Sets the Content-Type header to application/json with UTF-8 charset.
- Writes the given status code to the response.
- Encodes the provided data as JSON and writes it to the response writer.

If encoding fails, writes an internal server error response with the encoding
error message.
*/
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

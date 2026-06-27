package api

import (
	"Graduate-work/pkg/db"
	"net/http"
)

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

/*
tasksHandler retrieves a list of tasks from the database with optional search filtering.

Behavior:
  - Reads the "search" query parameter from the request URL to apply a search filter.
  - Calls db.Tasks to fetch up to 50 tasks matching the search criteria.
  - If the database call fails, returns a JSON error response with status 400.
  - Ensures that a nil result from the database is normalized to an empty slice
    to avoid sending null in the JSON response.
  - Returns a JSON response containing the list of tasks wrapped in a TaskResp struct.

Parameters:
- search: optional search string used to filter tasks (passed via query param).

Returns:
- A JSON object with a "Tasks" field containing the list of retrieved tasks on success.
- An error response if the database operation fails.
*/
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error receiving tasks"})
		return
	}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, http.StatusOK, TaskResp{
		Tasks: tasks,
	})
}

// Package main demonstrates using gomb with htmx for server-driven, dynamic UIs.
//
// htmx attributes (hx-get, hx-post, hx-target, hx-swap, …) are regular HTML
// attributes set with .A(). No client-side framework is required beyond the
// htmx CDN script. The server returns HTML fragments for partial page updates.
//
// Run:
//
//	go run main.go
//	open http://localhost:8080
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	. "github.com/ernlel/gomb"
)

// ── data ─────────────────────────────────────────────────────────────────────

type Task struct {
	ID   int
	Text string
	Done bool
}

var tasks = []Task{
	{1, "Learn gomb", true},
	{2, "Integrate htmx", false},
	{3, "Ship to production", false},
}

var nextID = 4

// ── components ───────────────────────────────────────────────────────────────

// taskRow renders a single <li> for the task list.
func taskRow(t Task) *Element {
	var checkbox *Element
	if t.Done {
		checkbox = E("input").
			A("type", "checkbox").
			A("checked", "").
			A("hx-patch", fmt.Sprintf("/tasks/%d/toggle", t.ID)).
			A("hx-target", "#task-list").
			A("hx-swap", "outerHTML")
	} else {
		checkbox = E("input").
			A("type", "checkbox").
			A("hx-patch", fmt.Sprintf("/tasks/%d/toggle", t.ID)).
			A("hx-target", "#task-list").
			A("hx-swap", "outerHTML")
	}

	label := IfElse(
		t.Done,
		E("span").A("class", "line-through text-gray-400").T(t.Text),
		E("span").T(t.Text),
	)

	deleteBtn := E("button").
		A("class", "ml-auto text-red-500 hover:text-red-700 text-sm").
		A("hx-delete", fmt.Sprintf("/tasks/%d", t.ID)).
		A("hx-target", "#task-list").
		A("hx-swap", "outerHTML").
		T("✕")

	return E("li").
		A("class", "flex items-center gap-3 p-2 border-b last:border-0").
		C(checkbox, label, deleteBtn)
}

// taskList renders the full <ul id="task-list"> fragment.
func taskList(ts []Task) *Element {
	return E("ul").
		A("id", "task-list").
		A("class", "divide-y rounded border").
		C(Map(ts, taskRow)...)
}

// addForm renders the task-creation form.
func addForm() *Element {
	return E("form").
		A("hx-post", "/tasks").
		A("hx-target", "#task-list").
		A("hx-swap", "outerHTML").
		A("hx-on::after-request", "this.reset()").
		A("class", "flex gap-2 mt-4").
		C(
			E("input").
				A("type", "text").
				A("name", "text").
				A("placeholder", "New task…").
				A("required", "").
				A("class", "flex-1 border rounded px-3 py-2 text-sm"),
			E("button").
				A("type", "submit").
				A("class", "bg-blue-600 text-white px-4 py-2 rounded text-sm hover:bg-blue-700").
				T("Add"),
		)
}

// page renders the full HTML document.
func page() *Element {
	return E("html").A("lang", "en").C(
		E("head").C(
			E("meta").A("charset", "UTF-8"),
			E("meta").A("name", "viewport").A("content", "width=device-width, initial-scale=1"),
			E("title").T("gomb + htmx – Task List"),
			E("link").
				A("rel", "stylesheet").
				A("href", "https://cdn.jsdelivr.net/npm/tailwindcss@3/dist/tailwind.min.css"),
			// In production pin the version and add integrity+crossorigin attributes.
			// Generate the hash at https://www.srihash.org/
			E("script").
				A("src", "https://unpkg.com/htmx.org@1.9.12").
				A("crossorigin", "anonymous").
				A("defer", ""),
		),
		E("body").A("class", "bg-gray-50 min-h-screen p-8").C(
			E("div").A("class", "max-w-md mx-auto bg-white rounded shadow p-6").C(
				E("h1").A("class", "text-2xl font-bold mb-4").T("Task List"),
				taskList(tasks),
				addForm(),
				// Progress indicator shown while htmx is loading
				E("div").
					A("class", "htmx-indicator mt-2 text-sm text-gray-500").
					T("Saving…"),
			),
		),
	)
}

// ── handlers ─────────────────────────────────────────────────────────────────

func main() {
	mux := http.NewServeMux()

	// Full page
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		page().Render(w)
	})

	// Return the updated task list fragment after adding a task
	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		text := r.FormValue("text")
		if text == "" {
			http.Error(w, "text required", http.StatusBadRequest)
			return
		}
		tasks = append(tasks, Task{ID: nextID, Text: text})
		nextID++

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		taskList(tasks).Render(w)
	})

	// Toggle a task done/undone; return updated list fragment
	mux.HandleFunc("PATCH /tasks/{id}/toggle", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Done = !tasks[i].Done
				break
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		taskList(tasks).Render(w)
	})

	// Delete a task; return updated list fragment
	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		for i, t := range tasks {
			if t.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				break
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		taskList(tasks).Render(w)
	})

	// JSON API endpoint – demonstrates that the same data can be served as
	// both HTML (for htmx) and JSON (for other consumers).
	mux.HandleFunc("GET /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	})

	fmt.Println("htmx example running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}

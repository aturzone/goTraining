package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Task represents a to-do item
type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Status   bool   `json:"status"`
	Priority int    `json:"priority"`
	Deadline string `json:"deadline"`
}

// TaskManager handles all task operations
type TaskManager struct {
	tasks  []Task
	nextID int
	mu     sync.RWMutex
}

var taskManager = &TaskManager{nextID: 1}

func main() {
	// Load tasks from file on startup
	taskManager.loadFromFile()

	// Set up routes
	mux := http.NewServeMux()

	// Task routes
	mux.HandleFunc("/tasks", tasksHandler)         // GET: list all, POST: create new
	mux.HandleFunc("/tasks/", taskHandler)         // GET, PUT, DELETE specific task
	mux.HandleFunc("/tasks/search", searchHandler) // GET: search tasks

	// Add middleware
	handler := loggingMiddleware(corsMiddleware(mux))

	fmt.Println("🚀 To-Do API Server running on http://localhost:8080")
	fmt.Println("📋 Endpoints:")
	fmt.Println("  GET    /tasks           - List all tasks")
	fmt.Println("  POST   /tasks           - Create new task")
	fmt.Println("  GET    /tasks/{id}      - Get specific task")
	fmt.Println("  PUT    /tasks/{id}      - Update task")
	fmt.Println("  DELETE /tasks/{id}      - Delete task")
	fmt.Println("  PUT    /tasks/{id}/done - Mark task as done")
	fmt.Println("  GET    /tasks/search?q=query - Search tasks")

	log.Fatal(http.ListenAndServe(":8080", handler))
}

// Handle /tasks endpoint
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getAllTasks(w, r)
	case "POST":
		createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle /tasks/{id} endpoint
func taskHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")

	// Handle special endpoints
	if strings.HasSuffix(path, "/done") {
		markTaskDone(w, r, path)
		return
	}

	switch r.Method {
	case "GET":
		getTask(w, r, path)
	case "PUT":
		updateTask(w, r, path)
	case "DELETE":
		deleteTask(w, r, path)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get all tasks
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	taskManager.mu.RLock()
	defer taskManager.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": taskManager.tasks,
		"count": len(taskManager.tasks),
	})
}

// Create new task
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask struct {
		Title    string `json:"title" binding:"required"`
		Priority int    `json:"priority"`
		Deadline string `json:"deadline"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if newTask.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	taskManager.mu.Lock()
	task := Task{
		ID:       taskManager.nextID,
		Title:    newTask.Title,
		Priority: newTask.Priority,
		Deadline: newTask.Deadline,
		Status:   false,
	}
	taskManager.tasks = append(taskManager.tasks, task)
	taskManager.nextID++
	taskManager.saveToFile()
	taskManager.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Task created successfully",
		"task":    task,
	})
}

// Get specific task
func getTask(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskManager.mu.RLock()
	defer taskManager.mu.RUnlock()

	for _, task := range taskManager.tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

// Update task
func updateTask(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updateData struct {
		Title    string `json:"title"`
		Priority int    `json:"priority"`
		Deadline string `json:"deadline"`
		Status   bool   `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	taskManager.mu.Lock()
	defer taskManager.mu.Unlock()

	for i, task := range taskManager.tasks {
		if task.ID == id {
			if updateData.Title != "" {
				taskManager.tasks[i].Title = updateData.Title
			}
			if updateData.Priority != 0 {
				taskManager.tasks[i].Priority = updateData.Priority
			}
			if updateData.Deadline != "" {
				taskManager.tasks[i].Deadline = updateData.Deadline
			}
			taskManager.tasks[i].Status = updateData.Status

			taskManager.saveToFile()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Task updated successfully",
				"task":    taskManager.tasks[i],
			})
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

// Delete task
func deleteTask(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskManager.mu.Lock()
	defer taskManager.mu.Unlock()

	for i, task := range taskManager.tasks {
		if task.ID == id {
			taskManager.tasks = append(taskManager.tasks[:i], taskManager.tasks[i+1:]...)
			taskManager.saveToFile()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Task deleted successfully",
				"task":    task,
			})
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

// Mark task as done
func markTaskDone(w http.ResponseWriter, r *http.Request, path string) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimSuffix(path, "/done")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskManager.mu.Lock()
	defer taskManager.mu.Unlock()

	for i, task := range taskManager.tasks {
		if task.ID == id {
			taskManager.tasks[i].Status = true
			taskManager.saveToFile()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Task marked as done",
				"task":    taskManager.tasks[i],
			})
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

// Search tasks
func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	taskManager.mu.RLock()
	defer taskManager.mu.RUnlock()

	var results []Task
	for _, task := range taskManager.tasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(query)) {
			results = append(results, task)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"query":   query,
		"results": results,
		"count":   len(results),
	})
}

// Save tasks to file
func (tm *TaskManager) saveToFile() error {
	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.json", data, 0644)
}

// Load tasks from file
func (tm *TaskManager) loadFromFile() error {
	data, err := os.ReadFile("tasks.json")
	if err != nil {
		// If file doesn't exist, start with empty tasks
		tm.tasks = []Task{}
		return nil
	}

	if err := json.Unmarshal(data, &tm.tasks); err != nil {
		return err
	}

	// Find next ID
	maxID := 0
	for _, task := range tm.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	tm.nextID = maxID + 1

	return nil
}

// Logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

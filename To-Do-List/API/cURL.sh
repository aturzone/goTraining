# 1. Get all tasks
curl http://localhost:8080/tasks

# 2. Create a new task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go HTTP package",
    "priority": 1,
    "deadline": "2024-12-31"
  }'

# 3. Create another task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Build REST API",
    "priority": 2,
    "deadline": "2024-12-25"
  }'

# 4. Get specific task (replace 1 with actual task ID)
curl http://localhost:8080/tasks/1

# 5. Mark task as done
curl -X PUT http://localhost:8080/tasks/1/done

# 6. Update a task
curl -X PUT http://localhost:8080/tasks/2 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Build Awesome REST API",
    "priority": 1,
    "deadline": "2024-12-20",
    "status": false
  }'

# 7. Search tasks
curl "http://localhost:8080/tasks/search?q=API"

# 8. Delete a task
curl -X DELETE http://localhost:8080/tasks/2

# 9. Get all tasks again to see changes
curl http://localhost:8080/tasks
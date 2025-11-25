# Task Management API Documentation

Base URL: http://localhost:8080

## Task Object Format

{
"id": 1,
"title": "Buy groceries",
"description": "Milk, eggs, bread",
"due_date": "2025-11-25",
"status": "pending"
}

---

GET /tasks
Description: Get all tasks.

Response 200 OK:
{
"data": [
{
"id": 1,
"title": "Buy groceries",
"description": "Milk, eggs, bread",
"due_date": "2025-11-25",
"status": "pending"
}
]
}

---

GET /tasks/:id
Description: Get a specific task by ID.

Example: GET /tasks/1

Response 200 OK:
{
"data": {
"id": 1,
"title": "Buy groceries",
"description": "Milk, eggs, bread",
"due_date": "2025-11-25",
"status": "pending"
}
}

Response 400 Bad Request:
{ "error": "invalid task ID" }

Response 404 Not Found:
{ "error": "task not found" }

---

POST /tasks
Description: Create a new task.

Request Body:
{
"title": "Buy groceries",
"description": "Milk, eggs, bread",
"due_date": "2025-11-25",
"status": "pending"
}

Response 201 Created:
{
"data": {
"id": 1,
"title": "Buy groceries",
"description": "Milk, eggs, bread",
"due_date": "2025-11-25",
"status": "pending"
}
}

Response 400 Bad Request:
{ "error": "invalid input: <details>" }

---

PUT /tasks/:id
Description: Update an existing task.

Example: PUT /tasks/1

Request Body:
{
"title": "Buy groceries and snacks",
"description": "Milk, eggs, bread, chips",
"due_date": "2025-11-26",
"status": "in-progress"
}

Response 200 OK:
{
"data": {
"id": 1,
"title": "Buy groceries and snacks",
"description": "Milk, eggs, bread, chips",
"due_date": "2025-11-26",
"status": "in-progress"
}
}

Response 400 Bad Request:
{ "error": "invalid task ID" }

Response 400 Bad Request:
{ "error": "invalid input: <details>" }

Response 404 Not Found:
{ "error": "task not found" }

---

DELETE /tasks/:id
Description: Delete a task by ID.

Example: DELETE /tasks/1

Response 204 No Content:
(no response body)

Response 400 Bad Request:
{ "error": "invalid task ID" }

Response 404 Not Found:
{ "error": "task not found" }

---

Summary:
GET /tasks - Get all tasks
GET /tasks/:id - Get a specific task
POST /tasks - Create a new task
PUT /tasks/:id - Update an existing task
DELETE /tasks/:id - Delete a task

# Task Management API Documentation (JWT + Role-Based Access + Clean Architecture)

Base URL: http://localhost:8080

CLEAN ARCHITECTURE NOTE
This project is refactored using Clean Architecture (Domain, Usecases, Repositories, Infrastructure, Delivery). The API endpoints and behavior remain the same.

Project Structure:
task_manager/
├── Delivery/
│ ├── main.go
│ ├── controllers/
│ │ └── controller.go
│ └── routers/
│ └── router.go
├── Domain/
│ └── domain.go
├── Infrastructure/
│ ├── auth_middleWare.go
│ ├── jwt_service.go
│ └── password_service.go
├── Repositories/
│ ├── task_repository.go
│ └── user_repository.go
└── Usecases/
├── task_usecases.go
└── user_usecases.go

ENVIRONMENT VARIABLES (Recommended)
Create a .env file in the project root:

MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=task_manager_db
JWT_SECRET=super_strong_secret_key

MongoDB Collections:

-   users
-   tasks

RUNNING THE SERVER
From the project root:

1. go mod tidy
2. go run ./Delivery/main.go

AUTHENTICATION

-   This API uses JWT (JSON Web Tokens).
-   Steps:
    1. POST /register
    2. POST /login -> get access_token
    3. Send header for protected endpoints:
       Authorization: Bearer <access_token>

ROLES & ACCESS RULES

-   If the database is empty, the first registered user becomes admin
-   Admins can promote other users to admin using POST /promote
-   Only admins can create, update, and delete tasks
-   Any authenticated user (admin or regular user) can get all tasks and get task by ID

USER OBJECT FORMAT
{
"id": "6926e9266694724e5e274283",
"username": "admin1",
"role": "admin"
}

TASK OBJECT FORMAT
{
"id": "6926e9266694724e5e274283",
"title": "Buy groceries",
"description": "Milk, eggs, bread",
"due_date": "2025-11-25",
"status": "pending"
}

---

POST /register (Public)
Description: Create a new user account (unique username). If the database has no users, the first registered user will be an admin.

Request Body:
{
"username": "admin1",
"password": "pass123"
}

Response 201 Created:
{
"data": {
"id": "6926e9266694724e5e274283",
"username": "admin1",
"role": "admin"
}
}

Response 400 Bad Request:
{ "error": "invalid input: <details>" }

Response 409 Conflict:
{ "error": "username already exists" }

---

POST /login (Public)
Description: Authenticate user and return a JWT access token.

Request Body:
{
"username": "admin1",
"password": "pass123"
}

Response 200 OK:
{
"data": {
"access_token": "<jwt_token_here>",
"user": {
"id": "6926e9266694724e5e274283",
"username": "admin1",
"role": "admin"
}
}
}

Response 400 Bad Request:
{ "error": "invalid input: <details>" }

Response 401 Unauthorized:
{ "error": "invalid username or password" }

---

POST /promote (Admin Only)
Description: Promote an existing user to admin.

Headers:
Authorization: Bearer <access_token>

Request Body:
{
"username": "user1"
}

Response 200 OK:
{
"data": {
"username": "user1",
"role": "admin"
}
}

Response 400 Bad Request:
{ "error": "invalid input: <details>" }

Response 401 Unauthorized:
{ "error": "missing Authorization header" }

Response 403 Forbidden:
{ "error": "admin access required" }

Response 404 Not Found:
{ "error": "user not found" }

---

GET /tasks (Authenticated Users)
Description: Get all tasks.

Headers:
Authorization: Bearer <access_token>

Response 200 OK:
{
"data": [
{
"id": "6926e9266694724e5e274283",
"title": "Buy groceries",
"description": "Milk, eggs, bread",
"due_date": "2025-11-25",
"status": "pending"
}
]
}

Response 401 Unauthorized:
{ "error": "missing Authorization header" }

---

GET /tasks/:id (Authenticated Users)
Description: Get a specific task by ID.

Headers:
Authorization: Bearer <access_token>

Example: GET /tasks/6926e9266694724e5e274283

Response 200 OK:
{
"data": {
"id": "6926e9266694724e5e274283",
"title": "Buy groceries",
"description": "Milk, eggs, bread",
"due_date": "2025-11-25",
"status": "pending"
}
}

Response 401 Unauthorized:
{ "error": "missing Authorization header" }

Response 404 Not Found:
{ "error": "task not found" }

---

POST /tasks (Admin Only)
Description: Create a new task.

Headers:
Authorization: Bearer <access_token>

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
"id": "6926e9266694724e5e274283",
"title": "Buy groceries",
"description": "Milk, eggs, bread",
"due_date": "2025-11-25",
"status": "pending"
}
}

Response 400 Bad Request:
{ "error": "invalid input: <details>" }

Response 401 Unauthorized:
{ "error": "missing Authorization header" }

Response 403 Forbidden:
{ "error": "admin access required" }

---

PUT /tasks/:id (Admin Only)
Description: Update an existing task by ID.

Headers:
Authorization: Bearer <access_token>

Example: PUT /tasks/6926e9266694724e5e274283

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
"id": "6926e9266694724e5e274283",
"title": "Buy groceries and snacks",
"description": "Milk, eggs, bread, chips",
"due_date": "2025-11-26",
"status": "in-progress"
}
}

Response 400 Bad Request:
{ "error": "invalid input: <details>" }

Response 401 Unauthorized:
{ "error": "missing Authorization header" }

Response 403 Forbidden:
{ "error": "admin access required" }

Response 404 Not Found:
{ "error": "task not found" }

---

DELETE /tasks/:id (Admin Only)
Description: Delete a task by ID.

Headers:
Authorization: Bearer <access_token>

Example: DELETE /tasks/6926e9266694724e5e274283

Response 204 No Content:
(no response body)

Response 401 Unauthorized:
{ "error": "missing Authorization header" }

Response 403 Forbidden:
{ "error": "admin access required" }

Response 404 Not Found:
{ "error": "task not found" }

---

CURL QUICK TESTS

Register:
curl -i -X POST http://localhost:8080/register -H "Content-Type: application/json" -d '{"username":"admin1","password":"pass123"}'

Login:
curl -i -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"admin1","password":"pass123"}'

Get tasks (replace TOKEN):
curl -i http://localhost:8080/tasks -H "Authorization: Bearer TOKEN"

Create task (admin only, replace TOKEN):
curl -i -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -H "Authorization: Bearer TOKEN" -d '{"title":"Task A","description":"demo","due_date":"2025-12-20","status":"pending"}'

---

SUMMARY
Public:

-   POST /register
-   POST /login

Authenticated (JWT required):

-   GET /tasks
-   GET /tasks/:id

Admin only (JWT + role=admin):

-   POST /tasks
-   PUT /tasks/:id
-   DELETE /tasks/:id
-   POST /promote

```

```

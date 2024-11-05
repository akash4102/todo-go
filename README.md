
# TODO App API

A simple TODO application built with Go, using the Chi router and MongoDB for data storage.

---

### Clone the Repository

```bash
git clone https://github.com/akash4102/todo-app
cd todo-go
```

### Install Dependencies

```bash
go mod tidy
```

### Running the Application

1. **Ensure MongoDB is running** on your machine or through a cloud provider.
2. Start the application:

   ```bash
   go run main.go
   ```

The server will start on the specified port (e.g., `localhost:8080`).

---

## Configuration

Configure the database and other settings in the `.env` file:

```plaintext
PORT=:3333
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=todo_app
MONGODB_COLLECTION=todos
```

- **DB_URI**: MongoDB connection URI.
- **DB_NAME**: MongoDB database name.
- **PORT**: Port for the application to run.

---

## API Endpoints

All routes are prefixed with `/todos` and managed by the `TodoController`.

### GET /todos

Retrieve a list of all TODO items.

- **Response**: List of TODO items with details.

### GET /todos/{id}

Retrieve a specific TODO item by its ID.

- **Parameters**:
  - `id`: The ID of the TODO item.

- **Response**: Details of the specific TODO item.

### POST /todos

Create a new TODO item.

- **Request Body**:
  ```json
  {
    "title": "Sample Todo",
    "description": "Sample description for the todo item",
    "status": "pending"
  }
  ```

- **Response**: Details of the newly created TODO item.

### PUT /todos/{id}

Update a specific TODO item.

- **Parameters**:
  - `id`: The ID of the TODO item.

- **Request Body**:
  ```json
  {
    "title": "Updated Todo Title",
    "description": "Updated description",
    "status": "completed"
  }
  ```

- **Response**: Details of the updated TODO item.

### DELETE /todos/{id}

Delete a specific TODO item by ID.

- **Parameters**:
  - `id`: The ID of the TODO item.

- **Response**: Confirmation of deletion.

---

## Technologies Used

- **Go**: Primary language for server-side logic.
- **Chi**: Lightweight router for handling HTTP requests.
- **MongoDB**: Database for storing TODO data.

---


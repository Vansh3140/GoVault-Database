# GOVault: A Lightweight Go-Based JSON Database

GOVault is a minimalist yet powerful Go-based database system that uses JSON files for data persistence. It provides a straightforward interface for managing collections and records, making it ideal for small-scale projects or applications requiring lightweight data storage solutions. 

---

## Features 🚀
- **RESTful API Interface**: Perform CRUD operations via HTTP.
- **JSON Storage**: Data is stored in JSON format for simplicity and readability.
- **Thread-Safe Operations**: Ensures data integrity using collection-specific mutexes.
- **Modular Design**: Separation of concerns through `drivers`, `routes`, and `main` modules.
- **Easy to Set Up**: No external dependencies, just Go and your code.

---

## Project Structure 📁

```
GOVault/
├── main.go             # Application entry point
├── drivers/            # Database management logic
│   └── drivers.go
├── routes/             # HTTP route definitions and handlers
│   └── routes.go
```

### Module Responsibilities
1. **`main.go`**: Initializes the app, sets up routes, and manages graceful shutdowns.
2. **`drivers/`**: Handles low-level database operations like reading, writing, and deleting data.
3. **`routes/`**: Defines RESTful API endpoints and links them to the database functions.

---

## Installation and Setup 🛠️

### Prerequisites
- **Go 1.18+** installed on your machine.

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/Vansh3140/GOVault.git
   ```
2. Navigate to the project directory:
   ```bash
   cd GOVault
   ```
3. Run the application:
   ```bash
   go run main.go
   ```
4. The application will start at `http://localhost:8080`.

---

## API Endpoints 🌐

### Base URL
```
http://localhost:8080/api/govault
```

### Endpoints Overview
| HTTP Method | Endpoint                   | Description                          |
|-------------|----------------------------|--------------------------------------|
| GET         | `/:collection`            | Fetch all records from a collection |
| GET         | `/:collection/:resource`  | Fetch a specific record             |
| POST        | `/:collection`            | Create a new record                 |
| PUT         | `/:collection/:resource`  | Update a specific record            |
| DELETE      | `/:collection`            | Delete all records in a collection  |
| DELETE      | `/:collection/:resource`  | Delete a specific record            |

---

## Example API Usage 📋

### 1. Fetch All Records in a Collection
**Endpoint:**  
`GET /api/govault/users`

**cURL Command:**
```bash
curl -X GET http://localhost:8080/api/govault/users
```

**Response:**
```json
[
  {
    "Name": "John Doe",
    "Age": "30",
    "Contact": "1234567890"
  },
  {
    "Name": "Jane Smith",
    "Age": "28",
    "Contact": "9876543210"
  }
]
```

---

### 2. Fetch a Specific Record
**Endpoint:**  
`GET /api/govault/users/:resource`

**cURL Command:**
```bash
curl -X GET http://localhost:8080/api/govault/users/John_Doe
```

**Response:**
```json
{
  "Name": "John Doe",
  "Age": "30",
  "Contact": "1234567890"
}
```

---

### 3. Add a New Record
**Endpoint:**  
`POST /api/govault/:collection`

**cURL Command:**
```bash
curl -X POST http://localhost:8080/api/govault/users -H "Content-Type: application/json" -d '{
  "Name": "Alice Johnson",
  "Age": "35",
  "Contact": "5678901234"
}'
```

**Response:**
```json
{
  "message": "Record added successfully."
}
```

---

### 4. Update a Record
**Endpoint:**  
`PUT /api/govault/:collection/:resource`

**cURL Command:**
```bash
curl -X PUT http://localhost:8080/api/govault/users/Alice_Johnson -H "Content-Type: application/json" -d '{
  "Name": "Alice Johnson",
  "Age": "36",
  "Contact": "5678901234"
}'
```

**Response:**
```json
{
  "message": "Record updated successfully."
}
```

---

### 5. Delete All Records in a Collection
**Endpoint:**  
`DELETE /api/govault/:collection`

**cURL Command:**
```bash
curl -X DELETE http://localhost:8080/api/govault/users
```

**Response:**
```json
{
  "message": "All records in collection 'users' deleted successfully."
}
```

---

### 6. Delete a Specific Record
**Endpoint:**  
`DELETE /api/govault/:collection/:resource`

**cURL Command:**
```bash
curl -X DELETE http://localhost:8080/api/govault/users/Alice_Johnson
```

**Response:**
```json
{
  "message": "Record 'Alice_Johnson' deleted successfully."
}
```

---

## Licensing 📜

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Author 👨‍💻
**Vansh3140**  
Feel free to reach out or open an issue if you encounter any problems or have suggestions.

Enjoy building with GOVault! 🎉

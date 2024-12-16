Here’s a detailed **`README.md`** file to help users set up and run the application locally:

---

# **Event Trigger Platform**

This is a Go-based platform for managing event triggers and logs, featuring support for scheduled and API-based triggers. It uses MySQL as the database and Swagger UI for API documentation.

---

## **Features**
- Create, list, and delete triggers.
- Handle scheduled and recurring triggers.
- Log trigger executions.
- Interactive API documentation via Swagger UI.

---

## **Prerequisites**
1. Install **Go** (v1.17 or higher): [https://golang.org/dl/](https://golang.org/dl/)
2. Install **MySQL** (v5.7 or higher): [https://dev.mysql.com/downloads/](https://dev.mysql.com/downloads/)
3. Install **Git**: [https://git-scm.com/](https://git-scm.com/)
4. Install **Swag CLI** for generating Swagger documentation:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

---

## **Setup Instructions**

### **1. Clone the Repository**
```bash
git clone https://github.com/your-username/event-trigger-platform.git
cd event-trigger-platform
```

### **2. Configure the Database**
1. **Start MySQL** and create the required database:
   ```sql
   CREATE DATABASE eventdb;
   ```
2. **Create the tables**:
   ```sql
   USE eventdb;

   CREATE TABLE IF NOT EXISTS triggers (
       id VARCHAR(255) PRIMARY KEY,
       type VARCHAR(50) NOT NULL,
       payload JSON,
       scheduled_at DATETIME,
       recurring BOOLEAN,
       interval INT,
       test BOOLEAN
   );

   CREATE TABLE IF NOT EXISTS event_logs (
       id INT AUTO_INCREMENT PRIMARY KEY,
       trigger_id VARCHAR(255),
       timestamp DATETIME NOT NULL,
       message TEXT,
       FOREIGN KEY (trigger_id) REFERENCES triggers(id)
   );
   ```

3. **Update `db.go`** with your MySQL credentials:
   - Open `db/db.go` and edit the `dsn` (Data Source Name):
     ```go
     dsn := "username:password@tcp(localhost:3306)/eventdb"
     ```

---

### **3. Install Dependencies**
Run the following command to download the necessary Go modules:
```bash
go mod tidy
```

---

### **4. Generate Swagger Documentation**
Generate the Swagger documentation using the Swag CLI:
```bash
swag init
```

This will create the `docs` folder with the necessary Swagger files.

---

### **5. Start the Server**
Run the application:
```bash
go run main.go
```

The server will start at `http://localhost:8080`.

---

## **Usage**

### **Swagger UI**
Access the API documentation and test the endpoints using Swagger UI:
```
http://localhost:8080/swagger/index.html
```

### **API Endpoints**
- **Create Trigger**:  
  **POST** `/triggers`  
  Request Body (JSON):
  ```json
  {
      "type": "Scheduled",
      "payload": {"key1": "value1"},
      "scheduled_at": "2024-12-20T15:00:00Z",
      "recurring": true,
      "interval": 3600,
      "test": false
  }
  ```

- **List Triggers**:  
  **GET** `/triggers`

- **Delete Trigger**:  
  **DELETE** `/triggers/delete?id={trigger_id}`

- **List Logs**:  
  **GET** `/logs`

---

## **Testing**
- Use **Postman** or Swagger UI to interact with the APIs.
- Verify data in the database using MySQL queries:
  ```sql
  SELECT * FROM triggers;
  SELECT * FROM event_logs;
  ```

---

## **Stopping the App**
Press `Ctrl + C` in the terminal to stop the application.

---

## **Troubleshooting**

### **1. MySQL Connection Issues**
- Ensure MySQL is running and accessible.
- Verify the credentials in `db/db.go`.

### **2. Swagger Not Loading**
- Ensure `swag init` was executed successfully.
- Verify the Swagger JSON file at `http://localhost:8080/swagger/swagger.json`.

### **3. Missing Dependencies**
- Run `go mod tidy` to resolve dependency issues.

---

## **Contributing**
Contributions are welcome! Please fork the repository, create a new branch, and submit a pull request.

---

## **License**
This project is licensed under the MIT License. See the `LICENSE` file for details.

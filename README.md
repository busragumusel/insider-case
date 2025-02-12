# 📌 Insider Case - Message Processing API

## 🚀 Project Setup Guide
This is a Go-based backend API that processes messages, integrates with PostgreSQL and Redis, and provides an HTTP interface using Chi router.

---

## **📌 Prerequisites**
Make sure you have the following installed:

- **Go** (>=1.23)
- **Docker & Docker Compose**
- **PostgreSQL**
- **Redis**

---

## **📌 Running the Project**

### **🐳 Using Docker (Recommended)**
```sh
docker-compose build --no-cache            
docker-compose up
```
This will start the app, PostgreSQL, and Redis inside Docker containers.

### **🖥️ Running Locally**
1. **Start PostgreSQL**
   ```sh
   brew services start postgresql  # macOS
   sudo systemctl start postgresql  # Linux
   ```
   Ensure PostgreSQL is running before proceeding.

2. **Start Redis**
   ```sh
   brew services start redis  # macOS
   sudo systemctl start redis  # Linux
   ```

3. **Run the Application**
   ```sh
   go run main.go
   ```

---

## **📌 API Documentation**
The project includes **Swagger API documentation**.
Once the app is running, open:
```
http://localhost:8080/swagger/index.html
```

---

## **📌 Running Tests**
Unit tests can be executed using:
```sh
go test -v ./...
```

To test using Docker:
```sh
docker-compose exec app go test -v ./...
```

---

## **📌 Endpoints**

### **🔹 Start Message Processing**
```http
GET /start
```
**Response:**
```json
{ "message": "Message processing started" }
```

### **🔹 Stop Message Processing**
```http
GET /stop
```
**Response:**
```json
{ "message": "Message processing stopped" }
```

### **🔹 Retrieve Messages**
```http
GET /messages?status=sent
```
**Response:**
```json
{
  "messages": [
    {
      "id": 1,
      "phone_number": "+905551111111",
      "content": "Test Message",
      "status": "pending"
    }
  ]
}
```

---

## **📌 Useful Commands**
### **Check Running Containers**
```sh
docker ps
```

### **Check Docker Logs**
```sh
docker logs insider-case
```

### **Restart Docker Services**
```sh
docker-compose restart
```

---

## **📌 Troubleshooting**

### **❌ PostgreSQL: Connection Refused**
- Make sure PostgreSQL is running: `brew services list`
- Restart PostgreSQL: `brew services restart postgresql`

### **❌ Redis: Connection Failed**
- If running locally, update `.env`: `REDIS_HOST=127.0.0.1`
- Restart Redis: `brew services restart redis`

### **❌ Swagger Not Working?**
- Regenerate Swagger docs: `swag init`

---

## **📌 Contributors**
- **[Your Name]** - Developer

---

## **📌 License**
This project is licensed under **UnLicense**. Feel free to use and modify it!


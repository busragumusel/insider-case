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

## **📌 Running SQL Inside a PostgreSQL Docker Container To Test**

If your PostgreSQL container is already running, follow these steps to execute SQL commands inside it.

### 1️⃣ Find Your PostgreSQL Container ID or Name

Run the following command to list running containers and find the one running PostgreSQL:

```sh
docker ps
```  

Look for a container with the `postgres` image and note its **container ID or name**.

### 2️⃣ Access the PostgreSQL Container's Shell

Use the following command to enter the PostgreSQL shell inside the running container:

```sh
docker exec -it <container_id_or_name> psql -U postgres -d insider_case_test
```  

Replace `<container_id_or_name>` with your actual container name or ID.  
If your database name is different, replace `insider_case_test` with the correct database name.

### 3️⃣ Run the SQL Command Inside `psql`

Once inside `psql`, execute the following SQL command to insert records into the `messages` table:

```sql
INSERT INTO "messages" ("id", "phone_number", "content", "status", "sent_at", "created_at") VALUES
(1, '+905714822412', 'Quis nostrud exercitation ullamco', 'sent', '2025-09-28 02:52:22.000000', '2025-05-02 07:44:29.000000'),
(2, '+905732050897', 'Consectetur adipiscing elit', 'sent', '2025-07-21 17:43:52.000000', '2025-04-19 03:59:46.000000'),
(3, '+905124620539', 'Sunt in culpa qui officia deserunt mollit anim id est laborum', 'sent', '2025-07-05 20:03:29.000000', '2025-08-28 16:28:42.000000'),
(4, '+905643631747', 'Excepteur sint occaecat cupidatat non proident', 'sent', '2025-01-14 19:56:13.000000', '2025-02-14 11:32:56.000000'),
(5, '+905443596825', 'Excepteur sint occaecat cupidatat non proident', 'sent', '2025-01-28 10:03:29.000000', '2025-09-11 06:44:02.000000'),
(6, '+905648658345', 'Duis aute irure dolor in reprehenderit', 'pending', '2025-05-26 13:06:04.000000', '2025-05-28 17:39:53.000000'),
(7, '+905653416977', 'Quis nostrud exercitation ullamco', 'sent', '2025-12-05 17:50:01.000000', '2025-12-23 21:22:58.000000');
```  

### 4️⃣ Exit `psql`

After running the SQL command, exit the PostgreSQL shell by typing:

```sh
\q
```  

Now, you have successfully inserted the data into the PostgreSQL database running inside the Docker container! 🚀

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

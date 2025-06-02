#  Skill Barter Platform (Backend)

This is the backend API for a **Skill Barter Platform for Students**, where users can offer and request skills, match with others, and exchange learning sessions. Built using **Go (Gin)** and **MongoDB**.

📽️ **Project Walkthrough:** 
I have included this video to shows live API interactions 
[Skill Barter Walkthrough (Google Drive)](https://drive.google.com/file/d/1ghLobAz04YG58X3ccJamHeiw0FZCTkl2/view?usp=drive_link)

---

##  Tech Stack

- **Backend:** Go, Gin
- **Database:** MongoDB
- **Authentication:** JWT, bcrypt
- **Real-time Communication:** WebSocket (Gorilla WebSocket)
- **Middleware:** JWT auth, CORS
- **API Format:** REST

---

##  Setup Instructions

```bash
# 1. Clone the repository
git clone https://github.com/your-username/skill-barter-backend.git
cd skill-barter-backend

# 2. Install dependencies
go mod tidy

# 3. Start MongoDB if not already running

# 4. Run the server
go run main.go
```

Server runs at: `http://localhost:8000`

---

##  Authentication Routes

### POST `/auth/signup`

Create a new user account.

**Request:**
```json
{
  "name": "Rupanjan",
  "email": "rupanjan@example.com",
  "password": "yourpassword"
}
```

---

### POST `/auth/login`

Authenticate user and receive JWT token.

**Request:**
```json
{
  "email": "rupanjan@example.com",
  "password": "yourpassword"
}
```

**Response:**
```json
{
  "token": "JWT-TOKEN"
}
```

---

##  Protected Routes

**Header Required:**  
`Authorization: Bearer <JWT-TOKEN>`

---

### GET `/auth/myprofile`

Fetch the current user's profile.

---

### PUT `/auth/myprofile`

Update user profile including skills and availability.

**Request:**
```json
{
  "skillsHave": ["Go", "React"],
  "skillsWant": ["Python", "Docker"],
  "availableDays": ["Saturday", "Sunday"]
}
```

---

##  Matching Route

### GET `/match`

Returns users whose skillsWant intersect with your skillsHave, and whose skillsHave intersect with your skillsWant.

**Response:**
```json
{
  "matchCount": 2,
  "names": ["Harish", "Karu"]
}
```

---

##  Skill Request Routes

### POST `/request/send`

Send a skill barter request to another user.

**Request:**
```json
{
  "toEmail": "harish@gmail.com",
  "skill": "football"
}
```

---

### POST `/request/respond`

Respond to a skill request with either `accepted` or `rejected`.

**Request:**
```json
{
  "fromName": "Rupanjan",
  "status": "accepted"
}
```

---

##  WebSocket Integration

Real-time communication is enabled using WebSockets.

### WebSocket Endpoint:

```
ws://localhost:8000/ws?email=<your-email>
```

### Example (using wscat):

```bash
wscat -c "ws://localhost:8000/ws?email=harish@gmail.com"
```

### Notifications:

- On request received:
  ```
  Hey! You got a new skill request from Rupanjan for: football
  ```

- On request accepted:
  ```
  Your request was accepted by Harish
  ```

---

##  Folder Structure

```
.
├── config/         # MongoDB config
├── controllers/    # Route logic
├── models/         # DB models (User, SkillRequest)
├── websocket/      # WebSocket handlers
├── main.go         # Entry point
```

---

##  Testing With Postman

1. Signup via `/auth/signup`
2. Login via `/auth/login` and copy JWT
3. Use token with protected routes (add header `Authorization: Bearer <token>`)
4. Test profile, match, and request routes
5. Connect WebSocket and observe real-time messages

---

##  CORS Setup

For local development with a frontend on `localhost:5173`, CORS is enabled:

```go
import "github.com/gin-contrib/cors"

r.Use(cors.Default())
```

---

##  Features

- 🔐 Secure signup/login with JWT & bcrypt
- 🔄 Real-time updates using WebSocket
- 🔎 Skill-based matching algorithm
- 📩 Skill request system
- 🛡️ Fully protected APIs with middleware

---

##  Contributions

Feel free to fork the repository and open a pull request. All contributions are welcome!

---

## 📄 License

This project is licensed under the MIT License.

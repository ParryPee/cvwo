# CVWO Assignment 2026: Gossip with Go - Web Forum

This repository contains the source code for a web forum application built for the CVWO AY2025/26 assignment. The application allows users to view topics, create posts, and engage in discussions via comments.

## Tech Stack

**Frontend:**
* **Framework:** React with TypeScript
* **Styling:** Material UI (MUI) & Tailwind CSS
* **State/Data Fetching:** Axios
* **Routing:** React Router

**Backend:**
* **Language:** Go (Golang)
* **Router:** Gorilla Mux
* **Database Driver:** go-sql-driver/mysql
* **Authentication:** JWT (JSON Web Tokens)

**Database:**
* MySQL

---

## Setup Instructions

Follow these steps to set up the application locally for development.

### 1. Database Setup
Ensure you have MySQL installed and running.

1.  Create a new MySQL database (e.g., `cvwo_db`).
2.  Import the provided SQL dump file to seed the database tables:
    ```bash
    mysql -u <your_username> -p <database_name> < backend_final.sql
    ```

### 2. Backend Setup

1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Create a `.env` file in the `backend/` root directory with the following configuration (adjust values to match your local MySQL setup):
    ```env
    # Database Configuration
    DB_USER=root
    DB_PASS=your_password
    DB_HOST=localhost
    DB_PORT=3306
    DB_NAME=cvwo_db
    
    # Security
    JWT_KEY=your_secret_jwt_key
    
    # CORS (Frontend URL)
    FRONTEND_URL=http://localhost:5173
    ```
3.  Install dependencies:
    ```bash
    go mod tidy
    ```
4.  Start the server:
    ```bash
    go run cmd/main.go
    ```
    The backend should now be running on `http://localhost:8080`.

### 3. Frontend Setup

1.  Navigate to the frontend directory:
    ```bash
    cd frontend-new
    ```
2.  Create a `.env` file in the `frontend-new/` root directory:
    ```env
    VITE_API_URL=http://localhost:8080
    ```
3.  Install dependencies:
    ```bash
    npm install
    ```
4.  Start the development server:
    ```bash
    npm run dev
    ```
5.  Open your browser and navigate to `http://localhost:5173` (or the port shown in your terminal).

---

## Features Implemented

* **User Authentication**: Register, Login, and Logout functionality using JWT.
* **Topics**: Browse existing topics in the community or Create and Update your own. 
* **Posts**: Create, read, update, and delete posts within topics.
* **Comments**: Comment on posts to discuss with other users. Sub-replies are also supported.
* **Likes**: Like posts and comments.
* **Search**: Search for specific posts or topics.
* **Protected Routes**: Certain actions (creating/editing content) are restricted to authorised logged-in users.

---

## AI Usage Declaration

*Per the assignment requirements, the following is a declaration of AI tools used during the development of this project.*

| Tool Used | Purpose |
| :--- | :--- |
| **Gemini** | Used AI to learn about Tailwind utility classes. |
| **Gemini** | Used AI to learn about the differences between the full-text search modes. |


---

## Author
* **Name:** Lee Ho Yang Larry
* **Matriculation Number:** A0315278Y

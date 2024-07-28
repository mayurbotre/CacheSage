# CacheSage

## Overview

CacheSage is a full-stack application designed to manage a Least Recently Used (LRU) cache with real-time updates. It features a Golang backend and a React frontend, providing a seamless user experience for cache management.

## Technologies

- **Backend**: Go (Golang), Gin Framework, Gorilla WebSocket
- **Frontend**: React.js, Axios

## Setup and Running

### 1. Move to the directory

```
cd CacheSage
```

### 2. Setup and Run the Backend Server
Navigate to the Backend Directory

```
cd server
```

Install Dependencies

```
go mod tidy
```

Run the Backend Server

```
go run main.go
```

The backend server will start on http://localhost:8080.

### 3. Setup and Run the Frontend Application

Navigate to the Frontend Directory

```
cd ../client
```

Install Dependencies

```
npm install
```

Start the Development Server

```
npm start
```

The frontend application will start on http://localhost:3000

## Contact

Feel free to reach out in case of any issues.
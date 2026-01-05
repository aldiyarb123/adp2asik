# adp2asik
# Assignment 2 – Concurrent Web Server in Go

**Course:** Advanced Programming 1  
**Assignment:** 2  
**Student:** Aldiyar Bazarbayev
**Group:** SE2416

---

## 📌 Project Overview

This project implements a **concurrent web server** using the Go programming language and the standard `net/http` package.

The server supports multiple REST endpoints, safely handles concurrent requests using mutexes, runs a background worker using goroutines, and shuts down gracefully using context and OS signals.

The project demonstrates understanding of:
- Go concurrency model
- Goroutines and mutexes
- Context-based cancellation
- RESTful API design
- Clean project structure

---

## 🏗️ Project Structure

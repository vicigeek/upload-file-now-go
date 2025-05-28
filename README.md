# Go Carbon File Uploader

A modern and minimal **file uploader written in Go**, styled with **IBM Carbon Design System**, featuring a responsive UI and real-time animated loader from **Uiverse.io**. Perfect for self-hosted upload utilities and minimal dashboards.

---

## ✨ Features

- ✅ Built in **pure Go** (Golang)
- 🎨 Styled with **Carbon Design System**
- ⏳ Includes animated **Uiverse.io loader**
- 📁 Saves uploads to `/tmp/uploads`
- 📋 Detailed terminal logging with IP + timestamps
- 🛑 Graceful shutdown on Ctrl+C
- ⚡ No frontend build tools required (vanilla HTML/CSS/JS embedded)

---

## 🚀 Demo

![image](https://github.com/user-attachments/assets/a887f48f-767a-4c90-b327-0994cc2f8070)


---

## 🧠 Use Cases

- Upload tool for internal servers
- Lightweight admin panels
- Self-hosted utility dashboards
- Building blocks for AI tools with file input

---

## 📦 Requirements

- Go 1.18+
- Internet connection (to load Carbon CDN)

---

## 🛠️ How to Run

```bash
git clone https://github.com/vicigeek/upload-file-now-go.git
cd upload-file-now-go
go run uploader.go

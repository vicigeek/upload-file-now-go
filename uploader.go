package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
    "time"
)

// Embedded HTML with Carbon Design + Smart Loader
const uploadForm = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Upload File</title>
  <link rel="stylesheet" href="https://unpkg.com/carbon-components/css/carbon-components.min.css">
  <script src="https://unpkg.com/carbon-components/scripts/carbon-components.min.js" defer></script>
  <style>
    body {
      background-color: #1c1c1c;
      font-family: 'IBM Plex Sans', sans-serif;
      display: flex;
      justify-content: center;
      align-items: center;
      height: 100vh;
      margin: 0;
      color: #fff;
    }

    .container {
      background: #2c2c2c;
      padding: 2rem;
      border-radius: 12px;
      box-shadow: 0 0 20px rgba(0, 0, 0, 0.4);
      width: 420px;
    }

    .bx--form-item {
      margin-bottom: 1rem;
    }

    .loader {
      --c1: #673b14;
      --c2: #f8b13b;
      width: 40px;
      height: 80px;
      border-top: 4px solid var(--c1);
      border-bottom: 4px solid var(--c1);
      background: linear-gradient(90deg, var(--c1) 2px, var(--c2) 0 5px,var(--c1) 0) 50%/7px 8px no-repeat;
      display: none;
      overflow: hidden;
      animation: l5-0 2s infinite linear;
      margin: 1rem auto;
    }

    .loader::before,
    .loader::after {
      content: "";
      grid-area: 1/1;
      width: 75%;
      height: calc(50% - 4px);
      margin: 0 auto;
      border: 2px solid var(--c1);
      border-top: 0;
      box-sizing: content-box;
      border-radius: 0 0 40% 40%;
      -webkit-mask: linear-gradient(#000 0 0) bottom/4px 2px no-repeat,
        linear-gradient(#000 0 0);
      -webkit-mask-composite: destination-out;
      mask-composite: exclude;
      background: linear-gradient(var(--d,0deg),var(--c2) 50%,#0000 0) bottom /100% 205%,
        linear-gradient(var(--c2) 0 0) center/0 100%;
      background-repeat: no-repeat;
      animation: inherit;
      animation-name: l5-1;
    }

    .loader::after {
      transform-origin: 50% calc(100% + 2px);
      transform: scaleY(-1);
      --s: 3px;
      --d: 180deg;
    }

    @keyframes l5-0 {
      80% { transform: rotate(0) }
      100% { transform: rotate(0.5turn) }
    }

    @keyframes l5-1 {
      10%,70% { background-size: 100% 205%,var(--s,0) 100% }
      70%,100% { background-position: top,center }
    }
  </style>
</head>
<body>
  <div class="container">
    <h4 class="bx--form__heading">Upload a File</h4>
    <form id="uploadForm" enctype="multipart/form-data" method="post">
      <div class="bx--form-item">
        <div data-file class="bx--file">
          <label for="file" class="bx--file-browse-btn bx--btn bx--btn--secondary" role="button">
            Choose file
            <input type="file" name="file" id="file" class="bx--file-input" required>
          </label>
          <div class="bx--file-container"></div>
        </div>
      </div>
      <div class="bx--form-item">
        <button type="submit" class="bx--btn bx--btn--primary">Upload</button>
      </div>
      <div id="spinner" class="loader"></div>
      <p id="resultMessage" style="margin-top: 1rem;"></p>
    </form>
  </div>

  <script>
    const form = document.getElementById('uploadForm');
    const spinner = document.getElementById('spinner');
    const resultMessage = document.getElementById('resultMessage');

    form.addEventListener('submit', function(event) {
      event.preventDefault();
      const fileInput = document.getElementById('file');
      const file = fileInput.files[0];
      if (!file) {
        resultMessage.textContent = "❌ No file selected.";
        return;
      }

      spinner.style.display = 'grid';
      resultMessage.textContent = "";

      const xhr = new XMLHttpRequest();
      xhr.open("POST", "/upload", true);

      xhr.onload = function () {
        spinner.style.display = 'none';
        if (xhr.status === 200) {
          resultMessage.innerHTML = "✅ " + xhr.responseText;
        } else {
          resultMessage.innerHTML = "❌ Upload failed: " + xhr.responseText;
        }
        form.reset();
      };

      xhr.onerror = function () {
        spinner.style.display = 'none';
        resultMessage.innerHTML = "❌ Network error during upload.";
      };

      const formData = new FormData();
      formData.append("file", file);
      xhr.send(formData);
    });
  </script>
</body>
</html>
`

// LogRequest logs detailed HTTP request info
func LogRequest(r *http.Request) {
    clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
    log.Printf("[INFO] %s %s from %s", r.Method, r.URL.Path, clientIP)
}

// uploadHandler handles file uploads
func uploadHandler(w http.ResponseWriter, r *http.Request) {
    LogRequest(r)

    if r.Method != http.MethodPost {
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, uploadForm)
        return
    }

    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        log.Printf("[ERROR] Parse form: %v", err)
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        log.Printf("[ERROR] Get file: %v", err)
        http.Error(w, "Error retrieving file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    uploadDir := "/tmp/uploads"
    os.MkdirAll(uploadDir, os.ModePerm)

    dstPath := filepath.Join(uploadDir, filepath.Base(handler.Filename))
    dst, err := os.Create(dstPath)
    if err != nil {
        log.Printf("[ERROR] Create file: %v", err)
        http.Error(w, "Unable to create file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    _, err = io.Copy(dst, file)
    if err != nil {
        log.Printf("[ERROR] Save file: %v", err)
        http.Error(w, "Unable to save file", http.StatusInternalServerError)
        return
    }

    log.Printf("[SUCCESS] Uploaded: %s (%d bytes)", handler.Filename, handler.Size)
    fmt.Fprintf(w, "File uploaded to /tmp/uploads/%s", handler.Filename)
}

func main() {
    log.Println("[START] Server running at http://localhost:8080/upload")
    log.Println("[INFO] Ctrl+C to stop. Make sure port 8080 is allowed in firewall.")

    http.HandleFunc("/upload", uploadHandler)
    srv := &http.Server{Addr: ":8080"}

    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("[FATAL] ListenAndServe: %v", err)
        }
    }()

    <-stop
    log.Println("[SHUTDOWN] Gracefully stopping...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("[ERROR] Shutdown failed: %v", err)
    }

    log.Println("[EXIT] Server stopped.")
}

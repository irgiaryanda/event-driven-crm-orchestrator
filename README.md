# Event-Driven CRM Orchestrator

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-3.44-003B57?style=for-the-badge&logo=sqlite&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Production%20Ready-success?style=for-the-badge)

---

## Overview

An AI-powered, event-driven webhook routing system built entirely in Go. This orchestrator captures incoming webhook payloads, intelligently categorizes them using a Large Language Model, dispatches real-time Telegram notifications, and serves a built-in monitoring dashboard—all with zero external dependencies beyond the Go standard library.

---

## System Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│   Webhook   │────▶│   Go Server  │────▶│     SQLite      │
│   Payload   │     │   (Port 8080)│     │   (Persistent)   │
└─────────────┘     └──────┬───────┘     └─────────────────┘
                           │
                    ┌──────▼───────┐
                    │  Goroutines   │
                    ├───────────────┤
                    │  ┌─────────┐  │
                    │  │ LLM API │──┼──▶ AI Categorization
                    │  └─────────┘  │
                    │  ┌─────────┐  │
                    │  │Telegram │──┼──▶ Push Notification
                    │  └─────────┘  │
                    └───────────────┘
```

---

## Key Features

- ⚡ **Pure Go HTTP Server** — No heavy frameworks, uses `net/http` from the standard library
- 🧠 **AI-Powered Payload Routing** — Multi-language support for categorization (SUPPORT, SALES, BILLING, OTHER)
- 📨 **Asynchronous Telegram Notifications** — Non-blocking alert dispatch via goroutines
- 📊 **Built-in Monolithic UI Dashboard** — Vanilla JS + Tailwind CSS, served statically
- 🗄️ **CGO-Free SQLite** — Uses `modernc.org/sqlite` for zero-dependency database operations
- 🔄 **Idempotent Webhook Processing** — Duplicate event IDs are handled gracefully
- 🌐 **CORS Enabled** — Ready for cross-origin frontend integrations

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.21+ (`net/http`) |
| Database | SQLite (`modernc.org/sqlite`) |
| AI Categorization | OpenAI-compatible API (GPT-4o) |
| Notifications | Telegram Bot API |
| Frontend | HTML5 / Vanilla JavaScript / Tailwind CSS |

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- OpenAI-compatible API key
- Telegram Bot Token (optional for notifications)

### Installation

```bash
# Clone the repository
git clone https://github.com/irgiaryanda/event-driven-crm-orchestrator.git
cd event-driven-crm-orchestrator

# Copy environment template
cp .env.example .env

# Edit .env with your credentials
# LLM_API_URL=https://api.openai.com/v1/chat/completions
# LLM_API_KEY=your_api_key_here
# LLM_MODEL=gpt-4o
# TELEGRAM_BOT_TOKEN=your_bot_token_here
# TELEGRAM_CHAT_ID=your_chat_id_here

# Run the server
go run cmd/server/main.go
```

### Access the Dashboard

Open your browser and navigate to:

```
http://localhost:8080
```

---

## Usage / Testing

### Send a Webhook (PowerShell)

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/webhook" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{
    "event_id": "evt-001",
    "data": {
      "message": "I need help with my billing statement",
      "customer": "John Doe"
    }
  }'
```

### Send a Webhook (Bash/cURL)

```bash
curl -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "evt-002",
    "data": {
      "message": "Saya ingin bertanya tentang produk baru",
      "customer": "Jane Smith"
    }
  }'
```

### Retrieve Events

```bash
curl http://localhost:8080/api/events
```

---

## API Reference

### POST /api/webhook

Receives webhook payloads and triggers async processing.

**Request Body:**
```json
{
  "event_id": "optional-custom-id",
  "data": {
    "message": "Webhook payload content",
    "customer": "Customer name"
  }
}
```

**Response (201 Created):**
```json
{
  "status": "created",
  "event_id": "evt-001"
}
```

### GET /api/events

Returns all stored events ordered by creation date (newest first).

**Response:**
```json
[
  {
    "id": "evt-001",
    "payload": "{\"event_id\":\"evt-001\",\"data\":{...}}",
    "status": "PROCESSED",
    "created_at": "2026-07-18T22:00:00Z"
  }
]
```

---

## Project Structure

```
event-driven-crm-orchestrator/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── db/
│   │   └── database.go       # SQLite operations
│   ├── handlers/
│   │   ├── events.go         # GET /api/events handler
│   │   └── webhook.go        # POST /api/webhook handler
│   ├── models/
│   │   └── event.go         # Event data model
│   └── services/
│       ├── llm.go            # LLM categorization
│       └── telegram.go       # Telegram notifications
├── static/
│   └── index.html            # Dashboard UI
├── .env.example              # Environment template
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## License

This project is licensed under the MIT License.
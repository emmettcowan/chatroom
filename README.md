# Go Chatroom

A simple, concurrent TCP chat application built with Go, featuring a Terminal User Interface (TUI) powered by Charm's Bubble Tea.

## Features

- **Concurrent Server**: Supports multiple simultaneous client connections using Go routines.
- **Broadcast Messaging**: Messages sent by any client are broadcasted to all connected participants.
- **TUI Client**: A polished terminal interface built with `bubbletea`, `bubbles`, and `lipgloss`.
- **Real-time Updates**: Asynchronous message reception and UI rendering.

## Project Structure

- `cmd/server/main.go`: Entry point for the chat server.
- `cmd/client/main.go`: Entry point for the TUI chat client.
- `internal/server/`: Core server logic and connection handling.
- `internal/client/`: TUI implementation and client-side networking.

## Prerequisites

- Go 1.26 or later (as specified in `go.mod`).

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/emmettcowan/chatroom.git
cd chatroom
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Run the Server

Start the chat server first. It listens on TCP port `8090` by default.

```bash
go run cmd/server/main.go
```

### 4. Run the Client(s)

Open one or more new terminal windows and run the client:

```bash
go run cmd/client/main.go
```

## How to Use

1. Once the client starts, you'll see a terminal interface with a text input at the bottom.
2. Type your message and press **Enter** to send it.
3. Your message (and messages from others) will appear in the main chat area above.
4. To exit the client, press `Ctrl+C`, `Esc`, or `q`.

## Technologies Used

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): A Go framework based on The Elm Architecture for building terminal applications.
- [Bubbles](https://github.com/charmbracelet/bubbles): TUI components for Bubble Tea.
- [Lip Gloss](https://github.com/charmbracelet/lipgloss): Style definitions for nice terminal layouts.
- Standard Library `net`: For TCP networking.

# NetCat TCP Chat

## Overview

The NetCat TCP Chat is a simple server-client program implemented in Go, designed to facilitate real-time chat communication over TCP. The server allows multiple clients to connect, exchange messages and participate in group discussions. Each message is logged to a file, and the chat history is sent to new clients upon joining.

## Features

- **Multiple Client Support**: Allows multiple users to connect to the chat simultaneously.
- **Chat History**: Sends the previous chat history to new clients when they join.
- **Message Broadcasting**: Broadcasts messages from one client to all others.
- **Logging**: Logs client connections, disconnections and all messages to a file for auditing purposes.
- **Name Validation**: Ensures that client names are unique and valid.
- **Graceful Exit**: Allows users to leave the chat by sending specific commands.

## How It Works

### Server

The server listens for incoming connections on a specified port.

When a client connects:
- The server prompts the client for a username.
- If the username is valid and unique, the client is added to the active user list.
- The chat history is sent to the new client.
- Messages sent by clients are broadcast to all connected users, except the sender.
- The server logs every connection, disconnection and message.

### Client

- The client connects to the server using a TCP connection with a command: **nc localhost 8989**
- After providing a valid username, the client receives the chat history and can start sending messages.
- Messages are displayed in real-time, along with usernames.
- The client can leave the chat by typing `/exit` or `/quit`.

## Usage

### Prerequisites

- Go programming language installed on your system.
- Network access to run the server and client.

### Server Setup

1. Clone the repository.
2. Navigate to the project directory.
3. Run the server:
   ```sh
   go run .

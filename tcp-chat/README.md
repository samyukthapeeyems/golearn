# TCP Chat Server in Go

This project is a straightforward TCP chat server built with Go. It’s a great example of how to handle multiple tasks at once with concurrency, manage clean shutdowns, process various commands, and validate messages.

This chat server lets several users connect at the same time. They can join different chat rooms, exchange messages, and use simple commands to interact. It’s designed to show how you can manage real-time communication and keep things running smoothly, even with multiple users chatting simultaneously.

## Project Structure

The project is organized into the following files:

- **`main.go`**: Initializes and runs the server, handles connections, and manages graceful shutdown.
- **`server.go`**: Contains the main server logic, including room management and command processing.
- **`client.go`**: Manages individual client connections, including message validation and rate limiting.
- **`room.go`**: Manages chat rooms and broadcasts messages to room members.
- **`commands.go`**: Defines command identifiers and structures.


### `main.go`

- **`main` Function**: 
  - Initializes the server with `s := newServer()`.
  - Starts command handling in a separate goroutine with `go s.run()`.
  - Listens for incoming TCP connections on port 8888 with `net.Listen("tcp", ":8888")`.
  - Implements graceful shutdown:
    - Uses `context.WithCancel` to create a shutdown context.
    - Listens for interrupt or terminate signals (`SIGINT`, `SIGTERM`) and closes the listener upon receiving a signal.
    - Ensures all connections are handled before shutting down using `sync.WaitGroup`.

### `server.go`

- **`server` Struct**:
  - Manages chat rooms (`rooms` map) and commands (`commands` channel).
- **Command Handling**:
  - **`run` Method**: Processes incoming commands from clients.
  - **Command Methods**:
    - `name(c *client, args []string)`: Changes the client's name.
    - `join(c *client, args []string)`: Adds the client to a chat room or creates a new room if needed.
    - `listRooms(c *client)`: Lists available chat rooms.
    - `msg(c *client, args []string)`: Sends a message to the current room.
    - `quit(c *client)`: Removes the client from their room and closes the connection.
    - `quitCurrentRoom(c *client)`: Handles client disconnection from their current room.

### `client.go`

- **`client` Struct**:
  - Manages the client connection, name, current room, and message rate limiting.
- **Message Handling**:
  - **`readInput()`**: Reads and processes messages from the client.
  - **Message Validation**:
    - `validateMessage(msg string)`: Checks message length.
    - `canSendMessage()`: Implements rate limiting to prevent spamming.
    - `handleInvalidMessage(err error)`: Handles errors related to message validation.
    - `err(err error)`: Sends error messages to the client.
    - `msg(msg string)`: Sends regular messages to the client.
    - `containsForbiddenWords(msg string)`: Filters out messages with inappropriate content.

### `room.go`

- **`room` Struct**:
  - Represents a chat room with a name and a map of members (`members` map).
- **`broadcast(sender *client, msg string)`**: Sends a message to all room members except the sender.

### `commands.go`

- **`commandID` Enum**: Represents different command types (e.g., `CMD_NAME`, `CMD_JOIN`, `CMD_ROOMS`, `CMD_MSG`, `CMD_QUIT`).
- **`command` Struct**: Encapsulates a command with its ID, client, and arguments.


## Running the Server

1. Clone the repository:
   ```bash
   git clone https://github.com/samyukthapeeyems/golearn/tree/main/tcp-chat
   cd tcp-chat

## Testing

To test the chat server, you can use `telnet` to interact with the server and simulate client connections. Follow these steps:

1. **Build and Run the Server**:

   Open a terminal window and navigate to the project directory. Build and run the server with the following commands:
   ```bash
   go build .
   ./tcp-chat

    //output
    2024/08/19 01:49:18 Started server on port 8888
    2024/08/19 01:49:25 New client has connected: [::1]:38824
    2024/08/19 01:49:29 New client has connected: [::1]:38838
    2024/08/19 01:50:27 Client has disconnected: [::1]:38824
    2024/08/19 01:50:32 Client has disconnected: [::1]:38838

2. **Client x**:

    Open another terminal window. Run the client with the following commands:
   ```bash  
    telnet localhost 8888

    //output
    Trying ::1...
    Connected to localhost.
    Escape character is '^]'.
    /name John
    > Nice! I will call you John
    /join #general
    > Welcome to the room #general!
    > Jane has joined the room
        /msg Hi everyone!
    > Jane: Hey all!
    /quit
    > Sad to see you go :(  Bye...
    Connection closed by foreign host.

3. **Client y**:

    Open another terminal window. Run another client with the following commands:
   ```bash  
    telnet localhost 8888

    //output
    Trying ::1...
    Connected to localhost.
    Escape character is '^]'.
    /name Jane
    > Nice! I will call you Jane
    /rooms
    > Available rooms are: #general
    /join #general
    > Welcome to the room #general!
    > John: Hi everyone!
    /msg Hey all!
    > John has left the room
    /quit
    > Sad to see you go :(  Bye...
    Connection closed by foreign host.


# Fleare CLI Documentation

## Overview
Fleare CLI is a command-line client designed for communicating with an Fleare in-memory server. It allows users to send commands, authenticate, and receive responses in a structured format.

## Installation
To install Fleare CLI, clone the repository and build the binary using Go:

```sh
# Clone the repository
git clone git@github.com:parashmaity/Fleare-cli.git
cd fleare-cli

# Build the executable
go build -o fleare-cli main.go

# Run the CLI
./fleare-cli -host <server_host> -port <server_port> -u <username> -p <password>
```

## Command-Line Arguments
The CLI accepts the following command-line arguments:

| Flag        | Description                                    | Default         |
|------------|--------------------------------|---------------|
| `-host`     | Server host address               | `127.0.0.1`   |
| `-port`     | Server port                       | `4775`        |
| `-u`        | Username for authentication      | Required      |
| `-p`        | Password for authentication      | Required      |

Example usage:
```sh
./fleare-cli -host 192.168.1.1 -port 4775 -u admin -p secret
```

## Features
- **Authentication**: Users must authenticate with a username and password.
- **Command Execution**: Send commands to the server with parameters.
- **Response Parsing**: Parses JSON responses from the server.
- **Command History**: Stores command history in `~/.fleare_history`.
- **Graceful Exit**: Handles `CTRL+C` to prompt for an exit command.

## Authentication
Once the CLI establishes a TCP connection with the server, it sends the username and password for authentication. If the authentication succeeds, the server responds with `CONNECTED`.

## Running Commands
After authentication, users can enter commands using the CLI prompt:
```sh
fleare-cli> GET -k "session123"
fleare-cli> POST -p "/users" -d '{"name":"John"}'
```

### Supported Commands
| Command  | Description                          |
|----------|--------------------------------------|
| `GET`    | Retrieve data using a key (`-k`)    |
| `POST`   | Send data to a specified path (`-p`) with a body (`-d`) |
| `PUT`    | Update an existing record (`-k`, `-p`, `-d`) |
| `DELETE` | Remove data using a key (`-k`)     |
| `exit`   | Close the client gracefully        |

### Example Usage
#### Retrieve Data
```sh
fleare-cli> GET -k "user:1234"
```
#### Create a New User
```sh
fleare-cli> POST -p "/users" -d '{"name":"Alice", "email":"alice@example.com"}'
```

## JSON Response Handling
Responses from the server are parsed into JSON format:
```json
{
    "requestId": "abc123",
    "data": {
        "status": "success",
        "message": "User created successfully"
    }
}
```

The CLI colorizes JSON output for better readability.

## Handling System Signals
- Pressing `CTRL+C` prompts the user to use the `exit` command.
- The CLI listens for termination signals to close the connection gracefully.

## Contribution
Feel free to contribute to the development of Fleare CLI by submitting issues or pull requests.

## License
Fleare CLI is open-source and distributed under the MIT License.

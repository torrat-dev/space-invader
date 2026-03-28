# Space Invader

A real-time space invader game built with Go and HTMX. Control your spaceship, shoot enemies, and avoid getting hit in this classic arcade-style game.

## Installation

1. Ensure you have Go 1.23.1 or later installed.
2. Clone the repository:
   ```
   git clone https://github.com/yourusername/space-invader.git
   cd space-invader
   ```
3. Install dependencies:
   ```
   go mod tidy
   ```

## Usage

Run the server:
```
go run main.go
```

Open your browser and navigate to `http://localhost:8080`.

Use arrow keys to move left/right, space to shoot.

## Build

To build the executable:
```
go build -o space-incremental.exe
```

## Test

Run tests:
```
go test
```

## Project Structure

- `main.go`: Main application code with game logic, HTTP handlers, and WebSocket broadcasting.
- `main_test.go`: Unit tests for the game functionality.
- `go.mod`: Go module file with dependencies.
- `templates/`: HTML templates for the web interface.
  - `index.html`: Main page with game container and controls.
  - `game.html`: Template for rendering the game state.

## Dependencies

- [gorilla/websocket](https://github.com/gorilla/websocket): For WebSocket connections.
- [HTMX](https://htmx.org/): For dynamic HTML updates without JavaScript.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project does not have a license specified.
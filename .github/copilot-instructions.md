# Project Guidelines

## Code Style
- Follow standard Go formatting (`go fmt`)
- Use meaningful variable names, e.g., `state` for game state, `mu` for mutex
- Reference [main.go](../main.go) for struct definitions and function organization

## Architecture
- Monolithic web application with game loop running in background goroutine
- Centralized `GameState` struct protected by mutex for thread-safe access
- Server-side rendering: HTML templates sent via WebSocket for real-time UI updates
- Components: HTTP handlers for actions, WebSocket for broadcasting, game logic in update loop

## Build and Test
- Run: `go run main.go` (starts server on localhost:8080)
- Test: `go test`
- Build: `go build -o space-incremental.exe`

## Conventions
- State mutations always lock `mu` before access
- Templates in `templates/` directory, loaded on init
- Collision detection uses bounding box with 10-20px tolerances
- Enemy firing from last enemy in slice (not random)
- Game tick at 60 FPS (16ms intervals)
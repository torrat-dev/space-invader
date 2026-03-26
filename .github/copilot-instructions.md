# Project Guidelines

## Code Style
- Use Go standard library for HTTP server and HTML templating.
- HTML with inline CSS for minimal dependencies.
- Leverage HTMX for client-side interactivity without heavy JavaScript frameworks.

## Architecture
- Server-driven architecture where the Go backend maintains the single source of truth for game state.
- GameState struct encapsulates player, enemies, bullets, score, and game over status.
- Frontend uses HTMX polling (GET /game every 200ms) to update the game board.
- User inputs handled via POST endpoints for movement and shooting.

## Build and Test
- Install dependencies: `go mod tidy`
- Run development server: `go run main.go`
- Build executable: `go build`
- No automated tests implemented yet; manual testing via browser at http://localhost:8080

## Conventions
- Entity struct represents positions with X and Y integers.
- Global state protected by sync.Mutex for thread safety in HTTP handlers.
- HTML templates stored in `templates/` directory, parsed with `template.ParseGlob`.
- Custom agent "Web Game Developer (HTMX + Go)" available for game-specific development tasks.
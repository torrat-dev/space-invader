---
name: game-refactor
description: '**WORKFLOW SKILL** — Refactor game code to improve specified features or functions. USE FOR: enhancing game mechanics, optimizing code structure, fixing bugs in game logic. DO NOT USE FOR: general Go refactoring, non-game projects. INVOKES: file system tools (read/write code), run_in_terminal for testing, subagents for code analysis.'
---

# Game Refactor

## Workflow Steps

1. **Identify the Feature/Function**: Clarify which specific game feature or function needs improvement (e.g., enemy movement, collision detection, scoring).

2. **Analyze Current Code**: Read and understand the existing implementation using file reading tools and semantic search.

3. **Plan Improvements**: Identify code smells, performance issues, or logic errors. Suggest better structures like separating concerns, improving readability, or optimizing algorithms.

4. **Implement Refactor**: Make targeted code changes using replace_string_in_file, ensuring to follow Go conventions (go fmt, meaningful names).

5. **Test Changes**: Run the game and tests to validate the refactor works correctly. Use run_in_terminal for `go run` or `go test`.

6. **Validate and Iterate**: Check for any regressions, ensure game loop runs smoothly, and iterate if needed.

## Quality Criteria

- Code follows Go formatting (`go fmt`)
- Meaningful variable names (e.g., `state` for game state)
- Thread-safe access to GameState with mutex
- Collision detection with appropriate tolerances
- Game tick at 60 FPS

## Go Best Practices

- **Formatting and Style**: Always run `go fmt` to ensure consistent formatting. Use `gofmt` or `goimports` for imports.
- **Naming Conventions**: Use camelCase for variables/functions, PascalCase for exported types. Keep names descriptive but concise.
- **Error Handling**: Check and handle errors appropriately; avoid ignoring them. Use `if err != nil` patterns.
- **Concurrency**: Use channels and goroutines for concurrent operations. Protect shared state with mutexes (e.g., `sync.Mutex`).
- **Interfaces**: Prefer interfaces for abstraction; use `io.Reader` or custom interfaces to decouple code.
- **Slices and Maps**: Initialize with `make()` for better performance. Use `append()` carefully to avoid unnecessary allocations.
- **Avoid Globals**: Minimize global variables; pass dependencies explicitly.
- **Testing**: Write unit tests with `testing` package. Use table-driven tests for multiple cases.
- **Performance**: Profile with `pprof` if needed. Avoid premature optimization, but be mindful of allocations in hot paths.
- **Documentation**: Add comments to exported functions/types using `// FunctionName ...` format.

## Assets

- Reference [main.go](main.go) for struct definitions and function organization
- Use templates in `templates/` for UI updates
- Build with `go build -o space-incremental.exe`

## Examples

### Improving Enemy AI: Random Firing

**Before (current code):**
```go
// Enemies fire bullets
state.EnemyFireCount++
if state.EnemyFireCount > 20 && len(state.Enemies) > 0 {
    state.EnemyFireCount = 0
    // Pick a random enemy to fire
    randomEnemy := state.Enemies[len(state.Enemies)-1]
    state.EnemyBullets = append(state.EnemyBullets, Entity{X: randomEnemy.X + 5, Y: randomEnemy.Y + 20})
}
```

**After (improved with random selection):**
```go
import (
    // ... existing imports
    "math/rand"
    "time"
)

// In init(), add seeding:
rand.Seed(time.Now().UnixNano())

// Enemies fire bullets
state.EnemyFireCount++
if state.EnemyFireCount > 20 && len(state.Enemies) > 0 {
    state.EnemyFireCount = 0
    // Pick a truly random enemy to fire
    randomIndex := rand.Intn(len(state.Enemies))
    randomEnemy := state.Enemies[randomIndex]
    state.EnemyBullets = append(state.EnemyBullets, Entity{X: randomEnemy.X + 5, Y: randomEnemy.Y + 20})
}
```

This follows Go best practices for randomness, making enemy firing unpredictable and more engaging.
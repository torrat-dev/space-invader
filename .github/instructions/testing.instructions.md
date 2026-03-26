---
description: "Use when implementing or modifying testing and validation protocols for the Go + HTMX Space Invaders project to ensure stability and correctness."
name: "Testing Strategy and Validation Protocol"
applyTo: "**/*.go"
---

## 1. Testing Philosophy
Because this game relies on a **Server-Side Source of Truth**, the bulk of our testing will occur in the Go backend. The frontend is a dumb terminal rendering HTML fragments, meaning we do not need heavy frontend testing frameworks (like Cypress or Selenium). 

Our testing strategy is divided into three layers:
1.  **Unit Tests (Go):** Verify the core game math, state mutations, and physics (collisions, boundaries).
2.  **HTTP Handler Tests (Go):** Verify that endpoints correctly parse inputs, update state, and return the correct HTTP status codes and HTML fragments.
3.  **Manual/Visual Validation:** Clear instructions provided to the user to visually verify the HTMX mechanics in the browser.

---

## 2. Unit Testing Game Logic (Go)
Use Go's standard `testing` package. Every time you introduce a new game mechanic (movement, shooting, collisions), you must write a corresponding unit test in a `_test.go` file.

**Example: Testing Player Movement Bounds**
```go
// game_test.go
package main

import "testing"

func TestPlayerMovement(t *testing.T) {
    // Setup initial state
    state := GameState{Player: Entity{X: 10, Y: 20}}
    
    // Simulate moving left
    movePlayerLeft(&state)
    if state.Player.X >= 10 {
        t.Errorf("Expected player X to decrease, got %d", state.Player.X)
    }

    // Simulate boundary check (assuming left boundary is X=0)
    state.Player.X = 0
    movePlayerLeft(&state)
    if state.Player.X < 0 {
        t.Errorf("Player moved out of bounds! X = %d", state.Player.X)
    }
}
```

---

## 3. HTTP Endpoint Testing (`httptest`)
Use the `net/http/httptest` package to simulate HTMX requests. You must verify that `POST` requests mutate the state correctly and that `GET` requests return valid HTML fragments.

**Example: Testing the `/move` Endpoint**
```go
// handlers_test.go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestMoveHandler(t *testing.T) {
    // Reset state
    state = GameState{Player: Entity{X: 280, Y: 20}}
    
    // Create a request to move left
    req, err := http.NewRequest("POST", "/move?dir=left", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(moveHandler)
    handler.ServeHTTP(rr, req)

    // Check status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Verify state mutation
    if state.Player.X >= 280 {
        t.Errorf("Expected player to move left from 280, got %d", state.Player.X)
    }
}
```

---

## 4. Validating HTML Output
When testing the `/game` endpoint, ensure the returned HTML fragment accurately reflects the `GameState`. You can do this by checking for specific substrings in the response body.

```go
func TestGameRender(t *testing.T) {
    state = GameState{Player: Entity{X: 150, Y: 20}}
    // ... setup httptest for GET /game ...
    
    body := rr.Body.String()
    expectedCSS := "left: 150px"
    if !strings.Contains(body, expectedCSS) {
        t.Errorf("Expected HTML fragment to contain %q", expectedCSS)
    }
}
```

---

## 5. Validation Checklist for Every Change
Before concluding any implementation phase, you **must** complete the following checklist:

1.  **Code Compiles:** Run `go build`. Ensure there are no syntax errors.
2.  **State Logic Unit Tested:** Are the mathematical boundaries and collisions tested? (Run `go test`).
3.  **Handlers Tested:** Do the HTTP endpoints return `200 OK` (or appropriate errors) and mutate state correctly?
4.  **Race Conditions Checked:** Because HTMX polls concurrently while user inputs arrive, run tests with the race detector: `go test -race`.
5.  **Provide Manual Test Steps:** You must output a brief set of manual instructions for the human developer. 

**Example Manual Test Output format:**
> ### Manual Verification Steps
> 1. Run `go run main.go`.
> 2. Open `http://localhost:8080`.
> 3. Press the 'Left' and 'Right' buttons (or arrow keys).
> 4. Verify the green player box moves without leaving the black container.
> 5. Check the terminal for any unexpected error logs.

---
name: game-testing
description: '**WORKFLOW SKILL** — Generate test cases and run automated tests for the space-invader game until they pass. USE FOR: testing game functionality, validating game logic, ensuring game runs correctly. DO NOT USE FOR: general Go testing, non-game projects.'
---

# Game Testing Skill

This skill provides a multi-step workflow for testing the space-invader game, focusing on generating comprehensive test cases and running automated tests iteratively until all tests pass.

## Workflow Steps

1. **Code Analysis**: Review the game code (`main.go`) to identify key components: player movement, bullet shooting, enemy spawning/movement/firing, collision detection, scoring, game over conditions, and WebSocket broadcasting.

2. **Test Case Generation**: Generate test cases covering:
   - Player movement (left/right with boundaries)
   - Bullet shooting and movement
   - Enemy spawning, movement (left/right/down), and firing
   - Collision detection (bullets vs enemies, enemy bullets vs player)
   - Scoring and score updates
   - Game over conditions (player hit, enemies reach bottom)
   - New wave spawning when enemies cleared
   - WebSocket state broadcasting
   - HTTP handlers for actions (move, shoot, restart)

3. **Test Implementation**: Update `main_test.go` with test functions using Go's testing framework. Ensure tests initialize state properly, use mutex locks, and cover edge cases.

4. **Test Execution**: Run `go test` to execute tests. Capture output and identify any failures.

5. **Iteration Loop**: If tests fail:
   - Analyze failure reasons (e.g., logic bugs, boundary issues)
   - Fix code in `main.go` or adjust test expectations
   - Rerun tests
   - Repeat until all tests pass

6. **Validation**: Confirm all tests pass and optionally verify game runs via `go run main.go`.

## Assets

- Test templates (`test_templates.go`) with example test functions for key features
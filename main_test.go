package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGameHandlerBulletsRender(t *testing.T) {
	mu.Lock()
	state = GameState{
		Player:  Entity{X: 280, Y: 350},
		Bullets: []Entity{{X: 150, Y: 100}},
		Enemies: []Entity{},
		Score:   0,
	}
	state.GameOver = false
	mu.Unlock()

	req := httptest.NewRequest("GET", "/game", nil)
	w := httptest.NewRecorder()

	gameHandler(w, req)

	res := w.Result()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	html := string(body)

	if !strings.Contains(html, `class="bullet"`) {
		t.Fatalf("expected bullet element in game HTML, got:\n%s", html)
	}

	if strings.Contains(html, `>|</div>`) {
		t.Fatalf("expected bullet to be rendered without text char, got fallback markup then: %s", html)
	}

	if !strings.Contains(html, `width: 4px`) && !strings.Contains(html, `width:4px`) {
		t.Fatalf("expected bullet CSS size to be defined in game template, got:\n%s", html)
	}
}
func TestShootThenGameShowsBullet(t *testing.T) {
	mu.Lock()
	state = GameState{
		Player:         Entity{X: 280, Y: 350},
		Bullets:        []Entity{},
		Enemies:        []Entity{},
		Score:          0,
		EnemyDirection: 1,
		EnemySpeed:     5,
	}
	state.GameOver = false
	mu.Unlock()

	shootReq := httptest.NewRequest("POST", "/shoot", nil)
	shootW := httptest.NewRecorder()
	shootHandler(shootW, shootReq)

	if len(state.Bullets) != 1 {
		t.Fatalf("shoot handler did not add bullet, bullets=%d", len(state.Bullets))
	}

	gameReq := httptest.NewRequest("GET", "/game", nil)
	gameW := httptest.NewRecorder()
	gameHandler(gameW, gameReq)

	res := gameW.Result()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read game response body: %v", err)
	}

	if !strings.Contains(string(body), `class="bullet"`) {
		t.Fatalf("expected game output to contain bullet after shoot, got:\n%s", string(body))
	}
}

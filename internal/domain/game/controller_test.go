//go:build integration

package game_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/viniciusfal/placar/internal/domain/game"
	"github.com/viniciusfal/placar/internal/platform/mongodb"
	"github.com/viniciusfal/placar/internal/testutil"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	db := testutil.SetupTestDB(t)
	repo := mongodb.NewGameRepository(db)
	srv := game.NewService(repo)
	ctrl := game.NewController(srv)

	r := gin.New()
	v1 := r.Group("v1")
	ctrl.RegisterRoutes(v1)

	return r
}

func TestCreateGame_E2E(t *testing.T) {
	r := setupTestRouter(t)

	body := `{"jogador1": "vinicius", "jogador2": "jose", "pontuacao_necessaria": 10 }`
	req := httptest.NewRequest(http.MethodPost, "/v1/games", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateGame_PayloadInvalido(t *testing.T) {
	router := setupTestRouter(t)

	body := `{"jogador1":"","jogador2":"jose","pontuacao_necessaria":10}`
	req := httptest.NewRequest(http.MethodPost, "/v1/games", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

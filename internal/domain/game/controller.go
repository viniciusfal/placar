package game

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciusfal/placar/internal/apierror"
	"github.com/viniciusfal/placar/internal/middleware"
)

type Controller struct {
	srv *Service
}

func NewController(srv *Service) *Controller {
	return &Controller{
		srv: srv,
	}
}

func (c *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("games", c.Create)
	r.GET("/games/:id", c.GetByID)
}

type createGameRequest struct {
	Jogador1            string `json:"jogador1" binding:"required"`
	Jogador2            string `json:"jogador2" binding:"required"`
	PontuacaoNecessaria int    `json:"pontuacao_necessaria" binding:"required"`
}

func (c *Controller) Create(ctx *gin.Context) {
	log := middleware.LoggerFromContext(ctx.Request.Context())

	var req createGameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		p := apierror.New(http.StatusBadRequest, "Payload inválido", err.Error())
		apierror.Respond(ctx, log, p)
		return
	}

	g, err := c.srv.Create(ctx.Request.Context(), CreateGameInput(req))
	if err != nil {
		apierror.Respond(ctx, log, mapError(err))
		return
	}

	ctx.JSON(http.StatusCreated, toResponse(g))
}

func (c *Controller) GetByID(ctx *gin.Context) {
	log := middleware.LoggerFromContext(ctx.Request.Context())

	g, err := c.srv.repo.BuscarPorID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		apierror.Respond(ctx, log, mapError(err))
		return
	}

	ctx.JSON(http.StatusOK, toResponse(g))
}

type gameResponse struct {
	ID                  string `json:"id"`
	Jogador1            string `json:"jogador1"`
	Jogador2            string `json:"jogador2"`
	StatusPartida       string `json:"status_partida"`
	PontuacaoNecessaria int    `json:"pontuacao_necessaria"`
}

func toResponse(g *Game) gameResponse {
	return gameResponse{
		ID:                  g.ID.String(),
		Jogador1:            g.Jogador1,
		Jogador2:            g.Jogador2,
		StatusPartida:       g.StatusPartida,
		PontuacaoNecessaria: g.PontuacaoNecessaria,
	}
}

func mapError(err error) *apierror.Problem {
	switch {
	case errors.Is(err, ErrJogador1Obrigatorio),
		errors.Is(err, ErrJogador2Obrigatorio),
		errors.Is(err, ErrPontuacaoNecessariaInvalida):
		return apierror.New(http.StatusBadRequest, "Erro de validação", err.Error())

	case errors.Is(err, ErrGameNaoEncontrado):
		return apierror.New(http.StatusNotFound, "Partida não encontrada", err.Error())

	default:
		return apierror.New(http.StatusInternalServerError, "Erro interno do servidor", "")
	}
}

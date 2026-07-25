package apierror

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ContentTypeProblemJSON = "application/problem+json"

// Problem segue RFC 7807 (Problem Details for HTTP APIs).
// É o ÚNICO formato de erro usado em toda a aplicação, não importa o domínio.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	TraceID  string `json:"trace_id,omitempty"`
}

func New(status int, title, detail string) *Problem {
	if status < http.StatusBadRequest {
		panic(fmt.Sprintf("apierror.New chamado com status não-erro: %d", status))
	}

	return &Problem{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
	}
}

// Respond escreve o Problem como resposta HTTP.
// Preenche trace_id e instance automaticamente, loga com severidade
// adequada (Error para 5xx, Warn para 4xx), e aborta o handler.
func Respond(c *gin.Context, log *slog.Logger, p *Problem) {
	if traceID, ok := c.Get("trace_id"); ok {
		if id, ok := traceID.(string); ok {
			p.TraceID = id
		}
	}
	p.Instance = c.Request.URL.Path

	if p.Status >= http.StatusInternalServerError {
		log.Error(
			"erro interno",
			slog.String("title", p.Title),
			slog.String("trace_id", p.TraceID),
			slog.Int("status", p.Status),
		)
	} else {
		log.Warn(
			"erro de cliente",
			slog.String("title", p.Title),
			slog.String("trace_id", p.TraceID),
			slog.Int("status", p.Status),
		)
	}

	c.Header("Content-Type", ContentTypeProblemJSON)
	c.AbortWithStatusJSON(p.Status, p)
}

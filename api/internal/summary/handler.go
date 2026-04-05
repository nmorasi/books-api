package summary

import (
	"books-api/internal/annotation"
	"books-api/internal/middleware"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/go-chi/chi/v5"
)

type Character struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

type SummaryResponse struct {
	Characters []Character `json:"characters"`
}

type Handler struct {
	annotationRepo *annotation.Repository
}

func NewHandler(repo *annotation.Repository) *Handler {
	return &Handler{annotationRepo: repo}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/books/{bookID}/summary", h.GetSummary)
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int64)
	bookID, err := strconv.ParseInt(chi.URLParam(r, "bookID"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	annotations, err := h.annotationRepo.ListByBook(r.Context(), bookID, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(annotations) == 0 {
		writeJSON(w, http.StatusOK, SummaryResponse{Characters: []Character{}})
		return
	}

	var sb strings.Builder
	for i, a := range annotations {
		sb.WriteString(strconv.Itoa(i+1))
		sb.WriteString(". ")
		sb.WriteString(a.Body)
		sb.WriteString("\n")
	}

	prompt := `You are a reading assistant. Below are a reader's annotations about a book.

If the book has characters (fiction or biography), identify all characters mentioned. A character may have different names or aliases — treat them as one, using the most common name. For each, write a concise summary (2-4 sentences) of what the reader observed about them.

If the book has no characters (e.g. technical, self-help), identify the key concepts, ideas, or people mentioned instead and summarize what the reader noted about each.

You MUST respond ONLY with valid JSON, no explanation, no markdown, exactly this format:
{"characters":[{"name":"Name or Concept","summary":"What the reader noted."}]}

If nothing meaningful is mentioned, return: {"characters":[]}

Annotations:
` + sb.String()

	client := anthropic.NewClient(
		option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
	)

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate summary")
		return
	}

	if len(msg.Content) == 0 {
		writeError(w, http.StatusInternalServerError, "empty response from AI")
		return
	}

	raw := msg.Content[0].Text
	// Strip markdown code fences if present
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```json")
		raw = strings.TrimPrefix(raw, "```")
		raw = strings.TrimSuffix(raw, "```")
		raw = strings.TrimSpace(raw)
	}

	var result SummaryResponse
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		log.Printf("AI parse error: %v\nraw response: %s", err, raw)
		writeError(w, http.StatusInternalServerError, "failed to parse AI response")
		return
	}

	if result.Characters == nil {
		result.Characters = []Character{}
	}
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

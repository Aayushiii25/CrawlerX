package scoring

import "context"

type URLMetadata struct {
	Domain      string `json:"domain"`
	Depth       int    `json:"depth"`
	ParentURL   string `json:"parent_url"`
	ContentType string `json:"content_type,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type ScoreResult struct {
	Priority float64 `json:"priority"`

	Confidence float64 `json:"confidence"`

	Reason string `json:"reason"`

	IsSpam bool `json:"is_spam"`
}

type Scorer interface {
	Score(ctx context.Context, url string, metadata URLMetadata) (ScoreResult, error)

	Name() string
}

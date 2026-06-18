package scoring

import (
	"context"
)

type LLMScorer struct {
	provider string
	apiKey   string
	model    string
	fallback Scorer
}

type LLMConfig struct {
	Provider string
	APIKey   string
	Model    string
}

func NewLLMScorer(cfg LLMConfig, fallback Scorer) *LLMScorer {
	return &LLMScorer{
		provider: cfg.Provider,
		apiKey:   cfg.APIKey,
		model:    cfg.Model,
		fallback: fallback,
	}
}

func (ls *LLMScorer) Name() string {
	if ls.provider == "" {
		return "llm-scorer(disabled→" + ls.fallback.Name() + ")"
	}
	return "llm-scorer(" + ls.provider + ")"
}

func (ls *LLMScorer) Score(ctx context.Context, url string, metadata URLMetadata) (ScoreResult, error) {

	if ls.provider == "" || ls.apiKey == "" {
		return ls.fallback.Score(ctx, url, metadata)
	}

	return ls.fallback.Score(ctx, url, metadata)
}

func buildScoringPrompt(url string, metadata URLMetadata) string {
	return `You are a web crawl prioritization engine. Evaluate the following URL and return a JSON object.

URL: ` + url + `
Domain: ` + metadata.Domain + `
Crawl Depth: ` + string(rune('0'+metadata.Depth)) + `
Parent URL: ` + metadata.ParentURL + `

Evaluate based on:
1. Content value: Is this likely to contain useful, original content?
2. Spam likelihood: Does this URL pattern suggest spam, ads, or low-quality content?
3. Domain authority: Is this domain known for quality content?
4. URL structure: Does the path suggest valuable pages (docs, articles, products)?

Return ONLY this JSON:
{
  "priority": 0.0-1.0 (higher = more valuable to crawl),
  "confidence": 0.0-1.0 (how confident you are),
  "is_spam": true/false,
  "reason": "brief explanation"
}`
}

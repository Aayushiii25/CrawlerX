package scoring

import (
	"context"
	"testing"
)

func TestRuleScorer_SeedURL(t *testing.T) {
	scorer := NewRuleScorer()
	ctx := context.Background()

	result, err := scorer.Score(ctx, "https://example.com/", URLMetadata{
		Domain: "example.com",
		Depth:  0,
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}

	if result.Priority < 0.7 {
		t.Errorf("Expected high priority for seed URL, got %.2f", result.Priority)
	}
	t.Logf("Seed URL: priority=%.2f confidence=%.2f reason=%s", result.Priority, result.Confidence, result.Reason)
}

func TestRuleScorer_HighValueDomain(t *testing.T) {
	scorer := NewRuleScorer()
	ctx := context.Background()

	result, err := scorer.Score(ctx, "https://mit.edu/research/papers", URLMetadata{
		Domain: "mit.edu",
		Depth:  1,
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}

	if result.Priority < 0.6 {
		t.Errorf("Expected high priority for .edu domain, got %.2f", result.Priority)
	}
	t.Logf(".edu domain: priority=%.2f reason=%s", result.Priority, result.Reason)
}

func TestRuleScorer_LowValuePath(t *testing.T) {
	scorer := NewRuleScorer()
	ctx := context.Background()

	result, err := scorer.Score(ctx, "https://example.com/login?redirect=/home", URLMetadata{
		Domain: "example.com",
		Depth:  2,
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}

	if result.Priority > 0.5 {
		t.Errorf("Expected low priority for login page, got %.2f", result.Priority)
	}
	t.Logf("Login page: priority=%.2f reason=%s", result.Priority, result.Reason)
}

func TestRuleScorer_SpamDetection(t *testing.T) {
	scorer := NewRuleScorer()
	ctx := context.Background()

	result, err := scorer.Score(ctx, "https://example.com/free-money-casino-viagra", URLMetadata{
		Domain: "example.com",
		Depth:  3,
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}

	if !result.IsSpam {
		t.Error("Expected spam detection for URL with spam keywords")
	}
	t.Logf("Spam URL: priority=%.2f spam=%v reason=%s", result.Priority, result.IsSpam, result.Reason)
}

func TestRuleScorer_DeepCrawl(t *testing.T) {
	scorer := NewRuleScorer()
	ctx := context.Background()

	result, err := scorer.Score(ctx, "https://example.com/a/b/c/d/e/f/g/h", URLMetadata{
		Domain: "example.com",
		Depth:  5,
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}

	if result.Priority > 0.5 {
		t.Errorf("Expected low priority for deep path, got %.2f", result.Priority)
	}
	t.Logf("Deep path: priority=%.2f reason=%s", result.Priority, result.Reason)
}

func TestRuleScorer_DocumentationPage(t *testing.T) {
	scorer := NewRuleScorer()
	ctx := context.Background()

	result, err := scorer.Score(ctx, "https://example.com/docs/getting-started", URLMetadata{
		Domain: "example.com",
		Depth:  1,
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}

	if result.Priority < 0.5 {
		t.Errorf("Expected decent priority for docs page, got %.2f", result.Priority)
	}
	t.Logf("Docs page: priority=%.2f reason=%s", result.Priority, result.Reason)
}

func TestLLMScorer_FallbackWhenDisabled(t *testing.T) {
	ruleScorer := NewRuleScorer()
	llmScorer := NewLLMScorer(LLMConfig{}, ruleScorer)
	ctx := context.Background()

	result, err := llmScorer.Score(ctx, "https://example.com/", URLMetadata{
		Domain: "example.com",
		Depth:  0,
	})
	if err != nil {
		t.Fatalf("Score failed: %v", err)
	}

	ruleResult, _ := ruleScorer.Score(ctx, "https://example.com/", URLMetadata{
		Domain: "example.com",
		Depth:  0,
	})

	if result.Priority != ruleResult.Priority {
		t.Errorf("LLM scorer (disabled) should match rule scorer: got %.2f, expected %.2f",
			result.Priority, ruleResult.Priority)
	}

	expectedName := "llm-scorer(disabled→rule-scorer)"
	if llmScorer.Name() != expectedName {
		t.Errorf("Expected name %s, got %s", expectedName, llmScorer.Name())
	}
}

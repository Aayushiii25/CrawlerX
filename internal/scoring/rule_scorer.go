package scoring

import (
	"context"
	"net/url"
	"path"
	"strings"
)

type RuleScorer struct{}

func NewRuleScorer() *RuleScorer {
	return &RuleScorer{}
}

func (rs *RuleScorer) Name() string {
	return "rule-scorer"
}

func (rs *RuleScorer) Score(ctx context.Context, rawURL string, metadata URLMetadata) (ScoreResult, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ScoreResult{Priority: 0.3, Confidence: 0.5, Reason: "unparseable URL"}, nil
	}

	score := 0.5
	reasons := make([]string, 0)
	rulesApplied := 0

	tld := extractTLD(metadata.Domain)
	if tldScore, ok := tldScores[tld]; ok {
		score += tldScore
		reasons = append(reasons, "tld:"+tld)
		rulesApplied++
	}

	pathDepth := countPathSegments(parsed.Path)
	if pathDepth == 0 || pathDepth == 1 {
		score += 0.15
		reasons = append(reasons, "shallow-path")
	} else if pathDepth > 5 {
		score -= 0.15
		reasons = append(reasons, "deep-path")
	}
	rulesApplied++

	if metadata.Depth == 0 {
		score += 0.2
		reasons = append(reasons, "seed-url")
	} else if metadata.Depth <= 2 {
		score += 0.05
		reasons = append(reasons, "near-seed")
	} else if metadata.Depth > 4 {
		score -= 0.1
		reasons = append(reasons, "deep-crawl")
	}
	rulesApplied++

	numParams := len(parsed.Query())
	if numParams > 3 {
		score -= 0.1
		reasons = append(reasons, "many-params")
	}
	if numParams > 6 {
		score -= 0.15
		reasons = append(reasons, "excessive-params")
	}
	rulesApplied++

	pathLower := strings.ToLower(parsed.Path)
	for _, pattern := range highValuePatterns {
		if strings.Contains(pathLower, pattern) {
			score += 0.1
			reasons = append(reasons, "valuable-pattern:"+pattern)
			rulesApplied++
			break
		}
	}

	for _, pattern := range lowValuePatterns {
		if strings.Contains(pathLower, pattern) {
			score -= 0.30
			reasons = append(reasons, "low-value:"+pattern)
			rulesApplied++
			break
		}
	}

	isSpam, spamReason := detectSpam(rawURL, parsed)
	if isSpam {
		return ScoreResult{
			Priority:   0.0,
			Confidence: 0.8,
			Reason:     "spam: " + spamReason,
			IsSpam:     true,
		}, nil
	}

	ext := strings.ToLower(path.Ext(parsed.Path))
	if ext == "" || ext == ".html" || ext == ".htm" {
		score += 0.05
		reasons = append(reasons, "html-content")
	} else if ext == ".php" || ext == ".asp" || ext == ".aspx" || ext == ".jsp" {

	} else {
		score -= 0.1
		reasons = append(reasons, "non-html:"+ext)
	}
	rulesApplied++

	if score > 1.0 {
		score = 1.0
	}
	if score < 0.0 {
		score = 0.0
	}

	confidence := float64(rulesApplied) / 8.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return ScoreResult{
		Priority:   score,
		Confidence: confidence,
		Reason:     strings.Join(reasons, ", "),
		IsSpam:     false,
	}, nil
}

var tldScores = map[string]float64{
	"edu":   0.15,
	"gov":   0.15,
	"org":   0.10,
	"ac":    0.10,
	"int":   0.05,
	"mil":   0.05,
	"info":  -0.05,
	"biz":   -0.10,
	"xyz":   -0.10,
	"top":   -0.15,
	"click": -0.15,
	"loan":  -0.20,
	"work":  -0.10,
}

var highValuePatterns = []string{
	"/docs/", "/documentation/",
	"/api/", "/reference/",
	"/blog/", "/article/", "/post/",
	"/about", "/team", "/company",
	"/product", "/features",
	"/guide", "/tutorial", "/learn",
	"/wiki/",
}

var lowValuePatterns = []string{
	"/login", "/signin", "/signup", "/register",
	"/cart", "/checkout", "/payment",
	"/admin", "/wp-admin", "/dashboard",
	"/search", "/tag/", "/category/",
	"/page/", "/comment",
	"/print/", "/share/", "/email/",
	"/feed", "/rss", "/atom",
	"/trackback", "/pingback",
}

func detectSpam(rawURL string, parsed *url.URL) (bool, string) {
	urlLower := strings.ToLower(rawURL)

	if len(rawURL) > 2000 {
		return true, "URL too long"
	}

	spamDomains := []string{
		"bit.ly", "tinyurl.com", "goo.gl",
	}
	for _, d := range spamDomains {
		if strings.Contains(parsed.Host, d) {
			return true, "known shortener domain"
		}
	}

	segments := strings.Split(parsed.Path, "/")
	if len(segments) > 15 {
		return true, "excessive path segments"
	}

	spamPatterns := []string{
		"casino", "viagra", "pharma", "crypto-",
		"free-money", "click-here", "make-money",
	}
	for _, pattern := range spamPatterns {
		if strings.Contains(urlLower, pattern) {
			return true, "spam keyword: " + pattern
		}
	}

	return false, ""
}

func countPathSegments(p string) int {
	segments := strings.Split(strings.Trim(p, "/"), "/")
	count := 0
	for _, s := range segments {
		if s != "" {
			count++
		}
	}
	return count
}

func extractTLD(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

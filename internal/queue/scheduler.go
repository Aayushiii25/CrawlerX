package queue

import (
	"context"
	"fmt"

	"github.com/crawlerx/crawlerx/internal/dedup"
	"github.com/crawlerx/crawlerx/internal/eventbus"
	"github.com/crawlerx/crawlerx/internal/pkg/logger"
	"github.com/crawlerx/crawlerx/internal/pkg/urlutil"
	"github.com/crawlerx/crawlerx/internal/scoring"
)

type Scheduler struct {
	queue    *PriorityQueue
	bloom    dedup.BloomFilter
	scorer   scoring.Scorer
	eventBus eventbus.EventBus
	maxDepth int
}

func NewScheduler(
	queue *PriorityQueue,
	bloom dedup.BloomFilter,
	scorer scoring.Scorer,
	eventBus eventbus.EventBus,
	maxDepth int,
) *Scheduler {
	return &Scheduler{
		queue:    queue,
		bloom:    bloom,
		scorer:   scorer,
		eventBus: eventBus,
		maxDepth: maxDepth,
	}
}

type ScheduleResult struct {
	URL       string  `json:"url"`
	Scheduled bool    `json:"scheduled"`
	Reason    string  `json:"reason"`
	Priority  float64 `json:"priority,omitempty"`
}

func (s *Scheduler) Schedule(ctx context.Context, rawURL string, depth int, parentURL string) (*ScheduleResult, error) {

	normalizedURL, err := urlutil.Normalize(rawURL)
	if err != nil || normalizedURL == "" {
		return &ScheduleResult{URL: rawURL, Scheduled: false, Reason: "invalid URL"}, nil
	}

	if !urlutil.IsValidCrawlURL(normalizedURL) {
		return &ScheduleResult{URL: normalizedURL, Scheduled: false, Reason: "not a crawlable URL"}, nil
	}

	if depth > s.maxDepth {
		return &ScheduleResult{URL: normalizedURL, Scheduled: false, Reason: "depth limit exceeded"}, nil
	}

	seen, err := s.bloom.MightContain(ctx, normalizedURL)
	if err != nil {
		return nil, fmt.Errorf("scheduler: bloom check: %w", err)
	}
	if seen {

		if s.eventBus != nil {
			s.eventBus.Publish(ctx, eventbus.Event{
				Type: eventbus.EventDuplicateDetected,
				Payload: eventbus.Payload{
					URL:    normalizedURL,
					Domain: urlutil.ExtractDomain(normalizedURL),
				},
			})
		}
		return &ScheduleResult{URL: normalizedURL, Scheduled: false, Reason: "duplicate"}, nil
	}

	domain := urlutil.ExtractDomain(normalizedURL)
	scoreResult, err := s.scorer.Score(ctx, normalizedURL, scoring.URLMetadata{
		Domain:    domain,
		Depth:     depth,
		ParentURL: parentURL,
	})
	if err != nil {
		logger.Warn(ctx, "scoring failed, using default priority",
			"url", normalizedURL, "error", err)
		scoreResult = scoring.ScoreResult{Priority: 0.5, Confidence: 0.0, Reason: "scoring error, default priority"}
	}

	if scoreResult.IsSpam {
		return &ScheduleResult{URL: normalizedURL, Scheduled: false, Reason: "spam detected: " + scoreResult.Reason}, nil
	}

	if err := s.bloom.Add(ctx, normalizedURL); err != nil {
		return nil, fmt.Errorf("scheduler: bloom add: %w", err)
	}

	task := CrawlTask{
		URL:       normalizedURL,
		Domain:    domain,
		Depth:     depth,
		Priority:  scoreResult.Priority,
		ParentURL: parentURL,
	}

	if err := s.queue.Enqueue(ctx, task); err != nil {
		return nil, fmt.Errorf("scheduler: enqueue: %w", err)
	}

	logger.Debug(ctx, "URL scheduled",
		"url", normalizedURL,
		"priority", scoreResult.Priority,
		"depth", depth,
		"reason", scoreResult.Reason,
	)

	return &ScheduleResult{
		URL:       normalizedURL,
		Scheduled: true,
		Reason:    scoreResult.Reason,
		Priority:  scoreResult.Priority,
	}, nil
}

func (s *Scheduler) ScheduleBatch(ctx context.Context, urls []string, depth int, parentURL string) ([]*ScheduleResult, error) {
	results := make([]*ScheduleResult, 0, len(urls))
	for _, u := range urls {
		result, err := s.Schedule(ctx, u, depth, parentURL)
		if err != nil {
			logger.Warn(ctx, "batch schedule error", "url", u, "error", err)
			results = append(results, &ScheduleResult{URL: u, Scheduled: false, Reason: err.Error()})
			continue
		}
		results = append(results, result)
	}
	return results, nil
}

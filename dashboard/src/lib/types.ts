// API types matching the Go backend responses

export interface MetricsSnapshot {
  total_crawled: number;
  total_duplicates: number;
  total_errors: number;
  active_nodes: number;
  total_nodes: number;
  queue_depth: number;
  retry_queue_size: number;
  dlq_size: number;
  crawl_rate_per_min: number;
  dedup_rate: number;
  error_rate: number;
  uptime: string;
  timestamp: string;
}

export interface NodeInfo {
  node_id: string;
  address: string;
  status: "active" | "failed" | "draining";
  worker_count: number;
  last_seen: string;
  pages_crawled: number;
}

export interface QueueStats {
  priority_queue: number;
  retry_queue: number;
  dead_letter: number;
}

export interface DomainCount {
  domain: string;
  count: number;
}

export interface BloomStats {
  num_bits: number;
  num_hash_funcs: number;
  items_added: number;
  expected_items: number;
  false_positive_rate: number;
  estimated_fp_rate: number;
}

export interface CrawlEvent {
  id: string;
  type: EventType;
  timestamp: string;
  node_id: string;
  payload: EventPayload;
}

export type EventType =
  | "node.joined"
  | "node.failed"
  | "node.left"
  | "page.crawled"
  | "duplicate.detected"
  | "queue.overflow"
  | "crawl.error"
  | "task.assigned"
  | "task.completed";

export interface EventPayload {
  url?: string;
  domain?: string;
  status_code?: number;
  depth?: number;
  worker_count?: number;
  node_addr?: string;
  queue_size?: number;
  queue_name?: string;
  error?: string;
  retry_count?: number;
  duration_ms?: number;
  content_length?: number;
  links_found?: number;
  message?: string;
}

export interface ThroughputPoint {
  timestamp: string;
  crawled: number;
  errors: number;
}

export interface SubmitCrawlResponse {
  submitted: number;
  results: Array<{
    url: string;
    scheduled: boolean;
    reason: string;
    priority?: number;
  }>;
}

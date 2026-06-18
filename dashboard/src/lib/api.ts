import type {
  MetricsSnapshot,
  NodeInfo,
  QueueStats,
  DomainCount,
  BloomStats,
  CrawlEvent,
  ThroughputPoint,
  SubmitCrawlResponse,
} from "./types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function fetchJSON<T>(path: string): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, { cache: "no-store" });
  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }
  return res.json();
}

export const api = {
  getMetrics: () => fetchJSON<MetricsSnapshot>("/api/metrics"),

  getNodes: () =>
    fetchJSON<{ nodes: NodeInfo[]; count: number }>("/api/nodes"),

  getQueueStats: () => fetchJSON<QueueStats>("/api/queue/stats"),

  getTopDomains: (count = 20) =>
    fetchJSON<{ domains: DomainCount[] }>(`/api/domains/top?count=${count}`),

  getBloomStats: () => fetchJSON<BloomStats>("/api/bloom/stats"),

  getRecentEvents: (count = 50) =>
    fetchJSON<{ events: CrawlEvent[]; count: number }>(
      `/api/events/recent?count=${count}`
    ),

  getThroughput: () =>
    fetchJSON<{ throughput: ThroughputPoint[] }>("/api/metrics/throughput"),

  submitCrawl: async (urls: string[]): Promise<SubmitCrawlResponse> => {
    const res = await fetch(`${API_URL}/api/crawl`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ urls, depth: 0 }),
    });
    return res.json();
  },
};

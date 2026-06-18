"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import { Header } from "@/components/header";
import { SystemHealth } from "@/components/system-health";
import { NodeTable } from "@/components/node-table";
import { QueueChart } from "@/components/queue-chart";
import { ThroughputChart } from "@/components/throughput-chart";
import { EventStream } from "@/components/event-stream";
import { DomainDistribution } from "@/components/domain-distribution";
import { useWebSocket } from "@/hooks/use-websocket";
import { useMetrics } from "@/hooks/use-metrics";

export default function Dashboard() {
  const { events, connected } = useWebSocket();
  const {
    metrics,
    nodes,
    queueStats,
    domains,
    bloomStats,
    throughput,
    loading,
    error,
  } = useMetrics(5000);

  // Build queue history from polling snapshots
  const [queueHistory, setQueueHistory] = useState<
    Array<{ time: string; priority: number; retry: number; dlq: number }>
  >([]);

  useEffect(() => {
    if (queueStats) {
      setTimeout(() => {
        setQueueHistory((prev) => {
          const next = [
            ...prev,
            {
              time: new Date().toISOString(),
              priority: queueStats.priority_queue,
              retry: queueStats.retry_queue,
              dlq: queueStats.dead_letter,
            },
          ];
          return next.slice(-60); // Keep last 60 data points (5 min at 5s interval)
        });
      }, 0);
    }
  }, [queueStats]);

  // URL submission
  const [urlInput, setUrlInput] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (!urlInput.trim()) return;

      setSubmitting(true);
      try {
        const urls = urlInput
          .split(/[\n,]/)
          .map((u) => u.trim())
          .filter(Boolean);

        const res = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/crawl`,
          {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ urls, depth: 0 }),
          }
        );

        if (res.ok) {
          setUrlInput("");
          inputRef.current?.focus();
        }
      } catch {
        // Silently fail — user can see errors in event stream
      } finally {
        setSubmitting(false);
      }
    },
    [urlInput]
  );

  if (loading && !metrics) {
    return (
      <div className="min-h-screen bg-zinc-950 flex items-center justify-center">
        <div className="text-zinc-500 text-sm">
          Connecting to CrawlerX coordinator...
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-950 flex flex-col">
      <Header
        connected={connected}
        activeNodes={metrics?.active_nodes ?? 0}
        uptime={metrics?.uptime ?? ""}
      />

      <main className="flex-1 px-6 py-4 space-y-4 max-w-[1800px] mx-auto w-full">
        {/* Error banner */}
        {error && (
          <div className="rounded-md border border-red-800/50 bg-red-950/30 px-4 py-2 text-[12px] text-red-400">
            ⚠ Connection error: {error}. Retrying...
          </div>
        )}

        {/* URL Input */}
        <form
          onSubmit={handleSubmit}
          className="flex items-center gap-3"
        >
          <div className="flex-1 relative">
            <input
              ref={inputRef}
              type="text"
              value={urlInput}
              onChange={(e) => setUrlInput(e.target.value)}
              placeholder="Enter URLs to crawl (comma-separated)..."
              className="w-full rounded-md border border-zinc-800 bg-zinc-900 px-3 py-2 text-[13px] text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-zinc-600 focus:ring-1 focus:ring-zinc-600 font-mono"
              disabled={submitting}
            />
          </div>
          <button
            type="submit"
            disabled={submitting || !urlInput.trim()}
            className="rounded-md bg-zinc-800 px-4 py-2 text-[13px] font-medium text-zinc-200 hover:bg-zinc-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors border border-zinc-700"
          >
            {submitting ? "Submitting..." : "Crawl"}
          </button>
        </form>

        {/* Stats row */}
        <SystemHealth
          metrics={metrics}
          queueStats={queueStats}
          bloomStats={bloomStats}
        />

        {/* Main grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <ThroughputChart data={throughput} />
          <QueueChart queueStats={queueStats} history={queueHistory} />
        </div>

        {/* Node table + Domain distribution */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
          <div className="lg:col-span-2">
            <NodeTable nodes={nodes} />
          </div>
          <DomainDistribution domains={domains} />
        </div>

        {/* Event stream — full width */}
        <EventStream events={events} />
      </main>

      {/* Footer */}
      <footer className="border-t border-zinc-800/50 px-6 py-2 text-[11px] text-zinc-600 flex items-center justify-between">
        <span>CrawlerX v1.0.0</span>
        <span>
          Last update:{" "}
          {metrics?.timestamp
            ? new Date(metrics.timestamp).toLocaleTimeString()
            : "—"}
        </span>
      </footer>
    </div>
  );
}

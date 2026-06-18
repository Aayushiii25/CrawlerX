"use client";

import { useState, useRef, useEffect } from "react";
import type { CrawlEvent } from "@/lib/types";
import { getEventColor, timeAgo } from "@/lib/utils";

interface EventStreamProps {
  events: CrawlEvent[];
}

export function EventStream({ events }: EventStreamProps) {
  const [paused, setPaused] = useState(false);
  const [filter, setFilter] = useState<string>("all");
  const containerRef = useRef<HTMLDivElement>(null);

  // Auto-scroll when not paused
  useEffect(() => {
    if (!paused && containerRef.current) {
      containerRef.current.scrollTop = 0;
    }
  }, [events, paused]);

  const filteredEvents =
    filter === "all"
      ? events
      : events.filter((e) => e.type === filter || e.type.startsWith(filter));

  const eventTypes = [
    { value: "all", label: "All" },
    { value: "page", label: "Crawled" },
    { value: "node", label: "Nodes" },
    { value: "crawl.error", label: "Errors" },
    { value: "duplicate", label: "Dupes" },
  ];

  return (
    <div className="rounded-md border border-zinc-800 bg-zinc-900/50 overflow-hidden flex flex-col h-full">
      {/* Header */}
      <div className="px-4 py-2.5 border-b border-zinc-800 flex items-center justify-between flex-shrink-0">
        <h3 className="text-[12px] font-medium text-zinc-400 uppercase tracking-wider">
          Event Stream
          <span className="ml-2 text-zinc-600">({events.length})</span>
        </h3>
        <div className="flex items-center gap-2">
          {/* Filter buttons */}
          <div className="flex items-center gap-0.5 bg-zinc-800/50 rounded p-0.5">
            {eventTypes.map((type) => (
              <button
                key={type.value}
                onClick={() => setFilter(type.value)}
                className={`px-2 py-0.5 rounded text-[10px] font-medium transition-colors ${
                  filter === type.value
                    ? "bg-zinc-700 text-zinc-200"
                    : "text-zinc-500 hover:text-zinc-300"
                }`}
              >
                {type.label}
              </button>
            ))}
          </div>
          {/* Pause button */}
          <button
            onClick={() => setPaused(!paused)}
            className={`px-2 py-0.5 rounded text-[10px] font-medium transition-colors ${
              paused
                ? "bg-amber-500/20 text-amber-400"
                : "bg-zinc-800 text-zinc-400 hover:text-zinc-200"
            }`}
          >
            {paused ? "▶ Resume" : "⏸ Pause"}
          </button>
        </div>
      </div>

      {/* Event list */}
      <div
        ref={containerRef}
        className="overflow-y-auto flex-1 min-h-0"
        style={{ maxHeight: "400px" }}
      >
        {filteredEvents.length === 0 ? (
          <div className="px-4 py-8 text-center text-[13px] text-zinc-600">
            Waiting for events...
          </div>
        ) : (
          <div className="divide-y divide-zinc-800/30">
            {filteredEvents.map((event, i) => (
              <div
                key={event.id || i}
                className="px-4 py-1.5 hover:bg-zinc-800/20 transition-colors"
              >
                <div className="flex items-start gap-3 text-[12px]">
                  <span className="text-zinc-600 font-mono w-[52px] flex-shrink-0 tabular-nums">
                    {event.timestamp ? timeAgo(event.timestamp) : "now"}
                  </span>
                  <span className={`font-mono flex-shrink-0 w-[120px] ${getEventColor(event.type)}`}>
                    {event.type}
                  </span>
                  <span className="text-zinc-500 truncate flex-1 font-mono">
                    {formatEventPayload(event)}
                  </span>
                  <span className="text-zinc-600 font-mono flex-shrink-0">
                    {event.node_id}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function formatEventPayload(event: CrawlEvent): string {
  const p = event.payload;
  switch (event.type) {
    case "page.crawled":
      return `${p.status_code} ${p.url} (${p.duration_ms}ms, ${p.links_found} links)`;
    case "node.joined":
      return `${p.node_addr} (${p.worker_count} workers)`;
    case "node.failed":
      return p.message || "heartbeat timeout";
    case "node.left":
      return p.message || "graceful shutdown";
    case "duplicate.detected":
      return p.url || "unknown URL";
    case "crawl.error":
      return `${p.url} — ${p.error} (retry ${p.retry_count})`;
    case "queue.overflow":
      return `${p.queue_name}: ${p.queue_size}`;
    default:
      return p.message || p.url || JSON.stringify(p);
  }
}

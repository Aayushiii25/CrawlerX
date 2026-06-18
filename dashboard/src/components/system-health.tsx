"use client";

import { formatNumber } from "@/lib/utils";
import type { MetricsSnapshot, QueueStats, BloomStats } from "@/lib/types";

interface SystemHealthProps {
  metrics: MetricsSnapshot | null;
  queueStats: QueueStats | null;
  bloomStats: BloomStats | null;
}

interface StatCardProps {
  label: string;
  value: string | number;
  subtext?: string;
  status?: "ok" | "warning" | "error" | "neutral";
}

function StatCard({ label, value, subtext, status = "neutral" }: StatCardProps) {
  const borderColor = {
    ok: "border-l-emerald-500",
    warning: "border-l-amber-500",
    error: "border-l-red-500",
    neutral: "border-l-zinc-700",
  }[status];

  return (
    <div
      className={`rounded-md border border-zinc-800 bg-zinc-900/50 px-4 py-3 border-l-2 ${borderColor}`}
    >
      <div className="text-[11px] font-medium text-zinc-500 uppercase tracking-wider mb-1">
        {label}
      </div>
      <div className="text-xl font-semibold font-mono text-zinc-100 tabular-nums">
        {value}
      </div>
      {subtext && (
        <div className="text-[11px] text-zinc-500 mt-0.5">{subtext}</div>
      )}
    </div>
  );
}

export function SystemHealth({ metrics, queueStats, bloomStats }: SystemHealthProps) {
  const m = metrics;

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
      <StatCard
        label="Pages Crawled"
        value={formatNumber(m?.total_crawled ?? 0)}
        status="ok"
      />
      <StatCard
        label="Crawl Rate"
        value={`${m?.crawl_rate_per_min?.toFixed(0) ?? 0}/min`}
        status={
          (m?.crawl_rate_per_min ?? 0) > 0
            ? "ok"
            : "neutral"
        }
      />
      <StatCard
        label="Queue Depth"
        value={formatNumber(queueStats?.priority_queue ?? m?.queue_depth ?? 0)}
        subtext={`Retry: ${formatNumber(queueStats?.retry_queue ?? 0)} · DLQ: ${formatNumber(queueStats?.dead_letter ?? 0)}`}
        status={
          (queueStats?.priority_queue ?? 0) > 100000
            ? "warning"
            : "neutral"
        }
      />
      <StatCard
        label="Dedup Rate"
        value={`${(m?.dedup_rate ?? 0).toFixed(1)}%`}
        subtext={`${formatNumber(m?.total_duplicates ?? 0)} duplicates`}
        status="neutral"
      />
      <StatCard
        label="Error Rate"
        value={`${(m?.error_rate ?? 0).toFixed(1)}%`}
        subtext={`${formatNumber(m?.total_errors ?? 0)} errors`}
        status={(m?.error_rate ?? 0) > 10 ? "error" : "ok"}
      />
      <StatCard
        label="Bloom Filter"
        value={formatNumber(bloomStats?.items_added ?? 0)}
        subtext={`FP: ${((bloomStats?.estimated_fp_rate ?? 0) * 100).toFixed(3)}%`}
        status="neutral"
      />
    </div>
  );
}

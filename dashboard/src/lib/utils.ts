import { clsx, type ClassValue } from "clsx";

export function cn(...inputs: ClassValue[]) {
  return clsx(inputs);
}

export function formatNumber(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
  if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`;
  return n.toLocaleString();
}

export function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`;
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
  return `${(ms / 60000).toFixed(1)}m`;
}

export function timeAgo(date: string): string {
  const seconds = Math.floor(
    (Date.now() - new Date(date).getTime()) / 1000
  );
  if (seconds < 5) return "just now";
  if (seconds < 60) return `${seconds}s ago`;
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  return `${Math.floor(seconds / 86400)}d ago`;
}

export function getEventColor(type: string): string {
  switch (type) {
    case "page.crawled": return "text-emerald-400";
    case "node.joined": return "text-blue-400";
    case "node.failed": return "text-red-400";
    case "node.left": return "text-amber-400";
    case "duplicate.detected": return "text-zinc-500";
    case "crawl.error": return "text-red-500";
    case "queue.overflow": return "text-orange-400";
    default: return "text-zinc-400";
  }
}

export function getStatusColor(status: string): string {
  switch (status) {
    case "active": return "bg-emerald-500";
    case "failed": return "bg-red-500";
    case "draining": return "bg-amber-500";
    default: return "bg-zinc-500";
  }
}

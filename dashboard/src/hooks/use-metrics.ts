"use client";

import { useState, useEffect, useCallback } from "react";
import { api } from "@/lib/api";
import type {
  MetricsSnapshot,
  NodeInfo,
  QueueStats,
  DomainCount,
  BloomStats,
  ThroughputPoint,
} from "@/lib/types";

interface MetricsData {
  metrics: MetricsSnapshot | null;
  nodes: NodeInfo[];
  queueStats: QueueStats | null;
  domains: DomainCount[];
  bloomStats: BloomStats | null;
  throughput: ThroughputPoint[];
  loading: boolean;
  error: string | null;
}

export function useMetrics(intervalMs = 5000): MetricsData {
  const [data, setData] = useState<MetricsData>({
    metrics: null,
    nodes: [],
    queueStats: null,
    domains: [],
    bloomStats: null,
    throughput: [],
    loading: true,
    error: null,
  });

  const fetchAll = useCallback(async () => {
    try {
      const [metrics, nodesRes, queueStats, domainsRes, bloomStats, throughputRes] =
        await Promise.allSettled([
          api.getMetrics(),
          api.getNodes(),
          api.getQueueStats(),
          api.getTopDomains(20),
          api.getBloomStats(),
          api.getThroughput(),
        ]);

      setData({
        metrics: metrics.status === "fulfilled" ? metrics.value : null,
        nodes: nodesRes.status === "fulfilled" ? nodesRes.value.nodes : [],
        queueStats: queueStats.status === "fulfilled" ? queueStats.value : null,
        domains: domainsRes.status === "fulfilled" ? domainsRes.value.domains : [],
        bloomStats: bloomStats.status === "fulfilled" ? bloomStats.value : null,
        throughput: throughputRes.status === "fulfilled" ? throughputRes.value.throughput : [],
        loading: false,
        error: null,
      });
    } catch (err) {
      setData((prev) => ({
        ...prev,
        loading: false,
        error: err instanceof Error ? err.message : "Failed to fetch metrics",
      }));
    }
  }, []);

  useEffect(() => {
    fetchAll();
    const interval = setInterval(fetchAll, intervalMs);
    return () => clearInterval(interval);
  }, [fetchAll, intervalMs]);

  return data;
}

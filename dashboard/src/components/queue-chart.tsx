"use client";

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import type { QueueStats } from "@/lib/types";
import { formatNumber } from "@/lib/utils";

interface QueueChartProps {
  queueStats: QueueStats | null;
  history: Array<{ time: string; priority: number; retry: number; dlq: number }>;
}

export function QueueChart({ queueStats, history }: QueueChartProps) {
  const data = history.length > 0
    ? history
    : [
        {
          time: new Date().toISOString(),
          priority: queueStats?.priority_queue ?? 0,
          retry: queueStats?.retry_queue ?? 0,
          dlq: queueStats?.dead_letter ?? 0,
        },
      ];

  return (
    <div className="rounded-md border border-zinc-800 bg-zinc-900/50 overflow-hidden">
      <div className="px-4 py-2.5 border-b border-zinc-800 flex items-center justify-between">
        <h3 className="text-[12px] font-medium text-zinc-400 uppercase tracking-wider">
          Queue Depth
        </h3>
        <div className="flex items-center gap-4 text-[11px]">
          <span className="flex items-center gap-1.5">
            <span className="h-2 w-2 rounded-full bg-blue-500" />
            <span className="text-zinc-500">Priority</span>
            <span className="font-mono text-zinc-300">
              {formatNumber(queueStats?.priority_queue ?? 0)}
            </span>
          </span>
          <span className="flex items-center gap-1.5">
            <span className="h-2 w-2 rounded-full bg-amber-500" />
            <span className="text-zinc-500">Retry</span>
            <span className="font-mono text-zinc-300">
              {formatNumber(queueStats?.retry_queue ?? 0)}
            </span>
          </span>
          <span className="flex items-center gap-1.5">
            <span className="h-2 w-2 rounded-full bg-red-500" />
            <span className="text-zinc-500">DLQ</span>
            <span className="font-mono text-zinc-300">
              {formatNumber(queueStats?.dead_letter ?? 0)}
            </span>
          </span>
        </div>
      </div>
      <div className="p-4 h-[200px]">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={data}>
            <CartesianGrid strokeDasharray="3 3" stroke="#27272a" />
            <XAxis
              dataKey="time"
              tick={{ fill: "#71717a", fontSize: 10 }}
              tickFormatter={(v) => {
                const d = new Date(v);
                return `${d.getHours()}:${d.getMinutes().toString().padStart(2, "0")}`;
              }}
              stroke="#3f3f46"
            />
            <YAxis
              tick={{ fill: "#71717a", fontSize: 10 }}
              stroke="#3f3f46"
              tickFormatter={formatNumber}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: "#18181b",
                border: "1px solid #3f3f46",
                borderRadius: "6px",
                fontSize: "12px",
              }}
              labelStyle={{ color: "#a1a1aa" }}
            />
            <Area
              type="monotone"
              dataKey="priority"
              stackId="1"
              stroke="#3b82f6"
              fill="#3b82f6"
              fillOpacity={0.15}
              strokeWidth={1.5}
            />
            <Area
              type="monotone"
              dataKey="retry"
              stackId="1"
              stroke="#f59e0b"
              fill="#f59e0b"
              fillOpacity={0.1}
              strokeWidth={1.5}
            />
            <Area
              type="monotone"
              dataKey="dlq"
              stackId="1"
              stroke="#ef4444"
              fill="#ef4444"
              fillOpacity={0.1}
              strokeWidth={1.5}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

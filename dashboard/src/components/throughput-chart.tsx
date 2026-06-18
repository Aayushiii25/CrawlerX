"use client";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import type { ThroughputPoint } from "@/lib/types";

interface ThroughputChartProps {
  data: ThroughputPoint[];
}

export function ThroughputChart({ data }: ThroughputChartProps) {
  const chartData = data.map((p) => ({
    time: p.timestamp,
    crawled: p.crawled,
    errors: p.errors,
  }));

  return (
    <div className="rounded-md border border-zinc-800 bg-zinc-900/50 overflow-hidden">
      <div className="px-4 py-2.5 border-b border-zinc-800 flex items-center justify-between">
        <h3 className="text-[12px] font-medium text-zinc-400 uppercase tracking-wider">
          Throughput
          <span className="ml-1 text-zinc-600 normal-case">(5min window)</span>
        </h3>
        <div className="flex items-center gap-4 text-[11px]">
          <span className="flex items-center gap-1.5">
            <span className="h-2 w-2 rounded-full bg-emerald-500" />
            <span className="text-zinc-500">Crawled</span>
          </span>
          <span className="flex items-center gap-1.5">
            <span className="h-2 w-2 rounded-full bg-red-500" />
            <span className="text-zinc-500">Errors</span>
          </span>
        </div>
      </div>
      <div className="p-4 h-[200px]">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#27272a" />
            <XAxis
              dataKey="time"
              tick={{ fill: "#71717a", fontSize: 10 }}
              tickFormatter={(v) => {
                const d = new Date(v);
                return `${d.getHours()}:${d.getMinutes().toString().padStart(2, "0")}:${d.getSeconds().toString().padStart(2, "0")}`;
              }}
              stroke="#3f3f46"
            />
            <YAxis
              tick={{ fill: "#71717a", fontSize: 10 }}
              stroke="#3f3f46"
              allowDecimals={false}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: "#18181b",
                border: "1px solid #3f3f46",
                borderRadius: "6px",
                fontSize: "12px",
              }}
              labelStyle={{ color: "#a1a1aa" }}
              labelFormatter={(v) => new Date(v).toLocaleTimeString()}
            />
            <Line
              type="monotone"
              dataKey="crawled"
              stroke="#10b981"
              strokeWidth={1.5}
              dot={false}
              activeDot={{ r: 3, stroke: "#10b981" }}
            />
            <Line
              type="monotone"
              dataKey="errors"
              stroke="#ef4444"
              strokeWidth={1.5}
              dot={false}
              activeDot={{ r: 3, stroke: "#ef4444" }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

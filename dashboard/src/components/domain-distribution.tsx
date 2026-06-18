"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import type { DomainCount } from "@/lib/types";
import { formatNumber } from "@/lib/utils";

interface DomainDistributionProps {
  domains: DomainCount[];
}

export function DomainDistribution({ domains }: DomainDistributionProps) {
  const chartData = domains.slice(0, 15).map((d) => ({
    domain: d.domain.length > 20 ? d.domain.slice(0, 18) + "…" : d.domain,
    fullDomain: d.domain,
    count: d.count,
  }));

  return (
    <div className="rounded-md border border-zinc-800 bg-zinc-900/50 overflow-hidden">
      <div className="px-4 py-2.5 border-b border-zinc-800">
        <h3 className="text-[12px] font-medium text-zinc-400 uppercase tracking-wider">
          Top Domains
          <span className="ml-2 text-zinc-600">({domains.length})</span>
        </h3>
      </div>
      <div className="p-4 h-[280px]">
        {chartData.length === 0 ? (
          <div className="flex items-center justify-center h-full text-[13px] text-zinc-600">
            No domains crawled yet
          </div>
        ) : (
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData} layout="vertical" margin={{ left: 10 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#27272a" horizontal={false} />
              <XAxis
                type="number"
                tick={{ fill: "#71717a", fontSize: 10 }}
                stroke="#3f3f46"
                tickFormatter={formatNumber}
              />
              <YAxis
                type="category"
                dataKey="domain"
                tick={{ fill: "#a1a1aa", fontSize: 11 }}
                stroke="#3f3f46"
                width={140}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: "#18181b",
                  border: "1px solid #3f3f46",
                  borderRadius: "6px",
                  fontSize: "12px",
                }}
                labelStyle={{ color: "#a1a1aa" }}
                formatter={(value) => [formatNumber(Number(value ?? 0)), "Pages"]}
                labelFormatter={(_, payload) =>
                  payload?.[0]?.payload?.fullDomain || ""
                }
              />
              <Bar
                dataKey="count"
                fill="#3b82f6"
                fillOpacity={0.7}
                radius={[0, 3, 3, 0]}
                maxBarSize={18}
              />
            </BarChart>
          </ResponsiveContainer>
        )}
      </div>
    </div>
  );
}

"use client";

import { getStatusColor } from "@/lib/utils";

interface HeaderProps {
  connected: boolean;
  activeNodes: number;
  uptime: string;
}

export function Header({ connected, activeNodes, uptime }: HeaderProps) {
  return (
    <header className="border-b border-zinc-800 bg-zinc-950 px-6 py-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <h1 className="text-[15px] font-semibold tracking-tight text-zinc-100">
            CrawlerX
          </h1>
          <span className="text-[11px] font-medium text-zinc-600 uppercase tracking-wider">
            Distributed Crawler
          </span>
        </div>

        <div className="flex items-center gap-6">
          {/* Uptime */}
          <div className="flex items-center gap-2 text-[12px]">
            <span className="text-zinc-500">Uptime</span>
            <span className="font-mono text-zinc-300">{uptime || "—"}</span>
          </div>

          {/* Active Nodes */}
          <div className="flex items-center gap-2 text-[12px]">
            <span className="text-zinc-500">Nodes</span>
            <span className="font-mono text-zinc-300">{activeNodes}</span>
          </div>

          {/* Connection Status */}
          <div className="flex items-center gap-2 text-[12px]">
            <div
              className={`h-2 w-2 rounded-full ${
                connected ? "bg-emerald-500" : "bg-red-500"
              }`}
            />
            <span className={connected ? "text-emerald-400" : "text-red-400"}>
              {connected ? "Live" : "Disconnected"}
            </span>
          </div>
        </div>
      </div>
    </header>
  );
}

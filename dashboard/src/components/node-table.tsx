"use client";

import {
  useReactTable,
  getCoreRowModel,
  flexRender,
  createColumnHelper,
} from "@tanstack/react-table";
import type { NodeInfo } from "@/lib/types";
import { getStatusColor, timeAgo, formatNumber } from "@/lib/utils";

const columnHelper = createColumnHelper<NodeInfo>();

const columns = [
  columnHelper.accessor("status", {
    header: "",
    cell: (info) => (
      <div className="flex items-center justify-center">
        <div className={`h-2 w-2 rounded-full ${getStatusColor(info.getValue())}`} />
      </div>
    ),
    size: 32,
  }),
  columnHelper.accessor("node_id", {
    header: "Node ID",
    cell: (info) => (
      <span className="font-mono text-[13px] text-zinc-200">{info.getValue()}</span>
    ),
  }),
  columnHelper.accessor("status", {
    id: "status_text",
    header: "Status",
    cell: (info) => {
      const status = info.getValue();
      const colorMap: Record<string, string> = {
        active: "text-emerald-400 bg-emerald-500/10",
        failed: "text-red-400 bg-red-500/10",
        draining: "text-amber-400 bg-amber-500/10",
      };
      return (
        <span className={`px-2 py-0.5 rounded text-[11px] font-medium uppercase ${colorMap[status] || "text-zinc-400"}`}>
          {status}
        </span>
      );
    },
  }),
  columnHelper.accessor("worker_count", {
    header: "Workers",
    cell: (info) => (
      <span className="font-mono text-[13px] text-zinc-300">{info.getValue()}</span>
    ),
  }),
  columnHelper.accessor("pages_crawled", {
    header: "Crawled",
    cell: (info) => (
      <span className="font-mono text-[13px] text-zinc-300">
        {formatNumber(info.getValue())}
      </span>
    ),
  }),
  columnHelper.accessor("last_seen", {
    header: "Last Seen",
    cell: (info) => (
      <span className="text-[12px] text-zinc-500">
        {info.getValue() ? timeAgo(info.getValue()) : "—"}
      </span>
    ),
  }),
];

interface NodeTableProps {
  nodes: NodeInfo[];
}

export function NodeTable({ nodes }: NodeTableProps) {
  const table = useReactTable({
    data: nodes,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="rounded-md border border-zinc-800 bg-zinc-900/50 overflow-hidden">
      <div className="px-4 py-2.5 border-b border-zinc-800">
        <h3 className="text-[12px] font-medium text-zinc-400 uppercase tracking-wider">
          Crawler Nodes
          <span className="ml-2 text-zinc-600">({nodes.length})</span>
        </h3>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-left">
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id} className="border-b border-zinc-800/50">
                {headerGroup.headers.map((header) => (
                  <th
                    key={header.id}
                    className="px-4 py-2 text-[11px] font-medium text-zinc-500 uppercase tracking-wider"
                    style={{ width: header.getSize() }}
                  >
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {table.getRowModel().rows.length === 0 ? (
              <tr>
                <td
                  colSpan={columns.length}
                  className="px-4 py-8 text-center text-[13px] text-zinc-600"
                >
                  No nodes registered. Start a crawler node to begin.
                </td>
              </tr>
            ) : (
              table.getRowModel().rows.map((row) => (
                <tr
                  key={row.id}
                  className="border-b border-zinc-800/30 hover:bg-zinc-800/20 transition-colors"
                >
                  {row.getVisibleCells().map((cell) => (
                    <td key={cell.id} className="px-4 py-2">
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

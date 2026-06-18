"use client";

import { useEffect, useRef, useState, useCallback } from "react";
import type { CrawlEvent } from "@/lib/types";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws";

export function useWebSocket() {
  const [events, setEvents] = useState<CrawlEvent[]>([]);
  const [connected, setConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttempts = useRef(0);
  const connectRef = useRef<() => void>(() => {});

  const connect = useCallback(() => {
    try {
      const ws = new WebSocket(WS_URL);

      ws.onopen = () => {
        setConnected(true);
        reconnectAttempts.current = 0;
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data) as CrawlEvent;
          setEvents((prev) => [data, ...prev].slice(0, 200));
        } catch {
          // Ignore malformed messages
        }
      };

      ws.onclose = () => {
        setConnected(false);
        // Exponential backoff reconnect
        const delay = Math.min(
          1000 * Math.pow(2, reconnectAttempts.current),
          30000
        );
        reconnectAttempts.current++;
        reconnectTimeoutRef.current = setTimeout(() => connectRef.current(), delay);
      };

      ws.onerror = () => {
        ws.close();
      };

      wsRef.current = ws;
    } catch {
      // Connection failed — retry
      const delay = Math.min(
        1000 * Math.pow(2, reconnectAttempts.current),
        30000
      );
      reconnectAttempts.current++;
      reconnectTimeoutRef.current = setTimeout(() => connectRef.current(), delay);
    }
  }, []);

  useEffect(() => {
    connectRef.current = connect;
  }, [connect]);

  useEffect(() => {
    connect();
    return () => {
      if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);
      wsRef.current?.close();
    };
  }, [connect]);

  const clearEvents = useCallback(() => setEvents([]), []);

  return { events, connected, clearEvents };
}

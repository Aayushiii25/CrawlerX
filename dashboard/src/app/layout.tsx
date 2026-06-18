import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "CrawlerX — Distributed Web Crawler Dashboard",
  description:
    "Operational dashboard for CrawlerX distributed web crawler. Monitor crawler nodes, queue depth, throughput, and crawl events in real time.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className="min-h-screen bg-zinc-950 antialiased">
        {children}
      </body>
    </html>
  );
}

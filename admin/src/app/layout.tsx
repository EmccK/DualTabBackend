import type { Metadata } from "next";
import "./globals.css";
import { EnvScript } from "@/components/env-script";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "DualTab 管理后台",
  description: "DualTab Chrome 扩展管理后台",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <head>
        <EnvScript />
      </head>
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}

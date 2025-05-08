// next.config.ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  env: {
    NEXT_PUBLIC_APP_DESCRIPTION:
      process.env.NEXT_PUBLIC_APP_DESCRIPTION! || "web-chatify-system",
    NEXT_PUBLIC_API_URL:
      process.env.NEXT_PUBLIC_API_URL! || "http://localhost:80",
    NEXTAUTH_URL: process.env.NEXTAUTH_URL! || "http://localhost:80",
  },
};

export default nextConfig;

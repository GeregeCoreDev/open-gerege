import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Allow images for profile picture (generic pattern or specific domain if known)
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
    ],
  },
  async rewrites() {
    return [
      {
        source: '/api/v1/:path*',
        destination: 'http://localhost:8000/:path*', // Proxy to backend, stripping /api/v1
      },
    ];
  },
};

export default nextConfig;

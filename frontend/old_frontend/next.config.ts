import createNextIntlPlugin from 'next-intl/plugin'
import type { NextConfig } from 'next'

const withNextIntl = createNextIntlPlugin('./i18n/request.ts')

// Security headers configuration
const securityHeaders = [
  {
    key: 'X-DNS-Prefetch-Control',
    value: 'on',
  },
  {
    key: 'Strict-Transport-Security',
    value: 'max-age=63072000; includeSubDomains; preload',
  },
  {
    key: 'X-Frame-Options',
    value: 'SAMEORIGIN',
  },
  {
    key: 'X-Content-Type-Options',
    value: 'nosniff',
  },
  {
    key: 'X-XSS-Protection',
    value: '1; mode=block',
  },
  {
    key: 'Referrer-Policy',
    value: 'strict-origin-when-cross-origin',
  },
  {
    key: 'Permissions-Policy',
    value: 'camera=(), microphone=(self), geolocation=()',
  },
  {
    key: 'Content-Security-Policy',
    value: [
      "default-src 'self'",
      "script-src 'self' 'unsafe-inline' 'unsafe-eval'",
      "style-src 'self' 'unsafe-inline'",
      "img-src 'self' data: https: blob:",
      "font-src 'self' data:",
      "connect-src 'self' https://template.gerege.mn https://sso.gerege.mn https://*.gerege.mn",
      "frame-src 'self' https://sso.gerege.mn",
      "frame-ancestors 'self'",
      "base-uri 'self'",
      "form-action 'self' https://sso.gerege.mn",
    ].join('; '),
  },
]

const nextConfig: NextConfig = {
  // Skip ESLint during build (run separately with `npm run lint`)
  // This avoids the circular structure error with ESLint 9 flat config
  eslint: {
    ignoreDuringBuilds: true,
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'template.gerege.mn',
        pathname: '/api/file/**',
      },
      {
        protocol: 'https',
        hostname: 'picsum.photos',
        pathname: '/1920/400/**',
      },
      {
        protocol: 'https',
        hostname: 'cdn.gerege.mn',
        pathname: '/file/bucket-ap/**',
      },
      {
        protocol: 'https',
        hostname: 'app.gerege.mn',
        pathname: '/api/file/**',
      },
      {
        protocol: 'https',
        hostname: 'business.gerege.mn',
        pathname: '/api/file/**',
      },
    ],
  },
  async headers() {
    return [
      {
        // Apply security headers to all routes
        source: '/:path*',
        headers: securityHeaders,
      },
    ]
  },
}

import { withSentryConfig } from '@sentry/nextjs';

export default withSentryConfig(withNextIntl(nextConfig), {
  // For all available options, see:
  // https://github.com/getsentry/sentry-webpack-plugin#options

  // Suppresses source map uploading logs during build
  silent: true,
  org: 'gerege-systems',
  project: 'frontend-refactor-v25',

  // For all available options, see:
  // https://docs.sentry.io/platforms/javascript/guides/nextjs/manual-setup/

  // Upload a larger set of source maps for prettier stack traces (increases build time)
  widenClientFileUpload: true,

  // Routes browser requests to Sentry through a Next.js rewrite to circumvent ad-blockers (increases server load)
  tunnelRoute: '/monitoring',
});


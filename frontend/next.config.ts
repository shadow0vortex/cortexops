import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // output: 'standalone' produces a minimal self-contained server at .next/standalone/
  // This is required for Docker deployment — without it, the full node_modules
  // directory would need to be present at container runtime.
  // Reference: https://nextjs.org/docs/app/api-reference/config/next-config-js/output
  output: "standalone",

  // Security: Remove the X-Powered-By: Next.js response header to reduce
  // technology fingerprinting surface.
  poweredByHeader: false,

  async redirects() {
    return [
      {
        source: '/docs',
        destination: '/docs/getting-started',
        permanent: true,
      },
      {
        source: '/deploy',
        destination: '/docs/deployment',
        permanent: false,
      },
    ];
  },

  // HTTP Security Headers (TD-024)
  // Applied to every page and API route in the application.
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          // Prevent the page from being embedded in an iframe on a different origin.
          { key: 'X-Frame-Options', value: 'SAMEORIGIN' },
          // Prevent MIME type sniffing.
          { key: 'X-Content-Type-Options', value: 'nosniff' },
          // Control referrer information sent with outbound requests.
          { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
          // Restrict browser features available to this page.
          {
            key: 'Permissions-Policy',
            value: 'camera=(), microphone=(), geolocation=()',
          },
          // Content Security Policy.
          // Note: 'unsafe-inline' is currently required for Tailwind CSS v4 runtime injection.
          // A future hardening pass can introduce nonces once a CSS extraction pipeline is
          // established. Tracked as a follow-up to TD-024.
          {
            key: 'Content-Security-Policy',
            value: [
              "default-src 'self'",
              "script-src 'self' 'unsafe-inline' 'unsafe-eval' https://app.posthog.com",
              "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
              "font-src 'self' https://fonts.gstatic.com",
              "img-src 'self' data: blob:",
              "connect-src 'self' https://app.posthog.com",
              "frame-ancestors 'self'",
            ].join('; '),
          },
        ],
      },
    ];
  },
};

export default nextConfig;

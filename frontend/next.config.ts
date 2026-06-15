import type { NextConfig } from "next";

const nextConfig: NextConfig = {
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
    ]
  },
};

export default nextConfig;

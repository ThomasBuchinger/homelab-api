/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "export",
  images: { unoptimized: true },
  async rewrites() {
    return [
      { source: '/api{/}?', destination: 'http://localhost:8080/api' },
      { source: '/api/component/:slug', destination: 'http://localhost:8080/api/component/:slug' },
      // { source: '/api/componenet/:slug', destination: 'http://localhost:8080/api/componenet/:slug' },
    ]
  },
}

module.exports = nextConfig

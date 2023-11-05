/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "export",
  images: { unoptimized: true },
  async rewrites() {
    return [
      { source: '/api/:slug', destination: 'http://localhost:8080/api/:slug' },
      { source: '/api/public/:slug', destination: 'http://localhost:8080/api/public/:slug' },
      { source: '/api/legacy/:slug', destination: 'http://localhost:8080/api/legacy/:slug' },
      { source: '/api/internal/:slug', destination: 'http://localhost:8080/api/internal/:slug' },
    ]
  },
}

module.exports = nextConfig

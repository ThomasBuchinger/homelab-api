/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "export",
  images: { unoptimized: true },
  async rewrites() {
    return [
      { source: '/api{/}?', destination: 'http://localhost:8080/api' },
      { source: '/api/public/bff/:slug', destination: 'http://localhost:8080/api/public/bff/:slug' },
      { source: '/api/private/bff/syncthing/:slug', destination: 'http://localhost:8080/api/private/bff/syncthing/:slug' },
      // { source: '/api/componenet/:slug', destination: 'http://localhost:8080/api/componenet/:slug' },
    ]
  },
}

module.exports = nextConfig

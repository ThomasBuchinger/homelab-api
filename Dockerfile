FROM cgr.dev/chainguard/node:latest AS builder

ADD ./package.json ./package-lock.json ./
RUN npm ci

Add ./public ./public
Add ./app ./app
Add ./.eslintrc.json ./next.config.js  ./README.md ./tsconfig.json ./
run npm run build

FROM cgr.dev/chainguard/nginx:latest AS app
COPY --from=builder /app/out/ /usr/share/html/

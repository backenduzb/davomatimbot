FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
COPY package*.json ./
RUN npm install --production
COPY --from=builder /app/.svelte-kit/output ./.svelte-kit/output
COPY --from=builder /app/static ./static
EXPOSE 5173
CMD ["npx", "vite", "preview", "--port", "5173", "--host"]
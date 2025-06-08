# Use Bun base image
FROM oven/bun:latest

# Set working directory
WORKDIR /app

# Copy package files and install dependencies
COPY bun.lockb package.json ./
RUN bun install

# Copy the rest of the application
COPY . .

# Pre-generate route tree (TanStack Router CLI)
RUN bunx @tanstack/router-cli generate
# RUN bunx tanstack-router watch

# Expose Vite's default dev server port
EXPOSE 3000

# Start both router watcher and Vite
CMD ["sh", "-c", "bunx @tanstack/router-cli watch & bun run start --host"]
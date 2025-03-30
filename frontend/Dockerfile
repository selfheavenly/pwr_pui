# Use an image with Bun installed
FROM oven/bun:latest
WORKDIR /app

# Copy package manifest and lock files to install dependencies
COPY package.json bun.lockb ./
RUN bun install

# Copy the rest of your source code
COPY . .

# Expose the port your Bun server will listen on
EXPOSE 3000

# Start the Bun server (change "start" to your actual start script if needed)
CMD ["bun", "run", "start"]

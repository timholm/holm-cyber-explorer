FROM node:20-alpine

WORKDIR /app

# Copy package files and install dependencies
COPY package.json ./
RUN npm install --omit=dev

# Copy application code
COPY server.js import.js ./
COPY public/ ./public/

# The html/ directory and manifest.json come from docs-framework
# and must be mounted or baked in at build time
# For now, they're expected to be in the build context
COPY html/ ./html/
COPY manifest.json ./

# Import documents then start server
CMD ["sh", "-c", "node import.js && node server.js"]

EXPOSE 3000

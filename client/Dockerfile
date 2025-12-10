# ------------------------------------------------------
# 1. Build Stage Client
# ------------------------------------------------------
FROM node:24-alpine AS vite_builder

WORKDIR /app/client

RUN echo "Node version: $(node -v)"
RUN echo "npm version: $(npm -v)"

# Cache dependencies
COPY package.json package-lock.json ./

# Install npm
RUN npm install

# Copy source files
COPY . .

# Build and install serve
RUN npm run build

# Nginx serve stage
FROM nginx:alpine

COPY --from=vite_builder /app/client/dist /usr/share/nginx/html
COPY --from=vite_builder /app/client/nginx_default.conf /etc/nginx/conf.d/default.conf

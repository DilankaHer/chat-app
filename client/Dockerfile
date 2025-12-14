# ------------------------------------------------------
# 1. Build Stage Client
# ------------------------------------------------------
FROM node:24-alpine AS vite_builder

WORKDIR /app

ARG VITE_MODE

RUN echo "Node version: $(node -v)"
RUN echo "npm version: $(npm -v)"
RUN echo "VITE_MODE: $VITE_MODE"

# Cache dependencies
COPY package.json package-lock.json ./

# Install npm
RUN npm install

# Copy source files
COPY . .

# Build
RUN npm run build -- --mode $VITE_MODE

# Nginx serve stage
FROM nginx:alpine

COPY --from=vite_builder /app/dist /usr/share/nginx/html
COPY --from=vite_builder /app/nginx_default.conf /etc/nginx/conf.d/default.conf

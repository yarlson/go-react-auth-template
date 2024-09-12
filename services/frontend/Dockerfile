ARG NODE_VERSION=20.11.1

FROM node:${NODE_VERSION}-alpine AS base
WORKDIR /app

FROM base AS pruner
COPY . .
RUN npx turbo prune --scope=frontend --docker

FROM base AS installer
COPY --from=pruner /app/out/json/ .
COPY --from=pruner /app/out/package-lock.json ./package-lock.json
RUN npm ci

FROM installer AS builder
COPY --from=pruner /app/out/full/ .
RUN npm run build

FROM nginx:alpine AS final
COPY --from=builder /app/services/frontend/dist/ /usr/share/nginx/html
COPY services/frontend/config/nginx/nginx.conf /etc/nginx/templates/default.conf.template

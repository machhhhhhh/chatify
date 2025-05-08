#Define argument variables
ARG NEXT_PUBLIC_APP_NAME
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_APP_DESCRIPTION

#Install dependencies 
FROM node:18-alpine as deps
# Define working directory
WORKDIR /tmp
# Provides cached layer for node_modules
# ADD .npmrc /tmp/.npmrc
COPY package.json ./
RUN yarn install
# RUN mkdir -p ${APP_ROOT} && cp -a /tmp/node_modules ${APP_ROOT}/

# Rebuild the source code 
FROM node:18-alpine as builder

# ENV APP_ROOT /usr/src/app

# Define working directory
WORKDIR /usr/src/app
COPY --from=deps /tmp/node_modules ./node_modules
COPY . .

ENV NEXT_PUBLIC_APP_NAME=$NEXT_PUBLIC_APP_NAME
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_APP_DESCRIPTION=$NEXT_PUBLIC_APP_DESCRIPTION
# ENV JWT_SECRET=$JWT_SECRET
# ENV SESSION_CSRF_SECRET=$SESSION_CSRF_SECRET
ENV NEXT_TELEMETRY_DISABLED 1

# RUN yarn cache clean --all
# RUN yarn
# Bundle app source
RUN yarn build

# Production image, copy all the files and run next
FROM node:18-alpine AS runner
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 -D nextjs
RUN chown -R nextjs:nodejs /usr/src/app

COPY --from=builder --chown=nextjs:nodejs /usr/src/app/*.config.js ./
COPY --from=builder --chown=nextjs:nodejs /usr/src/app/.next ./.next
COPY --from=builder --chown=nextjs:nodejs /usr/src/app/node_modules ./node_modules
COPY --from=builder --chown=nextjs:nodejs /usr/src/app/package.json ./package.json

# Set environment variables
ENV HOST 0.0.0.0
ENV PORT 3000
ENV NODE_OPTIONS --max-old-space-size=4096
ENV NEXT_TELEMETRY_DISABLED 1

# Update & Set timezone
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

USER nextjs

EXPOSE 3000

RUN npx next telemetry disable

CMD [ "yarn", "start" ]
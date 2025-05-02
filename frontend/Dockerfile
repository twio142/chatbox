# production
FROM node:22-alpine AS prod

# RUN apk add --no-cache python3 make g++

# ENV PYTHON=/usr/bin/python3

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY package*.json ./
RUN chown -R app:app /app

USER app

RUN npm install --omit=dev

COPY --chown=app:app . .

RUN npm run build

# dev
FROM node:22-alpine AS dev

# RUN apk add --no-cache python3 make g++

# ENV PYTHON=/usr/bin/python3

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY package*.json ./
RUN chown -R app:app /app

USER app

RUN npm install

COPY --chown=app:app . .

EXPOSE 3000

CMD ["npm", "run", "start"]

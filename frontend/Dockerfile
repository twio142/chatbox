# production
FROM node:22-alpine AS prod

# RUN apk add --no-cache python3 make g++

# ENV PYTHON=/usr/bin/python3

WORKDIR /app

COPY package*.json ./

RUN npm install --only=production

COPY . .

RUN npm run build

# dev
FROM node:22-alpine AS dev

# RUN apk add --no-cache python3 make g++

# ENV PYTHON=/usr/bin/python3

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

EXPOSE 3000

CMD ["npm", "run", "start"]

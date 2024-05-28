FROM node:18-alpine

WORKDIR /app

COPY . .

RUN npm install --only=prod

ENTRYPOINT ["node", "main.js"]
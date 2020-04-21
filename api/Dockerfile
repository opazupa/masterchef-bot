FROM node:10.15.0-alpine
WORKDIR /app
COPY package*.json /app/
RUN npm ci
COPY ./ /app/

RUN npm run build

ENV NODE_ENV=production

EXPOSE 3000

CMD ["npm", "start"]
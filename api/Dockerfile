FROM node:10.15.0-alpine
WORKDIR /app
COPY package*.json /app/
RUN npm ci --production
COPY ./ /app/

RUN npm run build

ENV NODE_ENV=production

CMD ["npm", "start"]
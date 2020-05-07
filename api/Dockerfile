FROM node:10.15.0-alpine
WORKDIR /app
COPY package*.json /app/
RUN npm ci --production
RUN npm i typescript
COPY ./ /app/

RUN npm run build
RUN npm run tslint

ENV NODE_ENV=production

CMD ["npm", "start"]
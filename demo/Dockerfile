FROM node:14-alpine3.10

ENV PORT 8181

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install

COPY src ./src

EXPOSE ${PORT}
CMD [ "npm", "start" ]
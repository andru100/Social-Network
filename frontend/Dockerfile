FROM node:14.17.4-alpine

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY . .

RUN npm install --silent

RUN npm install react-scripts@4.0.0 -g --silent

EXPOSE 3000

CMD ["npm", "start"]



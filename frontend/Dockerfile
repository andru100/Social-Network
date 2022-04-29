FROM node:14.17.4-alpine as build 

ARG DOMAIN

ARG PORT

WORKDIR /app 

COPY ./package.json . 
COPY ./package-lock.json . 

RUN npm install 

COPY ./public ./public/ 
COPY ./src ./src/ 

RUN REACT_APP_BACKEND_ADDRESS=${DOMAIN} \
  REACT_APP_BACKEND_PORT=${PORT} \ 
  npm run build 

FROM nginx:1.17-alpine 

COPY --from=build /app/build /usr/share/nginx/html 

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
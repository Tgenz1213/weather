FROM node:22 AS build
WORKDIR /app
COPY . .
RUN yarn install && yarn build --base=/app/dist/

FROM node:22 AS serve
WORKDIR /usr/share/nginx/html
RUN apt-get update && apt-get install -y nginx gettext-base
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf.template /etc/nginx/nginx.conf.template
RUN envsubst '${UPSTREAM_SERVER} ${SERVER_NAME}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

WORKDIR /
COPY package.json .
EXPOSE 80
CMD [ "sh", "-c", "yarn start_${MODE}" ]

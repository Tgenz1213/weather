# Step 1 - build app
FROM node:22 AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Step 2 - serve app with nginx
FROM nginx:1.26
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf.template /etc/nginx/nginx.conf.template
COPY start-nginx.sh /usr/local/bin/start-nginx.sh
RUN chmod +x /usr/local/bin/start-nginx.sh

EXPOSE 80
ENTRYPOINT [ "bin/sh", "/usr/local/bin/start-nginx.sh" ]

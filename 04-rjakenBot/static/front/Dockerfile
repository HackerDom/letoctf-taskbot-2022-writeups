FROM node:lts-slim as builder

WORKDIR /app

COPY package.json ./package.json

RUN npm install --only=prod

COPY . .

RUN npm run build

##############################################################

FROM nginx:1.16.0-alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80


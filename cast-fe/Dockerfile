FROM node:lts-alpine as builder
WORKDIR /build
COPY package.json yarn.lock ./
RUN yarn install
COPY . .
RUN yarn build --mode production

FROM nginx:stable-alpine as production
COPY --from=builder /build/dist /usr/share/nginx/html
COPY --from=builder /build/nginx.conf /etc/nginx/conf.d
RUN rm /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
# build front-end
FROM node:lts-alpine AS frontend
RUN npm config set registry https://registry.npmmirror.com
RUN npm install pnpm -g
COPY . ./app
WORKDIR /app
RUN pnpm bootstrap
RUN pnpm run build-only


FROM nginx:latest AS web 
COPY --from=frontend /app/dist/ /usr/share/nginx/html/webb/
ENTRYPOINT [ "/docker-entrypoint.sh" ]
CMD [ "nginx","-g","daemon off;" ]
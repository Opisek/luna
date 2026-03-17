# Deployment
Since Luna is not officially released yet, use the project with caution! Make frequent backups of your calendars. Assume that you have to nuke the database every time you update until official 1.0.0 release, because there may still be breaking changes ahead!

Keep in mind that Luna is provided with absolutely no warranty or liability from the authors.

## Docker
Currently, no first-party docker images are available. Instead, you can generate and run the images simply by typing `make` in the root directory of this repository.
Make sure you have **make** and **docker** installed.

Until docker images are generated officially, you can also use community-compiled images for the [frontend](https://hub.docker.com/r/tiritibambix/lunafrontend) and the [backend](https://hub.docker.com/r/tiritibambix/lunabackend). These images are provided with no warranty or liability from the main author.

A sample `docker-compose.yml` file is provided in the root directory of the project and below:
```yaml
name: luna
services:
  luna-frontend:
    container_name: luna-frontend
    ports:
      - "8080:8080"
    volumes:
      - /etc/timezone:/etc/timezone:ro
      - /etc/localtime:/etc/localtime:ro
    environment:
      PORT: 8080
      PUBLIC_URL: http://cal.example.com
      API_URL: http://luna-backend:3000
    build:
      context: frontend
      dockerfile: Dockerfile

  luna-backend:
    container_name: luna-backend
    volumes:
      - /srv/luna/data:/data
      - /etc/timezone:/etc/timezone:ro
      - /etc/localtime:/etc/localtime:ro
    environment:
      PUBLIC_URL: http://cal.example.com
      DB_HOST: luna-postgres
      DB_PORT: 5432
      DB_USERNAME: luna
      DB_PASSWORD: luna
      DB_DATABASE: luna
    depends_on:
      - luna-postgres
    build:
      context: backend
      dockerfile: Dockerfile

  luna-postgres:
    image: postgres:16-alpine
    container_name: luna-postgres
    volumes:
      - /srv/luna/postgres:/var/lib/postgresql/data
      - /etc/timezone:/etc/timezone:ro
      - /etc/localtime:/etc/localtime:ro
    environment:
      POSTGRES_USER: luna
      POSTGRES_PASSWORD: luna
      POSTGRES_DB: luna
```

## Baremetal
For baremetal deployment, you must ensure your system has:
- **make**
- **bun** (v1.2.5 or higher)
- **go** (go1.23 or higher)
- a running **postgres** (version 16 or higher) database

For the backend, create an `.env` file in the `backend/src` directory inside the repository and fill it out accordingly to `.env.example`. To start the backend in development mode, run `make` inside the `backend` directory.

If you want to build the backend instead, run `make build` inside the `backend` directory. The resulting binary is created in called `luna-backend` and is generated inside the `backend/src` directory.

For the frontend, create an `.env` file in the `frontend` directory inside the repository and fill it out accordingly to `.env.example`. To start the frontend in development mode, run `make` inside the `frontend` directory.

If you want to build the frontend instead, run `make build` inside the `frontend` directory. To start the compiled frontend, run `bun run ./build/index.js`

## Reverse Proxy
Make sure to put Luna behind a reverse proxy with configured TLS.

Here is a sample configuration file for nginx:
```nginx
server {
  listen 443 ssl;
  server_name cal.example.com;
  include /etc/nginx/certificates/cert.conf;
  client_max_body_size 15M;

  location / {
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header Host $http_host;
    proxy_pass http://127.0.0.1:8080;
  }
}
```

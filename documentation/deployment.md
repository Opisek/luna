# Deployment
Since Luna is not ready to be used yet, this only serves as instructions on how to get Luna up and running for development purposes!

Keep in mind that Luna is provided with absolutely no warranty or liability from the authors.

## Docker
Currently, no first-party docker images are available. Instead, you can generate and run the images simply by typing `make` in the root directory of this repository.
Make sure you have **make** and **docker** installed.

Until docker images are generated officially, you can also use community-compiled images for the [frontend](https://hub.docker.com/r/tiritibambix/lunafrontend) and the [backend](https://hub.docker.com/r/tiritibambix/lunabackend). These images are provided with no warranty or liability from the main author.

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
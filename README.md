# Plex (& co.) Monitoring
Monitoring stack for Plex &amp; related services.

# Setup
This project uses a yaml file to configure the services that it will monitor. An example file is provided in `config.example.yml`. You can copy this file to `config.yml` and edit it to your liking. This contains connection information for the database, info for the server, and Discord bot configuration.

In addition to this configuration, you need to setup a secret store. The application currently supports environment variables (`SECRET_KEY` and `AES_KEY`) to setup the JWT secret and the AES encryption secret. You can also use AWS Secrets Manager, which is the recommended way to do this in production. The application will look for the same secrets.

# Running
To create binaries to run for your platform, run `make`. To create a docker image, run `make build-docker`.

# Features

# Supported Services
- [Plex](https://plex.tv)
- [Sonarr](https://sonarr.tv/)
- [Radarr](https://radarr.video/)
- [Ombi](https://ombi.io/)

## Upcoming(?) Eventually(?) Supported Services
- [Requestrr](https://github.com/darkalfx/requestrr) -- alternatively, bake Discord bot into this application?
- [Deluge](https://deluge-torrent.org/) - should be doable via bash scripts with [Plugin/Execute](https://dev.deluge-torrent.org/wiki/Plugins/Execute)
- [Transmission](https://transmissionbt.com/) ???

# Docker
Docker is the ideal medium for deploying this application. There is a `docker-compose.example.yml` file that outlines one way to setup these containers. You can use them with an existing compose file for the rest of the services and run them all in the same Docker network.

# Architecture
Currently the system works entirely based on webhooks. Each supported service has webhooks built-in, that can be configured with this application. Then, the data gets stored into a Mongo database. This application exposes a dashboard that can be used to visualize the data that is being stored.
# Plex (& co.) Monitoring
Monitoring stack for Plex &amp; related services.

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

# Environment Variables
The `DATABASE_URL` and `SECRET_KEY` environment variables must be set. The database URL is just the typical connection string for a Mongo database in Golang, which looks like `mongodb://user:pass@127.0.0.1:27017`. Secret key should be set specifically for your environment.

# Docker
Docker is the ideal medium for deploying this application. There is a `docker-compose.example.yml` file that outlines one way to setup these containers. You can use them with an existing compose file for the rest of the services and run them all in the same Docker network.

# Architecture
Currently the system works entirely based on webhooks. Each supported service has webhooks built-in, that can be configured with this application. Then, the data gets stored into a Mongo database. This application exposes a dashboard that can be used to visualize the data that is being stored.
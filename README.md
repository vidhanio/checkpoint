# Woodlands Checkpoint

## Setup

1. Add [Woodlands Checkpoint](https://discord.com/api/oauth2/authorize?client_id=896067278393712651&permissions=402653184&scope=bot%20applications.commands) to your server (requires manage roles and manage nicknames permissions)
2. Create to role be given to verified users (e.g. `@Verified`)
3. Create a channel for non-verified users to verify themselves in (e.g. `#verification`)
4. Set up permissons for the role so that only the verified users can see the normal channels
5. Set up the verification channel so only non-verified users can see it (verified users cannot see it)
6. Create roles for grades 7-12
7. Use `/initialize` with the roles you made earlier
8. Woodlands Checkpoint should be set up! 😄

### Fixes to Try

- In both channel settings and role settings, make sure you are allowed to use application commands
- Make sure that the `@Woodlands Checkpoint` role is higher than your verified role
- The bot will not nickname you if your highest role is higher than the `@Woodlands Checkpoint` role

## Self-Hosting

1. Make copies of `students.example.json`, `guilds.example.json`, and `example.env`
2. Remove the `.example` from each of the filenames
3. Fill `students.json` with student information*
4. Put your Discord bot token in the `.env`
5. Run `go build main`
6. Run `./main` (or `./main.exe` for Windows users)

\*DM me on Discord ([`vidhan#0001`](https://discord.com/users/277507281652940800)) if you are interested in doing this step yourself.

## Docker

This repository has been dockerized to allow for running the bot in a Docker container which is portable across different hosts and compatible with kubernetes clusters with a containerd runtime.

To build the Docker image: `make docker-build`
To publish the Docker image to GHCR: `make publish`

### Docker Compose

There is a premade `docker-compose.yml` file for quick deployment to any Docker host. Simply run `docker-compose up` to run the compose file or visit Docker's documentation for more options of the command.

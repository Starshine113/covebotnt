# CoveBot

CoveBot(n't) is a general purpose custom-built Discord bot.
It currently has:
- a starboard
- notes
- a couple of user commands
- mod warnings
- a basic gatekeeper

The starboard, gatekeeper, and user commands are stable, the other parts haven't been tested thoroughly.

## Requirements

- a working [Go](https://golang.org/) environment (only tested with the latest, 1.15)
- a PostgreSQL database (tested down to 9.6)

## Installation

1. Clone the repository
2. Build the bot with `go build`
3. Copy `config.toml.sample` to `config.toml`, and fill in the fields
4. Run the executable

Alternatively, you can use `docker-compose`:
1. Copy `config.toml.sample` to `config.toml`, and fill in the fields
2. `docker-compose build`
3. `docker-compose up -d`

## License

CoveBot is licensed under the GNU Affero General Public License, version 3 or later.
# Passman
> Password and secret manager

This is a toy implementation of a LastPass-like password and secret manager. DO NOT USE THIS IN PRODUCTION.

## Installation
1. Install [Go](https://go.dev/doc/install)
1. Install [templ](https://templ.guide/quick-start/installation)
1. Install the [TailwindCSS](https://tailwindcss.com/blog/standalone-cli) standalone CLI
1. Install Make
1. Clone the project and install dependencies
```sh
git clone git@github.com:oalexander6/passman.git
cd ./passman
go mod download
```

## Notes
### Secrets
Secrets must be placed in the `./secrets` folder. The required files must be created containing the desired values:
- `JWT_SECRET`
- `SESSION_SECRET`
- `STORAGE_ENCRYPTION_KEY`
- `STORAGE_PASSWORD`
- `POSTGRES_ADMIN_PASSWORD`

### Docker Setup
1. Run `docker-compose up`
2. Use `ifconfig` and find the ipv4 for the interface `docker0`
3. Go to `localhost:8080`, select `postgres`, set host to `<IPV4>:5432`, set username to postgres, set password to contents of POSTGRES_ADMIN_PASSWORD secret
4. Create a new database named `authelia`
5. Create a new user with `CREATE USER authelia WITH PASSWORD 'STORAGE_PASSWORD';`
6. Grant new user full access to authelia database with `GRANT ALL ON SCHEMA public TO authelia;`
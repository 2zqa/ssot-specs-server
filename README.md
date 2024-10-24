# SSOT specs server

An API for [ssot-specs-collector](https://github.com/2zqa/ssot-specs-collector) clients to send specifications to, and for the [ssot-specs-dashboard](https://github.com/2zqa/ssot-specs-dashboard) to fetch and view the specifications from.

For the related API client and OpenAPI documentation, see [ssot-specs-api-client](https://github.com/2zqa/ssot-specs-api-client).

> [!NOTE]
> This project is part of a suite of projects that work together. For all other related projects, see this search query: [`owner:2zqa topic:ssot`](https://github.com/search?q=owner%3A2zqa+topic%3Assot&type=repositories)

## Getting started

### Prerequisites

- Go (tested with `go version go1.20.3 linux/amd64`)
- PostgreSQL (tested with `psql (PostgreSQL) 14.7 (Ubuntu 14.7-0ubuntu0.22.04.1)`)

### Setup locally

1. Install Go: https://go.dev/doc/install
2. Install and setup postgres: https://www.postgresql.org/docs/current/tutorial-start.html. You need to create a database named `ssotdb` with a user named `ssot`. You may provide a password yourself.
3. Create a GitLab application with the `openid` scope and a redirect to `http://localhost:3000/callback`. Ensure "Confidential" is unchecked. You can create the application on https://\<yourgitlabdomain\>/-/profile/applications
4. Copy the `.env.example` file to `.env` and fill in the values. The API key can be anything, and the DSN is in the following format: `postgres://<username>:<password>@<host>:<port>/<database>`. For example: `postgres://ssot:password@localhost/ssotdb`

### Installation and running

1. Clone and enter the repository: `git clone https://github.com/2zqa/ssot-specs-server.git && cd ssot-specs-server`
2. Run `go install ./cmd/api`
3. Run `api`

## License

SSOT specifications server is licensed under the [MIT](LICENSE) license.

## Acknowledgements

- [Voys](https://www.voys.nl/) for facilitating the internship where this project was developed
- Alex Edwards for providing an excellent resource on building an API with Go, _[Let's Go Further!](https://lets-go-further.alexedwards.net)_, upon which this project is based

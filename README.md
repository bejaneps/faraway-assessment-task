# Faraway Assessment Task

Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW Challenge. 

## Requirements

- Docker
- Go

## Getting Started

1. Create a `.env` file with `.env.example` content.

2. Run application

```
make run
```

## Makefile

To run unit tests

```
make test-unit
```

To run integration tests **without Docker**

```
make test-integration
```

To run integration tests **with Docker**

```
make docker-test-integration
```

To see test coverage

```
make test-coverage
```

## Commit guidelines

- Setup git hooks by running `make hook`.
- Commit messages are validated in by linter.
- Github PR should have a title following conventional commit guideline.
  - fix: increasing k8s termination deadline
  - build: caching docker image build stage
  - docs: updating commit guidelines
- All commits should follow [conventional commits guideline](https://www.conventionalcommits.org/en/v1.0.0/)
- After a new tag is created and pushed to Github, new release will be created from that tag

## Linter

golangci-lint **v1.53.3** is used as primary linter. To use it run:

```
make deps
make lint
```

## Algorithm

For this project I decided to use simple sha256 checksum algorithm. In simple words, server generates a random number in [0 ... 2_000_000] range, calculates it's sha256 hashsum and sends to client. Client reads the hashsum and starts searching for a number in range [0 ... 2_000_000] until hashes are equal. If they are equal, client sends found number to server, server checks validity and sends random quote back.

Quotes are stored statically in database and cache, the reason that they are not fetched from external service(API) is because almost all quotes services are paid and require authentication with token or username/password :() . In future if such service would be enhanced, then a separate endpoint could be added for adding new quotes to database and cache.

## Project Structure

A project is divided by client and server domains.

* Transport layer(internal/transport) is uppermost layer which is responsible for dealing with transport components
* Service layer(internal/service) is middle layer which is responsible for application's business logic
* Repository layet(internal/repository) is lowermost layer which is responsible for dealing with database or cache

Each layer is separated in a way that it doesn't know anything about it's lower layer except exposed interfaces, this way we ensure that upper layers don't depend on implementations of lower layers.

Seeds are used to pre-populate cache and database so we have initial data to work with and test, they are located at internal/pkg/quotes/quotes.json . They are run in cmd/migrations

Each service is built locally and then copied to debian image, this is faster approach for local testings

## TODO

1. Add github action for running all types of tests
2. Implement e2e tests ovh/venom
3. Add github action for release
4. Add github action for linter
5. Add tests for transport layer
6. Add unit tests for db package
7. Add github action for building and pushing image to registry
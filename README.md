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
$ make run
```

API should be available on http://localhost:4500.

## E2E Tests

ovh/venom package is used for e2e testing, it's easy to use and configure. All test cases are located in `test/` folder, test cases are written in .yml format.

To run e2e tests locally do:

```
make e2e-test-build
make e2e-test
```

To run e2e tests in Github pipeline:

```
Open pr
Set pr to draft and then ready for review
Check actions, e2e tests should be running
```

E2e tests are also triggered when pr is merged to main

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

We are using golangci-lint **v1.53.3** as primary linter. To use it run:

```
make deps
make lint
```

It's integrated into CI, so if your code fails on linter it won't pass CI.

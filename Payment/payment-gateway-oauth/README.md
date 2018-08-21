# Payment Gateway Oauth2 Server

## Prerequisites

* Pipenv
* Python 3.5 or higher

## Development Environment Setup

After clone this repo. Init environment by:

```
pipenv install
```

Config `.env`. Using the same database with Payment Gateway server.

Finally, run this project.

```
pipenv shell
make dev_run
```

## Temp Usage

* `/auth/token`: Endpoint for request Password Grant
* `/verify`: Endpoint for checking Oauth2 token valid
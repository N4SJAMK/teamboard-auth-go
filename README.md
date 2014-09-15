# Teamboard Authentication

```
teamboard-auth
```

## Dependencies

Requires `golang` and `MongoDB`.

## Installation

Install with:
```
go get github.com/N4SJAMK/teamboard-auth
```

Run with:
```
teamboard-auth
```
You can provide `HOST` and `PORT` environmental variables:
```
HOST=every.day PORT=4200 teamboard-auth
```
You can provide the following `MongoDB` config as env. vars:
- `MONGODB_URL` defaults to `mongodb://localhost`
- `MONGODB_NAME` defaults to `teamboard-dev-go`




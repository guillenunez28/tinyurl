# tinyurl

Convert long urls to short urls. This project is built in Golang using Gin.

## Setup
Software needed on the system
- Golang 1.21
- Docker Desktop

Have Golang version 1.21 installed. If not, the Golang software can be found
[here](https://go.dev/dl/).

Docker Desktop can be installed [here](https://www.docker.com/products/docker-desktop/).

## Run the project
From the root directory, run the following actions.

Setup MariaDB. 
```
docker run --name mariadbtest -e MYSQL_ROOT_PASSWORD=mypass -p 3306:3306 -d docker.io/library/mariadb:10.3
```

Run the short URL program.
```
go run cmd/tinyurl/main.go
```

During the first time bootup, the DB schema is executed. If the DB and tables already exists, then no
DB operations happen.


## API endpoints
> GET /

Returns the Tiny Url resources. Limited to 10

> GET /:short_url

Returns the tiny url resource with short name `short_url`. If the resource exists,
the endpoint will redirect the user to the long url. If no resource is found, then
the response states the resource was not found. 

> GET /:short_url/stats

Return the tiny url resource stats (last 24 hour, last week and all time hits) for
the resource with ID `short_url`. If no record is found, the response is {}

> POST /

Creates a tiny url resource in the DB. In the payload, the key `long_version` must be
specified. The key `expiration_date` is optional. The expected format is `YYYY-MM-DD`.
Other formats of the `expiration_date` will fail.

> DELETE /:short_url

Deletes the tiny url resource in the DB with short url `short_url`. If it's not found,
the response states the resource was not found.

## Limitations
- The short URL has a length of 6 characters and may contain only numbers.
- The long URL is assumed to start with `www.` and be on HTTPS. If it's not,
the redirect might not work as intended.
- Deleting a short URL will not erase the hits history.
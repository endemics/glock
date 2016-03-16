# Glock

A simple lock HTTP API written in GOLANG.

Once upon a time I needed a simple mechanism to protect access to a shared resource. All actors being trusted, I didn't need anything fancy. I didn't need high availability, performance and consistency. I needed something simple to set, so glock was born.

## API

*Glock* uses a very simple RESTFUL HTTP API:

- `GET /`: retrieve list of locks
- `GET /:id`: retrieve info for lock `id`, 404 if it does not exists
- `PUT /:id`: tries to acquire lock `id`. If not existing, creates it and returns 201. If already existing, returns 409
- `DELETE /:id`: deletes lock `id`. Removes the lock `id` and returns 200

### examples

You can manipulate the API using curl to return the HTTP codes.

- list all the locks:

```
curl -s -o /dev/null -w "%{http_code}" -X GET http://lockserver/
```

- list infos for a lock:
```
curl -s -o /dev/null -w "%{http_code}" -X GET http://lockserver/123
```

- create a lock

```
curl -s -o /dev/null -w "%{http_code}" -X PUT http://lockserver/123
```

- release a lock
```
curl -s -o /dev/null -w "%{http_code}" -X DELETE http://lockserver/123
```

## Runtime options

There is only one option that can be used when running *glock*:

- `-d` will activate debugging

## Compilation

You will need golang installed and the sources, then:

```
go build -a
```

This should produce a `glock` binary for your platform that you can then install where you want.

However, if you want to install the `glock` binary in a lightweight container, you will need to do a static compilation using:

```
CGO_ENABLED=0 go build -a -installsuffix cgo
```

## Installation

The recommanded installation method is in a docker container. *glock* when started will listen on port `8080`.

For this, build a static version of glock using the instructions above, then build the container using:

```
docker build -t glock .
```

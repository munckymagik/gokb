# Map synchronisation

Problem: you need to present a map type to be read concurrently. But you also need to periodically
trash and overwrite the entire map.

## Running

Run the server like:

```shell
cd mapsync
go run -race ./cmd implname
```

Where `implname` is one of: `naive`, `atomic`, `mutex`, `rwmutex` or `chan`

The db always contains a key call `x`. It will be incremented upon each update of the map.

To read the value of `x` use curl:

```shell
curl http://localhost:8070/get?name=x
```

If there are any races (as there are in `naive`) this will trigger a warning in the server log output.

## Load testing

Run the server as (note no `-race` so we don't carry the overhead):

```shell
go run ./cmd implname
```

Then use `ab` or `wrk` or `vegeta` to pummel it with requests.

```shell
ab -k -t30 -n 2000000 -c 200 "127.0.0.1:8070/get?name=x"
```

| impl      | rps (mean) cold | rps (mean) warm |
| --------- | ---------------: | -------------: |
| `atomic`  |         71692.12 |       81712.02 |
| `mutex`   |         85683.40 |       81509.31 |
| `rwmutex` |         80091.04 |       80014.63 |
| `chan`    |         77620.46 |       83610.24 |

So once my laptop was warmed up they all performed roughly the same.
It could be that the limiting factor is not the repo implementation but
the http server.

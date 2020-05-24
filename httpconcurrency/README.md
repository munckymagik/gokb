Based on the tutorial at: https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/

The following benchmarks were interesting:

```shell
$ ab -n 20000 -c 200 "127.0.0.1:8070/inc?name=i"
$ ab -k -n 20000 -c 200 "127.0.0.1:8070/inc?name=i"
$ ab -k -n 20000 -c 16000 "127.0.0.1:8070/inc?name=i"
```

The first example crapped out at around 16383. The -k keepalive option resolved this. The problem
is expected to be ab reaching the port limit.

For an explanation see: https://stackoverflow.com/a/30357879

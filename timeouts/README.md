Making an HTTP request using netcat

```bash
nc -c localhost 8888
```

Then paste the contents of `request.txt`. Don't redirect into STDIN, this
seems to result in immediate closing of the connection.


# build

```sh
docker build -t twi-meteor .
```

# run

```sh
docker run -d  --env-file .env twi-meteor
```

# check if it works

```sh
docker exec <hash> tail -f /var/log/app.log
```


とりあえずtweet取ってくるところまでやった。
# TODO
- filter の作成
- 

# build

```sh
docker build -t twi-meteor .
```

# run

```sh
docker run -d twi-meteor --env-file .env
```

# check if it works

```sh
docker exec <hash> tail -f /var/log/backup.log
```

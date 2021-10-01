
# build

```sh
docker build -t twi-meteor .
```

# run

```sh
docker run -d twi-meteor 
```

# check if it works

```sh
docker exec 34841d8f4459 tail -f /var/log/backup.log
```

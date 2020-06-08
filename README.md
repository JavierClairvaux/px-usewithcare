# px-usewithcare

### Pull image
```
docker pull javier1/px-usewithcare:0.0.1
```

### Run container
```
docker run --rm  -p 8080:8080 --name test --oom-kill-disable --cpus=2 javier1/px-usewithcare:0.0.1
```

### Get memEater status
```
curl localhost:8080/memeater
```

### Run memEater
```
curl localhost:8080/memeater/start/<mem in MB>
```

For example
```
curl localhost:8080/memeater/start/100
```

### Free memory
```
curl localhost:8080/memeater/free
```

To increase or decrease memory usage it is necessary to free memory and start memEater again since there can only be one memEater process at the time.

### Get cpuBurner status
```
curl localhost:8080/cpuburner
```

### Start cpuBurner
```
curl -X PUT localhost:8080/cpuburner/start
```

### Stop cpuBurner
```
curl localhost:8080/cpuburner/stop
```

There can only be one cpuBurner process at the time.

You can find the container [here.](https://hub.docker.com/repository/docker/javier1/px-usewithcare)

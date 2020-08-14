# px-usewithcare

### Pull image
```
docker pull javier1/px-usewithcare:0.0.3
```

### Run container
```
docker run --rm  -p 8080:8080 --name test --oom-kill-disable --cpus=2 javier1/px-usewithcare:0.0.3
```

### Start a memEater
```
$ curl -X POST -d "{\"memMb\": 100}" localhost:8080/memeater
{"mem_mb":100,"id":"2286d6fd-9ce5-424f-b7ba-62ce9b22495b"}
```
Where memMb is the load of memory in MBs.

### Get a memEater status
```
$ curl localhost:8080/memeater/<ID from the first stop>
```

For example:
```
$ curl localhost:8080/memeater/ed8b4c7e-8338-4792-b2b3-de9a5c2177e5
{"mem_mb":100,"id":"ed8b4c7e-8338-4792-b2b3-de9a5c2177e5"}
```

### Stop a memEater
```
$ curl -X DELETE localhost:8080/memeater/<ID from the first steop>
```

For example:
```
$ curl -X DELETE localhost:8080/memeater/ed8b4c7e-8338-4792-b2b3-de9a5c2177e5
```

### Get a list of active memEaters
```
curl localhost:8080/memeaters
{"IDs":["ed8b4c7e-8338-4792-b2b3-de9a5c2177e5","2286d6fd-9ce5-424f-b7ba-62ce9b22495b"]}
```

### Start a cpuBurner

```
$ curl -X POST -d "{\"num_burn\": 1, \"ttl\": 500000}" localhost:8080/cpuburner
{"running":true,"num_burn":1,"ttl":500000,"id":"80c3f0e4-25f5-4bbe-93ab-4f86a7f371a2"}
```
Where num_burn is the number of cores you want to burn and ttl is the Time To Live in miliseconds.

### Get a cpuBurners status

```
$ curl localhost:8080/cpuburner/<ID from the first step>
```

For example:
```
$ curl localhost:8080/cpuburner/80c3f0e4-25f5-4bbe-93ab-4f86a7f371a2
{"running":true,"num_burn":1,"ttl":500000,"id":"80c3f0e4-25f5-4bbe-93ab-4f86a7f371a2"}
```

### Stop cpuBurner

```
$ curl -X DELETE localhost:8080/cpuburner/<ID from the first step>
```

For example:
```
$ curl -X DELETE localhost:8080/cpuburner/80c3f0e4-25f5-4bbe-93ab-4f86a7f371a2
```

### List all existing cpuBurners
```
$ curl localhost:8080/cpuburners
{"IDs":["ee9d9ab9-6563-45ba-823a-330827403b90","b04c0db8-dbe4-4bd5-ac4b-270cc3f6007e"]}
```
You can find the container [here.](https://hub.docker.com/repository/docker/javier1/px-usewithcare)

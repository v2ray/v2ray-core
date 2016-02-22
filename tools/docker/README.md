Docker build for v2ray
=======================

Usage
-----

To build the image:

```bash
./build.sh
```

Then spin up a v2ray instance with:

```bash
./run.sh
```

The docker image will generate a server config file with random user id on first run.
You can get see it with:

```bash
docker logs v2ray
```

To tail the access log, run:

```bash
docker exec v2ray tail -F /v2ray/logs/access.log
```

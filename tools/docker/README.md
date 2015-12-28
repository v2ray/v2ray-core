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

The build script will generate a server config file with random user id. You
can get it from `server-cfg.json`.

To tail the access log, run:

```bash
docker exec v2ray tail -F /v2ray/logs/access.log
```

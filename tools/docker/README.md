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

You can also specify config file by manual with:

```bash
docker run -d --name=v2ray -p 27183:27183 -v /config/file.json:/go/server-config.json $USER/v2ray
```

To tail the access log, run:

```bash
docker exec v2ray tail -F /go/access.log
```

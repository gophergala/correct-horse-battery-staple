Working Demo
------------

http://where-you-at.com

![iOS ScreenShot](https://s3.amazonaws.com/f.cl.ly/items/0i3l1w1I3V023Z3o341Y/Image%202015-01-25%20at%204.39.23%20PM.png)

Installing
----------

You need to have Go installed.

```bash
go get -u github.com/gophergala/correct-horse-battery-staple
go get -u -d -tags=js github.com/gophergala/correct-horse-battery-staple/...
```

Running
-------

In the root project folder:

```bash
go build -o main && ./main
```

Deploy
------

```bash
go test ./... && GOOS=linux go build -o main && ./deploy.sh
```

Notes
-----

Features whiteboard: http://whiteboardfox.com/29046-5977-1449

Mockup
------

![](https://s3.amazonaws.com/f.cl.ly/items/0j2B2K2Y3T2u1O1m0h2y/Image%202015-01-23%20at%206.19.12%20PM.png)

# short-link
create short link

use fiber, redis

### setup redis

```
docker pull redis
docker run --name redis-lab -p 6379:6379 -d redis
```

### run main
``` 
go run main.go
```

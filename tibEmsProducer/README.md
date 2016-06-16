Build with 
```bash
go build -ldflags -extldflags=-static goTibemsMsgProducer.go
```

As there is an issue with relative paths, change ldflags to fit your directory structure to ems libraries

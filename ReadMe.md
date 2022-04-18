# docker-registry-gui
This projects generates an HTML frontend for a Docker registry which lacks this functionality.

## Build 
To build this application, golang is required.  
To install GoLang, view [the docs](https://golang.google.cn/doc/install)  

### GoLang
Build the application  
```bash
$ go build ./cmd/docker-registry-gui
```

### Docker
Build the docker image (will also build the binary)
```bash
$ docker build . -t docker-registry-gui
```

## Run
To run the application, either a precompiled binary is required, a Docker image or GoLang to build/run.  
By default, the application will not start as it needs a docker registry to connect to.  
If not configured, the server will listen under `0.0.0.0:8080`


### Binary
```bash
# Lists all flags
$ ./docker-registry-gui -help

# Provide configuration file
$ ./docker-registry-gui -config my_config.yaml

# Provide flags
# Without AUTH
$ ./docker-registry-gui -registry-name docker.example.com
# With AUTH 
$ ./docker-registry-gui -registry-name docker.example.com -registry-username user -registry-password password
```

### Docker
After creating the Docker image, it will not start without proper configuration.  
```bash
# Mount a configuration file and use that
$ docker run -v "$(pwd)/local_configuration.yaml:configuration.yaml" -p 8080:8080 docker-registry-gui -config configuration.yaml
```

## Configuration
### Options
#### YAML
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  gracefulTimeoutMs: 15000

docker:
  protocol: https://
  registry: registry.example.com
  port: 0 # 0 means protocol default is used
  username: user
  password: pass
```


### Order
Configuration is loaded in the following order  
1. [Default configuration file](./default_configuration.yaml)
2. Configuration file provided via flag `-config`
3. Configuration provided via flags
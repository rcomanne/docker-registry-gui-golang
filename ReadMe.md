# docker-registry-gui
This projects generates an HTML frontend for a Docker registry which lacks this functionality.

## Build 
To build this application, golang is required.  
To install GoLang, view [the docs](https://golang.google.cn/doc/install)  

Build the application  
```bash
$ go build .
```

## Run
To run the application, either a precompiled binary is required, a Docker image or GoLang to build/run  

```bash
# precompiled binary
# Lists all flags
$ ./docker-registry-gui -help

# Provide configuration file
$ ./docker-registry-gui -config my_config.yaml

# Provide flags
$ ./docker-registry-gui -registry-name docker.example.com -registry-username user -registry-password password
```

### Configuration order
Configuration is loaded in the following order  
1. [Default configuration file](./default_configuration.yaml)
2. Configuration file provided via flag `-config`
3. Configuration provided via flags
# Forum Authentication   
   
The goal of the project is to implement the alternative ways of authorising users. Two methods have been presented for user: Google and Github.  
  
## Architecture   
The architecture of the project is based on the gateway and two microservices: auth and app, which communicate HTTP and REST API.

## Launch   
Two ways are provied to launch the server:  
### 1. Simple: 
open each service in their root directory in different terminals and run each of them separately
### 2. Docker compose:
```
docker compose up
```
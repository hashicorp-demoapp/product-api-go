# product-api-go
Products API written in Go

[![CircleCI](https://circleci.com/gh/hashicorp-demoapp/product-api-go.svg?style=svg)](https://circleci.com/gh/hashicorp-demoapp/product-api-go)  

Docker Image: [https://hub.docker.com/repository/docker/hashicorpdemoapp/product-api](https://hub.docker.com/repository/docker/hashicorpdemoapp/product-api)


## Running the API

To test the API you can use the Shipyard blueprint in this repository.

To install shipyard run:

```
curl https://shipyard.run/install | bash
```

Then run the API:

```
➜ shipyard destroy
2020-02-21T11:33:54.004Z [INFO]  Destroy Ingress: ref=db
2020-02-21T11:33:54.004Z [INFO]  Destroy Ingress: ref=api
2020-02-21T11:33:54.459Z [INFO]  Destroy Container: ref=db
2020-02-21T11:33:54.489Z [INFO]  Destroy Container: ref=api
2020-02-21T11:33:55.051Z [INFO]  Destroy Network: ref=local

product-api-go on  remotes/origin/HEAD [!?] via 🐹 v1.13.5 on 🐳 v19.03.6 () took 2s 
➜ curl -v localhost:19090/health 

product-api-go on  remotes/origin/HEAD [!?] via 🐹 v1.13.5 on 🐳 v19.03.6 () 
➜ shipyard run ./blueprint
Running configuration from:  ./blueprint

2020-02-21T11:34:02.153Z [DEBUG] Statefile does not exist
2020-02-21T11:34:02.153Z [INFO]  Creating Network: ref=local

# ...

########################################################

Title Product API
Author Erik Veld, Nicholas Jackson

# ...
```

The API is available at `http://localhost:19090`

```
➜ curl -v localhost:19090/coffees
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 19090 (#0)
> GET /coffees HTTP/1.1
> Host: localhost:19090
> User-Agent: curl/7.58.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Date: Fri, 21 Feb 2020 11:31:48 GMT
< Content-Length: 1165
< Content-Type: text/plain; charset=utf-8
< 
[{"id":1,"name":"Packer Spiced Latte","teaser":"Packed with goodness to spice up your images","description":"","price":350,"image":"/packer.png","ingredients":[{"ingredient_id":1},{"ingredient_id":2},{"ingredient_id":4}]},{"id":2,"name":"Vaulatte","teaser":"Nothing gives you a safe and secure feeling like a Vaulatte","description":"","price":200,"image":"/vault.png","ingredients":[{"ingredient_id":1},{"ingredient_id":2}]},{"id":3,"name":"Nomadicano","teaser":"Drink one today and you will want to schedule another","description":"","price":150,"image":"/nomad.png","ingredients":[{"ingredient_id":1},{"ingredient_id":3}]},{"id":4,"name":"Terraspresso","teaser":"Nothing kickstarts your day like a provision of Terraspresso","description":"","price":150,"image":"/terraform.png","ingredients":[{"ingredient_id":1}]},{"id":5,"name":"Vagrante espresso","teaser":"Stdin is not a tty","description":"","price":200,"image":"/vagrant.png","ingredients":[{"ingredient_id":1}]},{"id":6,"name":"Connectaccino","teaser":"Discover t* Connection #0 to host localhost left intact
he wonders of our meshy service","description":"","price":250,"image":"/consul.png","ingredients":[{"ingredient_id":1},{"ingredient_id":5}]}]%   
```

## Endpoints

Some notes on select API endpoints:
| Endpoint | Description |
| --- | --- |
| '/health' | (DEPRECATED) Health check endpoint that verifies DB connectivity. This has been replaced by `/health/readyz` |
| '/health/livez' | Health check endpoint that verifies the server has started. |
| '/health/readyz' | Health check endpoint that verifies the server is connected to the DB and ready to serve requests. |

## Requesting changes / Governance
This API is shared by multiple teams and therefore we require some form of process to ensure new features or changes do not break functionality
relied on by others. To make changes to the API:

* Create an issue defining the changes your would like to make (e.g. https://github.com/hashicorp-demoapp/product-api-go/issues/1)
* Ensure that at least 1 other maintainer has approved your changes before starting to code
* Implement your change
* Ensure that you have created a functional test which refers to your Issue (e.g. https://github.com/hashicorp-demoapp/product-api-go/blob/tests/functional_tests/features/basic_functionality.feature). Functional tests ensure that changes by others do not change the behaviour of the API which may break your use case.
* Submit a PR


## Release process

Pushing a new tag to Github will trigger a release on CircleCI

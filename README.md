Welups-bridge: a bridge from other blockchains to Welups
========================================================

## Deployment
//TODO

## Development
### Prequisites
  * Go >= 1.18
  * [pressly/goose](https://github.com/pressly/goose): a database migration tool for Golang
  * [cap'n proto](https://capnproto.org/): gRPC-like data serialization and rpc framework
  * RMDB: postgres
  * Cache: redis
### Structure of the project

The project is structured as a monorepo to simplify integration and versioning of the
whole set of microservices in the system, which could grow quite large as new bridges
to/from different chains are implemented
* **micros/**: the code for the system's microservices, each in its own directory
  * **core/**: the core service, implements API gateway, monitoring and administration
    functionalities
  * **<bridge>/**: implements the corresponding bridging logic. Currently there's only
    **weleth/** for Welups<->Ethereum bridge.
* **common/**: models, consts and configs common for all microservices
* **libs/**: library for some commonly used functionalities
* **service-managers/**: wrapping low level interactions with external services/daemons,
  e.g. redis, concrete database, mailing service, web services etc...

Each microservice in **micros/** is further structured as follows:
* **config/**: configuration logic
* **model/**: datatypes modelling the business domain
* **dao/**: abstract interface to interact with persistent store of data, which is
  currently just postgres over sqlx interface. Mostly dumb CRUD operations for the
  respective domain
* **blogic/**: business logic for the respective domain. The innermost layer, does most of
  the heavy-lifting computation so that outer layer can stay dumb.
* **http/** or **rpc/**: the service-providing layer. Also stays as dumb and simple as
  possible, mostly just deserializes data from request, call **blogic** to get results,
  then serializes results into response.
* **migrations/**: database migrations SQL scripts.

### Migrate DB
  * Start a postgres cluster in your preferred way
  * For each microservice directory:
    ```
      cd migrations
      goose postgres "host=<addr> port=<port> user=<username> password=<password> \
      dbname=welbridge sslmode=disable" up
    ```
### Build
Install the prequisites above, then run the following command in the root of the project
```sh
  ./build.sh build
```

To cleanup build artifacts:
```sh
  ./build.sh cleanup
```

### Run
* Each microservice reads config from either (preferably) environment variables or from
  .env file in the same directory as the binary. Variable names and default values are
  defined in the **config/config.go** file of each microservice.
### Test
```sh
  go test
```

# code_challenge1

## Setup

The easiest way to setup the project is to use docker compose.

Just copy the [docker-compose.yml](./docker-compose.yml) file provided, copy it somewhere. open the directory in the terminal and execuate

```bash
docker compose up -d
```

Then you can access the service through the API endpoint http://localhost:8080

## Test API Use Postman

Import the file [postman_collection.json](./postman_collection.json) in Postman, it provided all the all the API collection of this project.

Change the parameter, test the APIs as you like.

## Code Structure

As you can see, the code of this project is divided into three parts. 

First, the logging package, the project use uber/zap to do logging. It has a simpler interface and a better performance according to [this benchmark](https://github.com/uber-go/zap#Performance). The log package of this provide a simple wrapper of uber/zap, make it much easier to use logging function in other modules.

Second is the database package. This project leverage the default database/sql package, use raw sql, provided some simple method needed by other packages.

Last is the server package. This is where the bussiness logic lives. It use golang [gin-gonic/gin](https://github.com/gin-gonic/gin) to expose all the server APIs.

Of course there is also the main package, which init everything and runs the program.

I borrow a lot code and ideas from [my another opensource project](https://github.com/simon-ding/polaris). If you look closer, you will see a lot similarities.

## Money Handling

Because float type is not precise, we cannot use it to represent money in real word. But because of money has a difined number of decimals, we can simply use int to represent it. i.e. we can use 1001 to represent 10.01, can converted it back on API returns

## Github Actions

There three github actions available. They will run on every push or pull requests, run all the tests, apply the lint rules and build the output docker image. Then if user want to run the project simply need to download the docker image and use the [docker-compose.yml](./docker-compose.yml) provided.

## Build You Self

This project requires go 1.23.1 and above. Clone this project and simply run 

```
go build ./
```

will download all required dependence and compile all the packages.


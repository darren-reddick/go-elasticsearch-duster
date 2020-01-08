# go-elasticsearch-duster
Repository for Go app to cleanup Elasticsearch indexes



The app runs once and pulls the following config in from the json config file:

* domain : the elasticsearch domain:port
* patterns : an array of {name, age} objects. Name is the index pattern, Age is the max age in days to retain the index.
* username : the username to authenticate to the domain. This is an optional field
* password : the password to authenticate to the domain. This is an optional field

## Getting Started

Pull the project to an appropriate location under **GOPATH**

    git pull https://github.com/devopsgoat/go-elasticsearch-duster.git ${GOPATH}/src/devops.goat/go-elasticsearch-duster

### Prerequisites

* To build and install the binary you will need to have Go installed https://golang.org/doc/install
* If creating a docker container, docker will need to be installed

### Installing

To install the binary

```
go install
```

Create a configuration file in the directory with format - example:

```
cat << EOF > my-config.json
{
    "domain" : "elastic-search-domain:9243",
    "username" : "elastic",
    "password" : "password",
    "patterns" : [
        {
            "name" : "twitter",
            "age" : "2"
        },
        {
            "name" : "metricbeat-host-7.2.0",
            "age" : "5"
        }
    ]
}
EOF
```

Run the app as a dry-run using the config file created - this will be informational and will not delete indexes

```
go-elasticsearch-duster -c my-config.json
```

Do a full cleanup run if happy with output from dry-run - this will delete indexes

```
go-elasticsearch-duster -c my-config.json -d
```



## Running the tests

To run the tests

    go test


## Docker

A docker container can be created from this project using the Dockerfile supplied - example:

```
docker build -t my-go-elasticsearch-duster .
```

This is designed to be a very lightweight container built from scratch with only the go binary.

The container can be run mapping in a config file as a volume:

```
docker run -v $(pwd)/es-cloud-config.json:/es-cloud-config.json es_client -c es-cloud-config.json -d
```

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Darren Reddick** - *Initial work* - [PurpleBooth](https://github.com/devopsgoat)


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments





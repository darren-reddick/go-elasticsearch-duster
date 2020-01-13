# go-elasticsearch-duster

## Summary

Repository for Go app to cleanup Elasticsearch indexes

## Technical Overview

The app will run a single cycle to query indexes in an Elasticsearch cluster and purge them based on their index pattern and date derived from the index name - this currently support indexes with naming convention **[index name]-YYYY-MM-DD** (this may change in the future). Indexes used internally by Elasticsearch that start with '.' are also excluded.

The config is pulled from a json config file and supports the following settings:

* domain : the elasticsearch domain:port (required)
* patterns : an array of {name, age} objects. Name is the index pattern, Age is the max age in days to retain the index.
* username : the username to authenticate to the domain. This is an optional field
* password : the password to authenticate to the domain. This is an optional field

An example config file can be seen in the [Installing](Installing) section

As the app runs a single cycle before exiting it is designed to be run in a Kubernetes CronJob or similar.

### go-elasticsearch-duster Arguments
| Arg | Description | Default |
| ------- | ------- | ------------|
| `-c` | The location of the config file | config.json  |
| `-d` | Whether to delete or not. Default is false which is essentially a dry-run | false  |


## Getting Started

Clone the project to an appropriate location under **GOPATH**

    git clone https://github.com/devopsgoat/go-elasticsearch-duster.git ${GOPATH}/src/devops.goat/go-elasticsearch-duster

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

### SSL

If connecting to an SSL endpoint the cert chain for the endpoint will need to be included in the container build for the http client to trust the endpoint. The cert for AWS Elasticsearch in eu-west-1 is included in this repo and the Dockerfile gives an example of how to include this in the correct container location.

Alternatively the cert file could be mapped in as volume under /etc/ssl/certs/

### Kubernetes

An example Kubernetes manifest for running as a Cronjob and mapping in the config as a ConfigMap:

```
---
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: es-cleanup
spec:
  concurrencyPolicy: Forbid # dont run a new job if the previous one is already running
  schedule: "00 01 * * *" # run at 1am nightly
  jobTemplate:
    metadata:
      name: es-cleanup
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
          - image: my-go-elasticsearch-duster-image
            name: es-cleanup
            resources:
              limits:
                memory: "256Mi"
            args: ["-c","config.json"]
            volumeMounts:
            - name: config
              mountPath: /config.json
              subPath: config.json
          volumes:
          - name: config
            configMap:
              name: es-cleanup
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: es-cleanup
data:
  config.json: |
{
    "domain" : "xxxxxxxxxxxxxx.eu-west-1.aws.found.io:9243",
    "username" : "elastic",
    "password" : "password1",
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





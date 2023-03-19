# Readme
1. Run Elasticsearch and Kibana in local
```shell
$ docker network create elastic
$ docker run -d --name elasticsearch --net elastic -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.17.9
$ docker run --name kibana --net elastic -p 127.0.0.1:5601:5601 -e "ELASTICSEARCH_HOSTS=http://elasticsearch:9200" -d docker.elastic.co/kibana/kibana:7.17.9
```

2. go install and execute main.go
```shell
$ go install
$ go run main.go
```
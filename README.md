# InfluxDB-VK-API

InfluxDB-VK-API is a Go-based application that communicates with VK API, retrieves information about a certain VK group and then writes that information into InfluxDB database.

## Features

- VK API Interface
- Retrieval of Subscriber Information from a VK Group
- Writing Information to InfluxDB

## Docker

We provide a Dockerfile for creating a Docker image for this application. You can build the image with the command:

```shell
docker build -t InfluxDB-VK-API .
```
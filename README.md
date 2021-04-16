# user-api-gokit

### Building the binary

```
# Windows

export GOOS=linux // without this go compiles to the current operating system
go build main.go

```

### Building & pushing the docker image

```
docker build -t udaykiranr/user-api:<version> .

docker push udaykiranr/user-api:<version>

Version should be a proper semver 

vMAJOR.MINOR.PATCH

ex: v1.0.0
```

### Nginx config for load-balancing

```
pid        /tmp/nginx.pid;

events {
    worker_connections  1024;
}

http {
  upstream user-api {
    server 127.0.0.1:8500;
    server 127.0.0.1:8501;
    server 127.0.0.1:8502;
    server 127.0.0.1:8503;
  }

  server {
    listen 80;
    server_name bindutest.rcplatform.io;
    location / {
      proxy_pass http://user-api;
    }
  }
}

```

### Running the nginx container

```
docker run -d --name nginx -v $PWD/nginx.conf:/etc/nginx/nginx.conf --network=host nginx

Where $PWD is the path where nginx.conf is present with the contents above
```

### Running etcd DB in the container

```
docker run -d -v /root/data:/root/data --network=host --name etcd quay.io/coreos/etcd:v3.2 /usr/local/bin/etcd --data-dir /root/data --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://<VM-IP>:2379

```

### Running user-api in the container

```
docker run -d -e ETCD_IP=<ETCD-VM-IP> -p 8502:8500 --name api2 udaykiranr/user-api:<version> 
```
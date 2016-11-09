`graphviz` is an HTTP server which calls
[GraphViz](http://www.graphviz.org/) to visualize a specified .dot
file.

## Build and Run from Source Code

if we already got GraphViz and [Go](https://golang.org/) installed, we
can run this program as follows

```
go get github.com/k8sp/graphviz
$GOPATH/bin/graphviz -addr=:9090
```

Then we can direct our Web browser to

```
http://localhost:9090/?dot=https://gist.githubusercontent.com/wangkuiyi/c4e0015211dd1b9bde2e20455a6cd38e/raw/4d5ec099f98a5f326cf6f108bcf510cadba1a0b4/ci-arch.dot
```

so to visualize the
[Gist file `ci-arch.dot`](https://gist.github.com/wangkuiyi/c4e0015211dd1b9bde2e20455a6cd38e)
in the Web browser window.

## Build and Run the Docker Image

If we have Docker installed, we can build `graphviz` into a Docker
image without installing GraphViz and Go:

```
go get github.com/k8sp/graphviz
cd $GOPATH/src/github.com/k8sp/graphviz
docker build -t graphviz .
```

and run it as a Docker container:

```
docker run -p 9090:9090 -v /tmp:/cache graphviz
```

Now we can direct our Web browser to above link.

## Run with Images on DockerHub.com

We can also run the Docker images built from stable releases by
DockerHub.com:

```
docker run -p 9090:9090 -v /tmp:/cache k8sp/graphviz
```

<!--  LocalWords:  graphviz GraphViz GOPATH addr ci cd
 -->

# Go blind Receptor

Main usage is to put up an HTTPServer that will notify out the hits, with certain prefix to not get fucked by chinese scanners

# Install

```bash
go get -u github.com/jcatala/gbrx
```

# Usage

```bash
$ gbrx
  -notify
        Notify the incoming request via telegram bot (Not recommended to listen to root directory)
  -port int
        Specify another port, default: 9080 (default 9080)
  -prefix string
        Receive and notify just the requests with some certain prefix.
  -rbody string
        Custom response body, default: UNIX TIMESTAMP
  -redirect string
        To make the server redirect somewhere
  -verbose
        To be verbose
```


## Use Case

To get a simple response server

```bash
# With a curl, you'll get a unix timestamp
gbrx -verbose

# Making the server to redirect with a 302
gbrx -verbose -redirect "https://f4d3.io"

# Making the server to redirect and NOTIFY the request via telegram
gbrx -verbose -redirect "https://f4d3.io" -notify

# Making the server to redirect and notify, but just respond to a request having some PREFIX
# I totally recommend this, to get the scanners out of the notify zone, lol
gbrx -verbose -redirect "https://f4d3.io" -notify -prefix "veryhiddenprefix"

```

## TODO

* Add TLS/SSL support.
* Find an elegant way to read from the socket.
* A :beer:

**NOTE: To make the notify flag works, you must to have working the [Go Quick Message](https://github.com/jcatala/gqm), its just a config file :) **

# Thanks for checking this out!
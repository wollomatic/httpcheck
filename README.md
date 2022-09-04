# httpcheck

A simple quick and dirty go application which checks if a list of websites are up or down.

## Getting Started

Put the binary file to your preferred location and create a yaml file with the list of websites you want to check.

```yaml
---
delay: 100
service:
  - name: Human friendly name
    url: https://example.com/health
    test: GET
    status: 200
    text: OK
    timeout: 1000
    retries: 0
    err_delay: 100
  - name: Second example
    url: http://example.com
    test: HEAD
    status: 301
  - name: minimal example
    url: https://stiftung-musica-sacra.de/
```

`delay` is the delay between each check in milliseconds.

`name` is the human friendly name of the service.

`url` is the url to check.

`test` is the http method to use. Default is `GET`.

`status` is the expected http status code. Default is `200`.

`text` is the expected text in the response body. If there is none, no checks on the body are done.

`timeout` is the timeout for the http request in milliseconds. Default is `1000`.

`retries` is the number of retries if the http request fails. Default is `0`.

`err_delay` is the delay between retries in milliseconds. Default is `100`.


## Usage

Usage: httpcheck file.yaml


Return codes:

0 - all services are ok

1 - 1 service is not ok or no filename is given or problem in yaml file

n - n services are not ok

### Examples

```
$ httpcheck services.yaml
[2022-09-04 17:02:08] starting service checks
+ Example redir                    HEAD   HTTP/2.0   301 Moved Permanently               141ms     0 retries   nginx
+ Whoami                           GET    HTTP/2.0   200 OK                              186ms     0 retries   nginx        your IP address and some more stuff
+ Site with basic auth             GET    HTTP/2.0   403 Forbidden                       128ms     0 retries   nginx        permission
+ srv.example.com                  GET    HTTP/2.0   418 I'm a teapot                     68ms     0 retries
+ some html site                   GET    HTTP/2.0   200 OK                              160ms     0 retries                <!DOCTYPE html>
---
No problems detected.

$ httpcheck someservicewithtimeout.yaml
[2022-09-04 17:20:01] starting service checks
- AlmNet redir                     Head "https://almnet.de/whoami": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
---
Unhealthy services: 1
```

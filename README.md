- [Run applications on remote computers from the web](#run-applications-on-remote-computers-from-the-web)
  * [Warning](#warning)
  * [Usage](#usage)
    + [POST run](#post-run)
    + [GET history](#get-history)
    + [POST history](#post-history)
    + [Web UI](#web-ui)

# Run applications on remote computers from the web

The go-runner is a super simple web service which allows you to run any command on the computer it is started.

## Warning

**This program opens your computer to anyone. Make sure to use it only in trusted environments or add some sort of security infront of it. I'm not responsible for any damage you or others do to your computer.**

## Usage

```help
Usage of ./go-runner:
  -host string
      Hostname to listen on (default "127.0.0.1")
  -port int
      Port to start the web server on (default 8080)
  -www string
      Directory path where you keep your HTML files for the web UI. If no path is specified only the API is provided.
```

The server exposes three API endpoints and optionally with the `-www` option a web UI.

### POST run

Send a POST to `/v2/run/[command]`, where [command] is then executed on the server PC.

In the args form data you can specify additional arguments you'd like to pass to the command.

If the program is not in the global path you can specify the path to the file with the `cwd` form value.

Example:

```bash
$ curl -X POST http://127.0.0.1:8080/v2/run/ls -d 'args=-a -p'
./
../
.git/
.gitignore
.travis.yml
.vscode/
LICENSE
README.md
go-runner
go-runner.go

# The cwd parameter needs an absolute path and a trailing slash.
$ curl -X POST http://127.0.0.1:8080/v2/run/my-program -d 'args=--help' -d "cwd=$PWD/my-stuff/"
my-program help message
```

You can also get the response in JSON format if you desire. Specify the query param `json=true` to do so.

Example:

```bash
$ curl -X POST http://127.0.0.1:8080/v2/run/ls?json=true -d 'args=-a -p'
{"output":"./\n../\n.git/\n.gitignore\n.travis.yml\n.vscode/\nLICENSE\nREADME.md\ngo-runner\ngo-runner.go\nwww/\n"}
```

### GET history

Simple GET reques to `/v2/history` results in a bash like history with all the previous commands.

Example:

```bash
$ curl http://127.0.0.1:8080/v2/history
0 ls -l
1 ls -a
2 ls -a -d
3 ls -a -p
```

You can also get the response in JSON format if you desire. Specify the query param `json=true` to do so.

Example:

```bash
$ curl http://127.0.0.1:8080/v2/history?json=true

[{"id":0,"cmd":"ls","args":"-l"},{"id":1,"cmd":"ls","args":"-a"},{"id":2,"cmd":"ls","args":"-a -d"},{"id":3,"cmd":"ls","args":"-a -p"}]
```

### POST history

Similar to the GET you can send a POST to `/v2/history/[id]` where [id] matches the number in front of the command you want to execute from the [GET history](#get-history).

Example:

```bash
$ curl -X POST http://127.0.0.1:8080/v2/history/3
./
../
.git/
.gitignore
.travis.yml
.vscode/
LICENSE
README.md
go-runner
go-runner.go
````

You can also get the response in JSON format if you desire. Specify the query param `json=true` to do so.

Example:

```bash
$ curl -X POST http://127.0.0.1:8080/v2/history/3?json=true
{"output":"./\n../\n.git/\n.gitignore\n.travis.yml\n.vscode/\nLICENSE\nREADME.md\ngo-runner\ngo-runner.go\nwww/\n"}
```

### Web UI

Specify the `-www` option with a path to where you store your static web files.

An example ui can be found in the repository under the www folder.

Example:

```bash
$ ./go-runner -www ./www/
```

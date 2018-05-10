# Run applications on remote computers from the web

The go-runner is a super simple web service which allows you to run any command on the computer it is started.

**This program opens your computer to anyone. Make sure to use it only in trusted environments or add some sort of security infront of it. I'm not responsible for any damage you or others do to your computer.**

## Usage

```bash
Usage of ./go-runner:
  -host string
        Hostname to listen on (default "127.0.0.1")
  -port int
        Port to start the web server on (default 8080)
```

The server exposes three API endpoints.

### POST run

Send a POST to `/v1/run/[command]`, where [command] is then executed on the server PC.

In the args data you can specify additional arguments you'd like to pass to the command.

Example:

```bash
$ curl -X POST http://127.0.0.1:8080/v1/run/ls -d 'args=-a -p'
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
```

### GET history

Simple GET reques to `/v1/history` results in a bash like history with all the previous commands.

Example:

```bash
$ curl http://127.0.0.1:8080/v1/history
0 ls -l
1 ls -a
2 ls -a -d
3 ls -a -p
```

### POST history

Similar to the GET you can send a POST to `/v1/history/[id]` where [id] matches the number infront of the command you want to execute from the GET history.

Example:

```bash
$ curl -X POST http://127.0.0.1:8080/v1/history/3
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

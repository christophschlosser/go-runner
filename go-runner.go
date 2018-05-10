package main

import (
	"flag"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} -> ${status}\n",
	}))
	e.Use(middleware.Recover())

	const apiPrefix = "/v1"

	// Routes
	e.GET(apiPrefix+"/history", history)
	e.POST(apiPrefix+"/history/:id", runHistory)
	e.POST(apiPrefix+"/run/:command", runCmd)

	// check for CLI options
	host := flag.String("host", "127.0.0.1", "Hostname to listen on")
	port := flag.Int("port", 8080, "Port to start the web server on")

	flag.Parse()

	// Start server
	e.Logger.Fatal(e.Start(*host + ":" + strconv.Itoa(*port)))
}

// History
var (
	cmds []string
	args []string
)

// Handler
func history(c echo.Context) error {
	his := ""
	for id, h := range cmds {
		his += (strconv.Itoa(id) + " " + h + " " + args[id] + "\n")
	}
	return c.String(http.StatusOK, his)
}

func runHistory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if len(cmds) <= id {
		return c.String(http.StatusUnprocessableEntity, "ID does not exist\n")
	}
	cmd := cmds[id]
	arg := args[id]
	out, err := run(cmd, arg)
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, "Command with ID does no longer exist\n")
	}
	return c.String(http.StatusOK, out)
}

func runCmd(c echo.Context) error {
	command := c.Param("command")
	arg := c.FormValue("args")

	out, err := run(command, arg)

	if err != nil {
		return c.String(http.StatusUnprocessableEntity, "Can not run comand: "+command+" "+arg+"\n")
	}
	return c.String(http.StatusOK, out)
}

func run(command string, arg string) (string, error) {
	cmds = append(cmds, command)
	args = append(args, arg)
	cmdArgs := strings.Fields(arg)
	out, err := exec.Command(command, cmdArgs...).Output()
	res := string(out[:])

	return res, err
}

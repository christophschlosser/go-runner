package main

import (
	"flag"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// check for CLI options
	host := flag.String("host", "127.0.0.1", "Hostname to listen on")
	port := flag.Int("port", 8080, "Port to start the web server on")
	html := flag.String("www", "", "Directory path where you keep your HTML files for the web UI. If no path is specified only the API is provided.")

	flag.Parse()
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} -> ${status}\n",
	}))
	e.Use(middleware.Recover())

	const apiPrefix = "/v2"

	// Routes
	if *html != "" {
		e.Static("/", *html)
	}
	e.GET(apiPrefix+"/history", history)
	e.POST(apiPrefix+"/history/:id", runHistory)
	e.POST(apiPrefix+"/run/:command", runCmd)

	// Start server
	e.Logger.Fatal(e.Start(*host + ":" + strconv.Itoa(*port)))
}

type output struct {
	Text string `json:"output"`
}

// History
type historyEntry struct {
	ID   int    `json:"id"`
	Cmd  string `json:"cmd"`
	Args string `json:"args"`
	Dir  string `json:"dir"`
}

var historyEntries []historyEntry

// Handler
func history(c echo.Context) error {
	useJSON := c.QueryParam("json") == "true"
	if useJSON {
		if len(historyEntries) == 0 {
			return c.String(http.StatusOK, "{}")
		}
		return c.JSON(http.StatusOK, historyEntries)
	}
	his := ""
	for _, c := range historyEntries {
		his += strconv.Itoa(c.ID) + " " + c.Dir + c.Cmd + " " + c.Args + "\n"
	}
	return c.String(http.StatusOK, his)
}

func runHistory(c echo.Context) error {
	useJSON := c.QueryParam("json") == "true"
	id, _ := strconv.Atoi(c.Param("id"))
	if len(historyEntries) <= id {
		txt := "ID does not exist\n"
		if useJSON {
			o := output{
				Text: txt,
			}
			return c.JSON(http.StatusUnprocessableEntity, o)
		}
		return c.String(http.StatusUnprocessableEntity, txt)
	}
	cmd := historyEntries[id].Cmd
	arg := historyEntries[id].Args
	dir := historyEntries[id].Dir
	out, err := run(cmd, arg, dir)
	if err != nil {
		txt := "Command (" + dir + cmd + ") with history id " + c.Param("id") + " does not exist\n"
		if useJSON {
			o := output{
				Text: txt,
			}
			return c.JSON(http.StatusUnprocessableEntity, o)
		}
		return c.String(http.StatusUnprocessableEntity, txt)
	}
	if useJSON {
		o := output{
			Text: out,
		}
		return c.JSON(http.StatusOK, o)
	}
	return c.String(http.StatusOK, out)
}

func runCmd(c echo.Context) error {
	useJSON := c.QueryParam("json") == "true"
	command := c.Param("command")
	arg := c.FormValue("args")
	dir := c.FormValue("cwd")

	out, err := run(command, arg, dir)

	if err != nil {
		txt := "Can not run comand: " + dir + command + " " + arg + "\n"
		if useJSON {
			o := output{
				Text: txt,
			}
			return c.JSON(http.StatusUnprocessableEntity, o)
		}
		return c.String(http.StatusUnprocessableEntity, txt)
	}
	if useJSON {
		o := output{
			Text: out,
		}
		return c.JSON(http.StatusOK, o)
	}
	return c.String(http.StatusOK, out)
}

func run(command string, arg string, dir string) (string, error) {
	h := historyEntry{
		ID:   len(historyEntries),
		Cmd:  command,
		Args: arg,
		Dir:  dir,
	}
	historyEntries = append(historyEntries, h)
	cmdArgs := strings.Fields(arg)
	cmd := exec.Command(dir + command, cmdArgs...)
	out, err := cmd.Output()
	res := string(out[:])

	return res, err
}

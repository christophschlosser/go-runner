$(function () {
    $('[data-toggle="tooltip"]').tooltip()
})

function runCmd() {
    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    request.open("POST", self.location.origin + "/v2/run/" + document.getElementById("cmd").value.trim() + "?json=true", true)

    // Refresh history when done
    request.onload = function () {
        document.getElementById("output").innerHTML = "[" + self.location.hostname + "] $> " + document.getElementById("cwd").value.trim() + document.getElementById("cmd").value.trim() + " " + document.getElementById("arg").value.trim()
        var output = ""
        if (request.status == 200) {
            output = JSON.stringify(JSON.parse(this.response).output).replace(/\\n/g, "<br/>").replace(/\"([^(\")"]*)\"/g, "$1")
        } else if (request.status == 422) {
            output = JSON.parse(this.response).output
        }

        document.getElementById("output").innerHTML += "<br/>" + output
        history()
    }

    document.getElementById("output").innerHTML = "[" + self.location.hostname + "] $> " + document.getElementById("cmd").value.trim() + " " + document.getElementById("arg").value.trim()

    var fd = new FormData()
    fd.append("args", document.getElementById("arg").value.trim())
    fd.append("cwd", document.getElementById("cwd").value.trim())
    // Send request
    request.send(fd)
}

function history() {
    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    request.open("GET", self.location.origin + "/v2/history?json=true", true)

    request.onload = function () {
        var data = JSON.parse(this.response);
        var history = document.getElementById("history")
        history.innerHTML = ""
        data.forEach(element => {
            var cmd = document.createElement("div")
            cmd.className = "col-sm-10"
            cmd.id = "history-entry-cmd-" + element.id
            var cmdArgs = document.createElement("kbd")
            cmdArgs.className = "history-cmd"
            cmdArgs.textContent = element.dir + element.cmd + " " + element.args
            cmd.appendChild(cmdArgs)

            var btn = document.createElement("button")
            btn.className = "btn btn-outline-primary col-sm-2"
            btn.setAttribute("onclick", "runHistory(" + element.id + ")")
            btn.textContent = "Run again"

            var entry = document.createElement("div")
            entry.className = "row history-entry"
            entry.id = "history-entry-" + element.id
            entry.appendChild(cmd)
            entry.appendChild(btn)

            history.appendChild(entry)
            history.appendChild(document.createElement("hr"))
        });
    }

    // Send request
    request.send()
}

function runHistory(id) {
    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    request.open("POST", self.location.origin + "/v2/history/" + id + "?json=true", true)

    // Refresh history when done
    request.onload = function () {
        document.getElementById("output").innerHTML = "[" + self.location.hostname + "] $> " + document.getElementById("history-entry-cmd-" + id).innerText
        var output = ""
        if (request.status == 200) {
            output = JSON.stringify(JSON.parse(this.response).output).replace(/\\n/g, "<br/>").replace(/\"([^(\")"]*)\"/g, "$1")
        } else if (request.status == 422) {
            output = JSON.parse(this.response).output
        }

        document.getElementById("output").innerHTML += "<br/>" + output
        history()
    }

    // Send request
    request.send()
}
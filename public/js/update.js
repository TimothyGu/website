window.onload = function() {
    var url = 'ws://localhost:3000/socket';
    var socket = new WebSocket(url);

    socket.onopen = function() {
        socket.send("" + WS_ID);
    };

    var results = document.getElementById("query-results");

    socket.onmessage = function(event) {
        var parts = event.data.split(",");
        var percent = parts[0], files = parts[1], lines = parts[2];

        results.innerHTML = "Indexed <b>" + lines + "</b> lines in <b>" + files + "</b> files (<b>" + percent + "</b>%)";

        document.getElementById("bar").setAttribute("style", "width: " + percent + "%");
        document.getElementById("progress-label").innerText = percent + "%";

        if (percent == "100") {
            var pBar = document.getElementById("progress-bar");
            var cssClasses = pBar.getAttribute("class");
            pBar.setAttribute("class", cssClasses + " success");
        }
    }
}

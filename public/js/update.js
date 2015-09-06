window.onload = function() {
    var socket = new WebSocket('ws://localhost:3000/socket');

    var counter = document.getElementsByClassName("counter")[0];

    socket.onmessage = function(event) {
        var parts = event.data.split(",");
        var percent = parts[0], files = parts[1], lines = parts[2];

        counter.innerText = "Indexed " + lines + " lines in " + files + " files (" + percent + "%)";

        document.getElementById("bar").setAttribute("style", "width: " + percent + "%");
        document.getElementById("progress-label").innerText = percent + "%";
    }
}

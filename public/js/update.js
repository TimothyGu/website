window.onload = function() {
    var url = 'ws://localhost:3000/socket';
    var socket = new WebSocket(url);

    socket.onopen = function() {
        socket.send("" + WS_ID);
    };

    var results = document.getElementById("query-results");

    socket.onmessage = function(event) {
        if (event.data.charAt(0) == '{') {
            populateCode(JSON.parse(event.data));
            return;
        } else if (event.data.charAt(0) == '#') {
            searchResCount(event.data.substring(1));
            return;
        }
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

function populateCode(json) {
    var table = document.createElement("TABLE");
    table.setAttribute("class", "codeview");
    Object.keys(json).forEach(function(k) {
        var line = json[k];

        var tr = document.createElement("TR");
        var td = document.createElement("TD");
        td.innerText = k;
        tr.appendChild(td);
        var td = document.createElement("TD");
        var pre = document.createElement("pre");
        pre.innerText = line;
        td.appendChild(pre);
        tr.appendChild(td);

        table.appendChild(tr);
    });

    document.getElementById("search-results").appendChild(table);
    var br = document.createElement("BR");
    document.getElementById("search-results").appendChild(br);
}

function searchResCount(c) {
    $('#search-res-count').text("Found results in " + c + " files");
}

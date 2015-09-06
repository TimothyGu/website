window.onload = function() {
    var stxt = document.getElementsByClassName("stxt")[0];
    var sbtn = document.getElementsByClassName("sbtn")[0];
    if (stxt != null && sbtn != null) {
        sbtn.onclick = function() {
            window.location = "/search?query=" + encodeURIComponent(stxt.value);
        }

        var expanded = false;
        $(stxt).keydown(function(event) {
            if (event.keyCode == 13) {
                if (!expanded)
                    $('#collapse-search').slideToggle();
                expanded = true;
                event.preventDefault();
            }
        });    
    }
}

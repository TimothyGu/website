window.onload = function() {
    var stxt = document.getElementsByClassName("stxt")[0];
    var sbtn = document.getElementsByClassName("sbtn")[0];

    sbtn.onclick = function() {
        window.location = "/search?query=" + encodeURIComponent(stxt.innerText);
    }
}    
    

window.onload = function() {
    var stxt = document.getElementsByClassName("stxt")[0];
    var sbtn = document.getElementsByClassName("sbtn")[0];

    var expanded = false;
    var input = $('#repo-input');
    input.keydown(function(event) {
        if (event.keyCode == 13 && !expanded) {
            var val = input[0].value;
            if (val.split("/").length != 2) {
                event.preventDefault();
                return;
            }
            console.log("here");
            var ele = $('#repo-search .github.icon');
            ele.animate({
                top: "-=50px"
            });
            var ele2 = $('#repo-search .code.icon');
            ele2.animate({
                top: "-=50px"
            });
            $('#repo-form-value').val(input[0].value);
            input[0].value = "";
            input.attr('placeholder', 'enter semantic query');
            expanded = true;

            input.attr('style', 'font-family: Monaco, monospace; font-size: 24px;');
            event.preventDefault();
        }
        $('#query-form-value').val(input[0].value);
    });    
}

var arr = [[1, 2, 3, 4],
[5, 6, 7, 8],
[9, 10, 11, 12],
[13, 14, 15, 16]];
var step = 0;
var is_start = false;
window.onload = function () {
    $("#start").click(startGame);
    $(".chips").click(move);

}

function startGame(event) {

    var a = new Array();
    for (var i = 1; i < 16; i++)
        a[i] = i;

    for (var i = 1; i < 16; i++) {
        var random_idx = parseInt(Math.random() * 15 + 1);//[1,16)
        var tmp = a[i];
        a[i] = a[random_idx];
        a[random_idx] = tmp;
    }


    current = $("#whole-picture").children();
    for (var i = 0; i < 15; i++) {
        current[i].id = "p" + a[i + 1];
    }

    current[15].id = "p" + 16;
    is_start = true;
    $("#game-status").text("");
    $("#step-board").val(0);
}

function move(event) {
    if (is_start) {
        str = this.id;
        id = parseInt(str.substr(1, str.length));

        this_pos = $(this).index() + 1;
        blank_pos = $("#p16").index() + 1;

        var samerow = Boolean(parseInt((this_pos - 1) / 4) == parseInt((blank_pos - 1) / 4));
        if ((this_pos + 1 == blank_pos && samerow) || (this_pos - 1 == blank_pos &&samerow) || this_pos - 4 == blank_pos || this_pos + 4 == blank_pos) {

            $("#p16").attr("id", this.id);
            this.id = "p16";
            step++;
            $("#step-board").val(step);
        }

        if (isWin()) {
            is_start = false;
            $("#game-status").text("You Win!");
        }

    }
}

function isWin() {
    current = $("#whole-picture").children();
    for (var i = 0; i < current.length; i++) {
        str = current[i].id;
        id = parseInt(str.substr(1, str.length));
        if (id != i + 1)
            return false;
    }
    return true;
}
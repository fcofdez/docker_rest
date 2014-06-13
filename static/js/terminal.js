$(document).ready(function($) {
    $('body').terminal("runDocker", {
        login: false,
        greetings: "Press enter to submit commands",
        onBlur: function() {
            // the height of the body is only 2 lines initialy
            return false;
        }
    });
});

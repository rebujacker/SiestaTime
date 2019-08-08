
$.getScript("./node_modules/xterm/dist/xterm.js");
$.getScript("./static/lib/scripts/local-echo.js");

$(document).ready(function() {


		/*
        var term = new Terminal();
        term.open(document.getElementById('terminal'));
        term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ')
        */

		// Start an xterm.js instance
		const term = new Terminal();
		term.open(document.getElementById('terminal'));

		// Create a local echo controller
		const localEcho = new LocalEchoController(term);

		// Read a single line from the user
		localEcho.read("[bichito]> ")
    		.then(input => alert(`User entered: ${input}`))
    		.catch(error => alert(`Error reading: ${error}`));

});
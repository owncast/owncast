define(
	"main",
	[
		"MessageList"
	],
	function(MessageList) {
		var ws = new WebSocket("ws://" + location.host + "/entry");
		var list = new MessageList(ws);
		ko.applyBindings(list);
	}
);

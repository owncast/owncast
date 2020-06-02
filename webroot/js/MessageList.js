define(
	"MessageList",
	[
		"Message"
	],
	function(Message) {

		function MessageList(ws) {
			var that = this;
			this.messages = ko.observableArray();

			this.editingMessage = ko.observable(new Message());

			this.send = function() {
				var model = this.editingMessage().toModel();
				ws.send($.toJSON(model));
				var message = new Message();
				message.author(model.author);
				message.image("https://robohash.org/" + model.author);
				console.log("https://robohash.org/" + model.author)
				this.editingMessage(message);
			};

			ws.onmessage = function (e) {
				console.log("on message")
				var model = $.evalJSON(e.data);
				var msg = new Message(model);
				that.messages.push(msg);

				// Scroll DIV to the bottom
				scrollSmoothToBottom("messages-container")
			};
		}
		
		return MessageList;
	}
);


function scrollSmoothToBottom(id) {
	var div = document.getElementById(id);
	$('#' + id).animate({
		scrollTop: div.scrollHeight - div.clientHeight
	}, 500);
}
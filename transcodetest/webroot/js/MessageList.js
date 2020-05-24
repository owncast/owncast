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
				this.editingMessage(message);
			};

			ws.onmessage = function(e) {
				var model = $.evalJSON(e.data);
				var msg = new Message(model);
				that.messages.push(msg);
			};
		}
		
		return MessageList;
	}
);

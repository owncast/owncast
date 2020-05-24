define(
	"Message",
	[],
	function() {

		function Message(model) {
			if (model !== undefined) {
				this.author = ko.observable(model.author);
				this.body = ko.observable(model.body);
			} else {
				this.author = ko.observable("Anonymous");
				this.body = ko.observable("");
			}

			this.toModel = function() {
				return {
					author: this.author(),
					body: this.body()
				};
			}
		}

		return Message;
	}
);

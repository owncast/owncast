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
			this.image = ko.observable("https://robohash.org/" + this.author() + "?set=set3&size=50x50")


			this.toModel = function() {
				return {
					author: this.author(),
					body: this.body(),
					image: this.image()
				};
			}
		}

		return Message;
	}
);

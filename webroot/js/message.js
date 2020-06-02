class Message {
	constructor(model) {
		this.author = model.author
		this.body = model.body
		this.image = "https://robohash.org/" + model.author
	}

	linkedText() {
		return autoLink(this.body, { embed: true })
	}
	
	toModel() {
		return {
			author: this.author(),
			body: this.body(),
			image: this.image()
		}
	}
}
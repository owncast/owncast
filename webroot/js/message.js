class Message {
	constructor(model) {
		this.author = model.author
		this.body = model.body
		this.image = "https://robohash.org/" + model.author
		this.id = model.id
	}

	linkedText() {
		return autoLink(this.body, { embed: true })
	}

	toModel() {
		return {
			author: this.author(),
			body: this.body(),
			image: this.image(),
			id: this.id
		}
	}
}


// convert newlines to <br>s

class Messaging {
	constructor() {
		this.chatDisplayed = false;
		this.username = "";
		this.avatarSource = "https://robohash.org/";

		this.messageCharCount = 0;
		this.maxMessageLength = 500;


		this.tagChatToggle = document.querySelector("#chat-toggle");

		this.tagUserInfoDisplay = document.querySelector("#user-info-display");
		this.tagUserInfoChanger = document.querySelector("#user-info-change");

		this.tagUsernameDisplay = document.querySelector("#username-display");
		this.imgUsernameAvatar = document.querySelector("#username-avatar");

		this.tagMessageAuthor = document.querySelector("#self-message-author");

		this.tagMessageFormWarning = document.querySelector("#message-form-warning");
		
		this.tagAppContainer = document.querySelector("#app-container");

		this.inputChangeUserName = document.querySelector("#username-change-input"); 
		this.btnUpdateUserName = document.querySelector("#button-update-username"); 
		this.btnCancelUpdateUsername = document.querySelector("#button-cancel-change"); 

		this.formMessageInput = document.querySelector("#inputBody");

	}
	init() {
		this.tagChatToggle.addEventListener("click", this.handleChatToggle);
		this.tagUsernameDisplay.addEventListener("click", this.handleShowChangeNameForm);
		
		this.btnUpdateUserName.addEventListener("click", this.handleUpdateUsername);
		this.btnCancelUpdateUsername.addEventListener("click", this.handleHideChangeNameForm);
		
		this.inputChangeUserName.addEventListener("keydown", this.handleUsernameKeydown);
		this.formMessageInput.addEventListener("keyup", this.handleMessageInputKeydown);
	}

	handleChatToggle = () => {
		if (this.chatDisplayed) {
			this.tagAppContainer.className = "flex no-chat";
			this.chatDisplayed = false;
		} else {
			this.tagAppContainer.className = "flex";
			this.chatDisplayed = true;
		}
	}
	handleShowChangeNameForm = () => {
		this.tagUserInfoDisplay.style.display = "none";
		this.tagUserInfoChanger.style.display = "flex";
	}
	handleHideChangeNameForm = () => {
		this.tagUserInfoDisplay.style.display = "flex";
		this.tagUserInfoChanger.style.display = "none";
	}
	handleUpdateUsername = () => {
		var newValue = this.inputChangeUserName.value;
		newValue = newValue.trim();
		// do other string cleanup?

		if (newValue) {
			this.userName = newValue;
			this.inputChangeUserName.value = newValue;
			this.tagMessageAuthor.innerText = newValue;
			this.tagUsernameDisplay.innerText = newValue;
			this.imgUsernameAvatar.src = this.avatarSource + newValue;
		}
		this.handleHideChangeNameForm();
	}

	handleUsernameKeydown = event => {
		if (event.keyCode === 13) { // enter
			this.handleUpdateUsername();
		} else if (event.keyCode === 27) { // esc
			this.handleHideChangeNameForm();
		}
	}
	handleMessageInputKeydown = event => {
		console.log("keydown text", event.keyCode + ", " + this.formMessageInput.value + " , ", this.formMessageInput.value.length)
		// if event.isComposing || 
		if (event.keyCode === 13) { // enter
			if (!this.prepNewLine) {
				// submit
				return;
			}
			this.prepNewLine = false;
		} else {
			this.prepNewLine = false;
		}
		if (event.keyCode === 16 || event.keyCode === 17) { // ctrl, shift
			this.prepNewLine = true;
		}

		// var numChars = this.value.
	}



}
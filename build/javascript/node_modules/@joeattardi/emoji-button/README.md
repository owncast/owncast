![Emoji Button](https://user-images.githubusercontent.com/219285/88242435-2864b400-cc5b-11ea-9b77-b4574ad45f4c.png)


Vanilla JavaScript emoji picker 😎

## Screenshot

![Screenshot](https://user-images.githubusercontent.com/219285/88242157-690ffd80-cc5a-11ea-8b40-fc148d1f7eb7.png)


## Demo and Documentation

[https://emoji-button.js.org](https://emoji-button.js.org)

## Features

* 💻 Vanilla JS, use with any framework
* 😀 Use native or Twemoji emojis
* 🔎 Emoji search
* 👍🏼 Skin tone variations
* ⏱ Recently used emojis
* ⌨️ Fully keyboard accessible
* 🎨 Dark, light, and auto themes
* ⚙️ Add your own custom emoji images
* 🧩 Extend functionality with plugins

## Browser support

Emoji Button is supported on all modern browsers supporting the latest JavaScript features. Internet Explorer is not supported.

## Installation

If you are using a package manager like `yarn` or `npm`, you can install Emoji Button directly from the npm registry:

    npm install @joeattardi/emoji-button

## Basic usage

First, you'll need a trigger element. This is typically a button, and is used to toggle the emoji picker.

```html
<button id="emoji-trigger">Emoji</button>
```

Once you've added the trigger, you will need to import the `EmojiButton` class and create a new instance. Various options can be passed to the constructor, which is covered in the [API documentation](https://emoji-button.js.org/docs/api). 

After constructing a picker, it can be shown by calling `showPicker` or `togglePicker` on it. These functions expect a reference element as a parameter. The picker is positioned relative to this reference element.

When an emoji is selected, the picker will emit an `emoji` event, passing an object containing data about the emoji that was selected. You can then handle the selected emoji however you want.

For more in depth documentation and examples, please visit [https://emoji-button.js.org](https://emoji-button.js.org).

```javascript
import { EmojiButton } from '@joeattardi/emoji-button';

const picker = new EmojiButton();
const trigger = document.querySelector('#emoji-trigger');

picker.on('emoji', selection => {
  // handle the selected emoji here
  console.log(selection.emoji);
});

trigger.addEventListener('click', () => picker.togglePicker(trigger));
```

## Development

The easiest way to hack on Emoji Button is to use the documentation site. First, you should fork this repository.

### Clone the forked repository

    git clone https://github.com/your-name/emoji-button.git

### From the repository root

#### Install dependencies

    npm install

#### Set up the link

    npm link

#### Start the build/watch loop

    npm run build:watch

### From the `site` subdirectory

#### Install dependencies

    npm install

#### Link the library

    npm link @joeattardi/emoji-button

#### Start the documentation site

    npm run develop

### Open the page

http://localhost:8000

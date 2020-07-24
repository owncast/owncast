import { EmojiButton } from '/js/vendor/emojibutton.min.js'
const picker = new EmojiButton({
  zIndex: 100,
  theme: 'dark',
  custom: [
    {
      name: 'Owncast',
      emoji: 'apple-icon.png'
    }
  ]
});
const trigger = document.querySelector('#emoji-button');


picker.on('emoji', selection => {
  // handle the selected emoji here
  console.log(selection);
});

trigger.addEventListener('click', () => picker.togglePicker(trigger));
picker.on('emoji', emoji => {
  if (emoji.url) {
    const url = location.protocol + "//" + location.host + "/" + emoji.url;
    document.querySelector('#message-body-form').innerHTML += "<img width=30px src=\"" + url + "\"/>";
  } else {
    document.querySelector('#message-body-form').innerHTML += emoji.emoji;
  }
});

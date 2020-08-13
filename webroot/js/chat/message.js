import { h, Component, createRef } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
// Initialize htm with Preact
const html = htm.bind(h);

import {messageBubbleColorForString } from '../utils/user-colors.js';

export default class Message extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      displayForm: false,
    };

    this.handleKeydown = this.handleKeydown.bind(this);
    this.handleDisplayForm = this.handleDisplayForm.bind(this);
    this.handleHideForm = this.handleHideForm.bind(this);
    this.handleUpdateUsername = this.handleUpdateUsername.bind(this);
  }



  render(props) {
    const { message, type } = props;
    const { image, author, text } = message;

    const styles = {
      info: {
        display: displayForm || narrowSpace ? 'none' : 'flex',
      },
      form: {
        display: displayForm ? 'flex' : 'none',
      },
    };

    return (
      html`
        <div class="message flex">
          <div class="message-avatar rounded-full flex items-center justify-center" v-bind:style="{ backgroundColor: message.userColor() }">
            <img
              v-bind:src="message.image"
            />
          </div>
          <div class="message-content">
            <p class="message-author text-white font-bold">{{ message.author }}</p>
            <p class="message-text text-gray-400 font-thin " v-html="message.formatText()"></p>
          </div>
      </div>
    `);
  }
}

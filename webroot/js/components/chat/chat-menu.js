import { h, createContext } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import { useState, useEffect, useRef } from '/js/web_modules/preact/hooks.js';
import UsernameForm from './username.js';
import { ChatIcon, UserIcon, CaretDownIcon, AuthIcon } from '../icons/index.js';

const html = htm.bind(h);

const moderatorFlag = html`
  <img src="/img/moderator-nobackground.svg" class="moderator-flag" />
`;

const Context = createContext();

export const ChatMenu = (props) => {
  const {
    username,
    isModerator,
    chatDisabled,
    noVideoContent,
    handleChatPanelToggle,
    onUsernameChange,
    showAuthModal,
    onFocus,
    onBlur,
  } = props;

  const [chatMenuOpen, setChatMenuOpen] = useState(false);
  const [view, setView] = useState('main');

  const chatMenuRef = useRef(undefined);
  closeOnOutsideClick(chatMenuRef, setChatMenuOpen);

  useEffect(() => {
    if (chatMenuOpen) setView('main');
  }, [chatMenuOpen]);

  const authMenuItem =
    showAuthModal &&
    html`<li>
      <button type="button" id="chat-auth" onClick=${showAuthModal}>
        Authenticate
        <span><${AuthIcon} /></span>
      </button>
    </li>`;

  return html`
  <${Context.Provider} value=${props}>
    <div class="chat-menu p-2 relative shadow-lg" ref=${chatMenuRef}>
      <button
        id="chat-menu-button"
        class="flex items-center p-1 bg-transparent rounded-md overflow-hidden text-gray-200 transition duration-150"
        onClick="${() => setChatMenuOpen(!chatMenuOpen)}"
      >
      ${
        !isModerator
          ? html`<${UserIcon} className="w-6 h-6 mr-2" />`
          : moderatorFlag
      }
      <span
        id="username-display"
        class="text-indigo-100 text-sm font-bold truncate overflow-hidden whitespace-no-wrap ${
          isModerator && 'moderator-flag'
        }"
      >
        ${username}
      </span>
       <${CaretDownIcon} className="w-8 h-8"/>
      </button>
      ${
        chatMenuOpen &&
        html` <div
          class="chat-menu-popout shadow-2xl text-gray-100 absolute w-max top-full right-0 z-50 rounded-md p-2 bg-gray-900 fadeIn "
          style=${{ minWidth: '20rem' }}
        >
          ${view === 'main' &&
          html`<ul class="chat-menu-options w-max">
            <li>
              <${UsernameForm}
                username=${username}
                isModerator=${isModerator}
                onUsernameChange=${onUsernameChange}
                onFocus=${onFocus}
                onBlur=${onBlur}
              />
            </li>
            ${authMenuItem}
            <li>
              <button
                type="button"
                id="chat-toggle"
                onClick=${handleChatPanelToggle}
                style=${{
                  display: chatDisabled || noVideoContent ? 'none' : 'flex',
                }}
              >
                <span>Toggle Chat</span>
                <span><${ChatIcon} /></span>
              </button>
            </li>
          </ul>`}
          ${view != 'main' &&
          html`<${SubMenuView} view=${view} setView=${setView} />`}
        </div>`
      }
    </div>
  </${Context.Provider}>`;
};

const SubMenuView = ({ view, setView }) => {
  return html`
    <div className=${`chat-view fadeInRight`}>
      <ul>
        <li>
          <button onClick=${() => setView('main')}>
            <span
              ><${CaretDownIcon} className="w-6 h-6 transform rotate-90"
            /></span>
            <span>Back</span>
          </button>
        </li>
      </ul>
      ${view === 'change_username' && html`<${ChangeUsernameView} />`}
    </div>
  `;
};

function closeOnOutsideClick(ref, setter) {
  useEffect(() => {
    function handleClickOutside(event) {
      if (ref.current && !ref.current.contains(event.target)) {
        setter(undefined);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [ref]);
}

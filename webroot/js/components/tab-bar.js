import { h, Component } from '/js/web_modules/preact.js';

import htm from '/js/web_modules/htm.js';

const html = htm.bind(h);

export default class TabBar extends Component {
  constructor(props) {
    super(props);
    this.state = {
      activeIndex: 0,
    };

    this.handleTabClick = this.handleTabClick.bind(this);
  }

  handleTabClick(index) {
    this.setState({ activeIndex: index });
  }

  render() {
    const { tabs, ariaLabel } = this.props;
    if (!tabs.length) {
      return null;
    }

    if (tabs.length === 1) {
      return html` ${tabs[0].content} `;
    } else {
      return html`
        <div class="tab-bar">
          <div role="tablist" aria-label=${ariaLabel}>
            ${tabs.map((tabItem, index) => {
              const handleClick = () => this.handleTabClick(index);
              return html`
                <button
                  role="tab"
                  aria-selected=${index === this.state.activeIndex}
                  aria-controls=${`tabContent${index}`}
                  id=${`tab-${tabItem.label}`}
                  onclick=${handleClick}
                >
                  ${tabItem.label}
                </button>
              `;
            })}
          </div>
          ${tabs.map((tabItem, index) => {
            return html`
              <div
                tabindex="0"
                role="tabpanel"
                id=${`tabContent${index}`}
                aria-labelledby=${`tab-${tabItem.label}`}
                hidden=${index !== this.state.activeIndex}
              >
                ${tabItem.content}
              </div>
            `;
          })}
        </div>
      `;
    }
  }
}

import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

var refreshTimer;

export default class VideoPoster extends Component {
  constructor(props) {
    super(props);

    this.state = {
      // flipped is the state of showing primary/secondary image views
      flipped: false,

      // url is the currently displayed image url
      url: this.props.src,

      // active states if the refresh timer is active
      active: false,
    };
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.props.src !== prevProps.src) {
      this.setState(
        {
          url: this.props.src,
          oldUrl: this.props.src,
          flipped: false,
        },
        () => {
          if (this.props.active && !this.state.active) {
            this.startRefreshTimer();
          } else if (!this.props.active && this.state.active) {
            this.stopRefreshTimer();
          }
        }
      );
    }
  }

  startRefreshTimer() {
    this.setState({ active: true });
    clearInterval(refreshTimer);

    // Download a new copy of the image every n seconds
    const refreshInterval = 15000;
    refreshTimer = setInterval(() => {
      this.fire();
    }, refreshInterval);

    this.fire();
  }

  fire() {
    const cachebuster = Math.round(new Date().getTime() / 1000);
    const cbUrl = this.props.src + '?cb=' + cachebuster;
    const oldUrl = this.state.url;

    // Download the image
    const img = new Image();
    img.onload = () => {
      this.setState({
        flipped: !this.state.flipped,
        url: cbUrl,
        oldUrl: oldUrl,
      });
    };
    img.src = cbUrl;
  }

  stopRefreshTimer() {
    clearInterval(refreshTimer);
    this.setState({ active: false });
  }

  render() {
    const { url, oldUrl, flipped } = this.state;
    console.log({ url, oldUrl, flipped }, this.props)
    return html`
      <div id="oc-custom-poster">
        <${ThumbImage} url=${!flipped ? oldUrl : url } />
        <${ThumbImage} url=${flipped ? oldUrl : url } visible=${!flipped} />
      </div>
    `;
  }

  // secondaryImageView(url) {
  //   return html` ${this.imageLayer(url, true)} `;
  // }

  // primaryImageView(url, visible) {
  //   return html` ${this.imageLayer(url, visible)} `;
  // }

  // imageLayer(url, visible) {
  //   if (!url) {
  //     return null;
  //   }

  //   return html`
  //     <div
  //       class="custom-thumbnail-image"
  //       style=${{
  //         opacity: visible ? 1.0 : 0,
  //         'background-image': `url(${url})`,
  //       }}
  //     />
  //   `;
  // }
}

function ThumbImage({ url, visible }) {
  if (!url) {
    return null;
  }
  return html`
    <div
      class="custom-thumbnail-image"
      style=${{
        opacity: visible ? 1.0 : 0,
        'background-image': `url(${url})`,
      }}
    />
  `;
}

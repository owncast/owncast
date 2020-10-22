import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import { TEMP_IMAGE } from '../utils/constants.js';

const REFRESH_INTERVAL = 15000;
const POSTER_BASE_URL = '/thumbnail.jpg';

export default class VideoPoster extends Component {
  constructor(props) {
    super(props);

    this.state = {
      // flipped is the state of showing primary/secondary image views
      flipped: false,
      oldUrl: TEMP_IMAGE,
      url: TEMP_IMAGE,
    };

    this.refreshTimer = null;
    this.startRefreshTimer = this.startRefreshTimer.bind(this);
    this.fire = this.fire.bind(this);
    this.setLoaded = this.setLoaded.bind(this);
  }
  componentDidMount() {
    if (this.props.active) {
      this.fire();
      this.startRefreshTimer();
    }
  }
  shouldComponentUpdate(prevProps, prevState) {
    return this.props.active !== prevProps.active ||
      this.props.offlineImage !== prevProps.offlineImage ||
      this.state.url !== prevState.url ||
      this.state.oldUrl !== prevState.oldUrl;
  }
  componentDidUpdate(prevProps) {
    const { active } = this.props;
    const { active: prevActive } = prevProps;

    if (active && !prevActive) {
      this.startRefreshTimer();
    } else if (!active && prevActive) {
      this.stopRefreshTimer();
    }
  }
  componentWillUnmount() {
    this.stopRefreshTimer();
  }

  startRefreshTimer() {
    this.stopRefreshTimer();
    this.fire();
    // Load a new copy of the image every n seconds
    this.refreshTimer = setInterval(this.fire, REFRESH_INTERVAL);
  }

  // load new img
  fire() {
    const cachebuster = Math.round(new Date().getTime() / 1000);
    this.loadingImage = POSTER_BASE_URL + '?cb=' + cachebuster;
    const img = new Image();
    img.onload = this.setLoaded;
    img.src = this.loadingImage;
  }

  setLoaded() {
    const { url: currentUrl, flipped } = this.state;
    this.setState({
      flipped: !flipped,
      url: this.loadingImage,
      oldUrl: currentUrl,
    });
  }

  stopRefreshTimer() {
    clearInterval(this.refreshTimer);
    this.refreshTimer = null;
  }

  render() {
    const { active, offlineImage } = this.props;
    const { url, oldUrl, flipped } = this.state;
    if (!active) {
      return html`
      <div id="oc-custom-poster">
        <${ThumbImage} url=${offlineImage} visible=${true} />
      </div>
    `;
    }
    return html`
      <div id="oc-custom-poster">
        <${ThumbImage} url=${!flipped ? oldUrl : url } visible=${true} />
        <${ThumbImage} url=${flipped ? oldUrl : url } visible=${!flipped} />
      </div>
    `;
  }
}

function ThumbImage({ url, visible }) {
  if (!url) {
    return null;
  }
  return html`
    <div
      class="custom-thumbnail-image"
      style=${{
        opacity: visible ? 1 : 0,
        backgroundImage: `url(${url})`,
      }}
    />
  `;
}

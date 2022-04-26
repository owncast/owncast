/*
VideoPoster is the image that covers up the video component and shows a 
preview of the video, refreshing every N seconds.
It's more complex than it needs to be, using the "double buffer" approach to
cross-fade the two images. Now that we've moved to React we may be able to
simply use some simple cross-fading component.
*/

import { useEffect, useLayoutEffect, useState } from 'react';
import { ReactElement } from 'react-markdown/lib/react-markdown';

const REFRESH_INTERVAL = 15000;
const TEMP_IMAGE = 'http://localhost:8080/logo';
const POSTER_BASE_URL = 'http://localhost:8080/';

export default function VideoPoster(props): ReactElement {
  const { active } = props;
  const [flipped, setFlipped] = useState(false);
  const [oldUrl, setOldUrl] = useState(TEMP_IMAGE);
  const [url, setUrl] = useState(props.url);
  const [currentUrl, setCurrentUrl] = useState(TEMP_IMAGE);
  const [loadingImage, setLoadingImage] = useState(TEMP_IMAGE);
  const [offlineImage, setOfflineImage] = useState(TEMP_IMAGE);

  let refreshTimer = null;

  const setLoaded = () => {
    setFlipped(!flipped);
    setUrl(loadingImage);
    setOldUrl(currentUrl);
  };

  const fire = () => {
    const cachebuster = Math.round(new Date().getTime() / 1000);
    setLoadingImage(`${POSTER_BASE_URL}?cb=${cachebuster}`);
    const img = new Image();
    img.onload = setLoaded;
    img.src = loadingImage;
  };

  const stopRefreshTimer = () => {
    clearInterval(refreshTimer);
    refreshTimer = null;
  };

  const startRefreshTimer = () => {
    stopRefreshTimer();
    fire();
    // Load a new copy of the image every n seconds
    refreshTimer = setInterval(fire, REFRESH_INTERVAL);
  };

  useEffect(() => {
    if (active) {
      fire();
      startRefreshTimer();
    } else {
      stopRefreshTimer();
    }
  }, [active]);

  // On component unmount.
  useLayoutEffect(
    () => () => {
      stopRefreshTimer();
    },
    [],
  );

  // TODO: Replace this with React memo logic.
  // shouldComponentUpdate(prevProps, prevState) {
  //   return (
  //     this.props.active !== prevProps.active ||
  //     this.props.offlineImage !== prevProps.offlineImage ||
  //     this.state.url !== prevState.url ||
  //     this.state.oldUrl !== prevState.oldUrl
  //   );
  // }

  if (!active) {
    return (
      <div id="oc-custom-poster">
        <ThumbImage url={offlineImage} visible />
      </div>
    );
  }

  return (
    <div id="oc-custom-poster">
      <ThumbImage url={!flipped ? oldUrl : url} visible />
      <ThumbImage url={flipped ? oldUrl : url} visible={!flipped} />
    </div>
  );
}

function ThumbImage({ url, visible }) {
  if (!url) {
    return null;
  }
  return (
    <div
      className="custom-thumbnail-image"
      style={{
        opacity: visible ? 1 : 0,
        backgroundImage: `url(${url})`,
      }}
    />
  );
}

import { FC } from 'react';
import { OwncastLogo } from '../../common/OwncastLogo/OwncastLogo';
import styles from './Noscript.module.scss';

export const Noscript: FC = () => (
  <noscript className={styles.noscript}>
    <div className={styles.scrollContainer}>
      <div className={styles.content}>
        <OwncastLogo className={styles.logo} />
        <br />
        <p>
          This website is powered by&nbsp;
          <a href="https://owncast.online" rel="noopener noreferrer" target="_blank">
            Owncast
          </a>
          .
        </p>
        <p>
          Owncast uses JavaScript for playing the HTTP Live Streaming (HLS) video, and its chat
          client. But your web browser does not seem to support JavaScript, or you have it disabled.
        </p>
        <p>
          For the best experience, you should use a different browser with JavaScript support. If
          you have disabled JavaScript in your browser, you can re-enable it.
        </p>
        <h2>How can I watch this stream without JavaScript?</h2>
        <p>
          You can open the URL of this website in your media player (such as&nbsp;
          <a href="https://mpv.io" rel="noopener noreferrer" target="_blank">
            mpv
          </a>
          &nbsp;or&nbsp;
          <a href="https://www.videolan.org/vlc/" rel="noopener noreferrer" target="_blank">
            VLC
          </a>
          ) to watch the stream.
        </p>
        <h2>How can I chat with the others without JavaScript?</h2>
        <p>Currently, there is no option to use the chat without JavaScript.</p>
      </div>
    </div>
  </noscript>
);

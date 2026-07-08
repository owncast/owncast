/* eslint-disable max-classes-per-file */
import React, { FC, useContext, useEffect } from 'react';
import { useAtom, useAtomValue } from 'jotai';
import { useHotkeys } from 'react-hotkeys-hook';
import { useTranslation } from 'next-export-i18n';
import classNames from 'classnames';
import { ErrorBoundary, getErrorMessage } from 'react-error-boundary';
import { VideoJS } from '../VideoJS/VideoJS';
import ViewerPing from '../viewer-ping';
import { VideoPoster } from '../VideoPoster/VideoPoster';
import { getLocalStorage, setLocalStorage } from '../../../utils/localStorage';
import { autoplayModeForSetting } from '../../../utils/autoplay';
import { Localization } from '../../../types/localization';
import { isVideoPlayingAtom, clockSkewAtom } from '../../stores/ClientConfigStore';
import PlaybackMetrics from '../metrics/playback';
import { createVideoSettingsMenuButton, LATENCY_COMPENSATION_ENABLED } from '../settings-menu';
import LatencyCompensator from '../latencyCompensator';
import styles from './OwncastPlayer.module.scss';
import { VideoSettingsServiceContext } from '../../../services/video-settings-service';
import { ComponentError } from '../../ui/ComponentError/ComponentError';

const PLAYER_VOLUME = 'owncast_volume';

const ping = new ViewerPing();
let playbackMetrics = null;
let latencyCompensator = null;
let latencyCompensatorEnabled = false;

export type OwncastPlayerProps = {
  source: string;
  online: boolean;
  initiallyMuted?: boolean;
  autoplay?: string;
  title: string;
  className?: string;
};

export const OwncastPlayer: FC<OwncastPlayerProps> = ({
  source,
  online,
  initiallyMuted = false,
  autoplay = 'off',
  title,
  className,
}) => {
  const VideoSettingsService = useContext(VideoSettingsServiceContext);
  const playerRef = React.useRef(null);
  const [videoPlaying, setVideoPlaying] = useAtom(isVideoPlayingAtom);
  const clockSkew = useAtomValue(clockSkewAtom);
  const { t } = useTranslation();

  const setSavedVolume = () => {
    try {
      playerRef.current.volume(getLocalStorage(PLAYER_VOLUME) || 1);
    } catch (err) {
      console.warn(err);
    }
  };

  const handleVolume = () => {
    setLocalStorage(PLAYER_VOLUME, playerRef.current.muted() ? 0 : playerRef.current.volume());
  };

  const togglePlayback = () => {
    if (playerRef.current.paused()) {
      playerRef.current.play();
    } else {
      playerRef.current.pause();
    }
  };

  const toggleMute = () => {
    if (playerRef.current.muted() || playerRef.current.volume() === 0) {
      playerRef.current.volume(0.7);
    } else {
      playerRef.current.volume(0);
    }
  };

  const toggleFullScreen = () => {
    if (playerRef.current.isFullscreen()) {
      playerRef.current.exitFullscreen();
    } else {
      playerRef.current.requestFullscreen();
    }
  };

  const startLatencyCompensator = () => {
    if (latencyCompensator) {
      // Fully tear down the old instance so its check timer doesn't leak.
      latencyCompensator.disable();
    }

    latencyCompensatorEnabled = true;

    latencyCompensator = new LatencyCompensator(playerRef.current);
    latencyCompensator.setClockSkew(clockSkew);
    latencyCompensator.enable();
    setLocalStorage(LATENCY_COMPENSATION_ENABLED, true);
  };

  const stopLatencyCompensator = () => {
    if (latencyCompensator) {
      latencyCompensator.disable();
    }
    latencyCompensator = null;
    latencyCompensatorEnabled = false;
    setLocalStorage(LATENCY_COMPENSATION_ENABLED, false);
  };

  // Toggle minimized latency mode. Return the new state.
  const toggleLatencyCompensator = () => {
    if (latencyCompensatorEnabled) {
      stopLatencyCompensator();
    } else {
      startLatencyCompensator();
    }
    return latencyCompensatorEnabled;
  };

  const setupLatencyCompensator = player => {
    const tech = player.tech({ IWillNotUseThisInPlugins: true });

    // VHS is required.
    if (!tech || !tech.vhs) {
      return;
    }

    const latencyCompensatorEnabledSaved = getLocalStorage(LATENCY_COMPENSATION_ENABLED);

    if (latencyCompensatorEnabledSaved === 'true' && tech && tech.vhs) {
      startLatencyCompensator();
    } else {
      stopLatencyCompensator();
    }
  };

  const createSettings = async (player, videojs) => {
    const videoQualities = await VideoSettingsService.getVideoQualities();
    const menuButton = createVideoSettingsMenuButton(
      player,
      videojs,
      videoQualities,
      toggleLatencyCompensator,
    );
    player.controlBar.addChild(
      menuButton,
      {},
      // eslint-disable-next-line no-underscore-dangle
      player.controlBar.children_.length - 2,
    );
    setupLatencyCompensator(player);
  };

  const setupAirplay = (player, videojs) => {
    // eslint-disable-next-line no-prototype-builtins
    if (window.hasOwnProperty('WebKitPlaybackTargetAvailabilityEvent')) {
      const VJSButtonClass = videojs.getComponent('Button');

      class ConcreteButtonClass extends VJSButtonClass {
        constructor() {
          super(player);
        }

        // eslint-disable-next-line class-methods-use-this
        handleClick() {
          try {
            const videoElement = document.getElementsByTagName('video')[0];
            (videoElement as any).webkitShowPlaybackTargetPicker();
          } catch (e) {
            console.error(e);
          }
        }
      }

      const ccbc = new ConcreteButtonClass();
      const concreteButtonInstance = player.controlBar.addChild(ccbc);
      concreteButtonInstance.addClass('vjs-airplay');
    }
  };

  const setupUnmuteButton = (player, videojs) => {
    const VJSButtonClass = videojs.getComponent('Button');

    class UnmuteButton extends VJSButtonClass {
      constructor() {
        super(player);
      }

      // eslint-disable-next-line class-methods-use-this
      handleClick() {
        player.muted(false);
        if (player.volume() === 0) {
          const saved = parseFloat(getLocalStorage(PLAYER_VOLUME));
          player.volume(saved > 0 ? saved : 0.7);
        }
      }
    }

    const unmuteButton = new UnmuteButton();
    unmuteButton.addClass('vjs-big-unmute-button');
    // Localize the control text (exposed to assistive tech and tooltips),
    // falling back to English when the key has no translation, mirroring the
    // Translation component's missing-key behavior.
    const unmuteLabel = t(Localization.Frontend.unmute);
    unmuteButton.controlText(unmuteLabel === Localization.Frontend.unmute ? 'Unmute' : unmuteLabel);
    // The overlay glyph is antd's MutedFilled speaker (the icon family used by
    // the rest of the UI) instead of the video.js icon font. Inlined as SVG
    // because this button renders into video.js DOM, outside React. Path from
    // @ant-design/icons-svg MutedFilled. fill=currentColor picks up the theme
    // action color set in VideoJS.scss.
    const iconPlaceholder = unmuteButton.el().querySelector('.vjs-icon-placeholder');
    if (iconPlaceholder) {
      iconPlaceholder.innerHTML =
        '<svg viewBox="64 64 896 896" fill="currentColor" fill-rule="evenodd" aria-hidden="true" focusable="false"><path d="M771.91 115a31.65 31.65 0 00-17.42 5.27L400 351.97H236a16 16 0 00-16 16v288.06a16 16 0 0016 16h164l354.5 231.7a31.66 31.66 0 0017.42 5.27c16.65 0 32.08-13.25 32.08-32.06V147.06c0-18.8-15.44-32.06-32.09-32.06"></path></svg>';
    }
    player.addChild(unmuteButton);

    // Mirror the big play button: show a large, obvious unmute affordance only
    // while the stream is autoplaying muted. When paused, the native big play
    // button already covers that case, so this stays hidden.
    const updateUnmuteButton = () => {
      if (!player.paused() && (player.muted() || player.volume() === 0)) {
        unmuteButton.show();
      } else {
        unmuteButton.hide();
      }
    };

    player.on(['playing', 'pause', 'ended', 'volumechange'], updateUnmuteButton);
    updateUnmuteButton();
  };

  // Register keyboard shortcut for the space bar to toggle playback
  useHotkeys('space', e => {
    e.preventDefault();
    togglePlayback();
  });

  // Register keyboard shortcut for f to toggle full screen
  useHotkeys('f', toggleFullScreen, {
    enableOnContentEditable: false,
  });

  // Register keyboard shortcut for the "m" key to toggle mute
  useHotkeys('m', toggleMute, {
    enableOnContentEditable: false,
  });

  useHotkeys('0', () => playerRef.current.volume(playerRef.current.volume() + 0.1), {
    enableOnContentEditable: false,
  });
  useHotkeys('9', () => playerRef.current.volume(playerRef.current.volume() - 0.1), {
    enableOnContentEditable: false,
  });

  // Resolve the video.js autoplay value from the instance setting, honoring the
  // viewer's data-saver and reduced-motion preferences. Guarded because this
  // runs during render (no window during SSR).
  const prefersReducedMotion =
    typeof window !== 'undefined' &&
    typeof window.matchMedia === 'function' &&
    window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  const navConnection =
    typeof navigator !== 'undefined'
      ? (navigator as Navigator & { connection?: { saveData?: boolean } }).connection
      : undefined;
  // The player restores the persisted volume on ready, and a mute is persisted
  // as volume 0. Starting muted (embed request) or at volume 0 makes a
  // sound-only autoplay attempt succeed silently in Firefox, which allows
  // inaudible autoplay, so the mapping needs to know about it.
  const savedVolume = typeof window !== 'undefined' ? getLocalStorage(PLAYER_VOLUME) : null;
  const autoplayMode = autoplayModeForSetting(autoplay, {
    prefersReducedMotion,
    saveData: navConnection?.saveData === true,
    startsInaudible: initiallyMuted || (savedVolume !== null && Number(savedVolume) === 0),
  });

  const videoJsOptions = {
    autoplay: autoplayMode,
    controls: true,
    responsive: true,
    fluid: false,
    fill: true,
    playsinline: true,
    liveui: true,
    preload: 'auto',
    muted: initiallyMuted,
    controlBar: {
      progressControl: {
        seekBar: false,
      },
    },
    html5: {
      vhs: {
        // used to select the lowest bitrate playlist initially. This helps to decrease playback start time. This setting is false by default.
        enableLowInitialPlaylist: true,
        experimentalBufferBasedABR: true,
        useNetworkInformationApi: true,
        maxPlaylistRetries: 30,
      },
    },
    liveTracker: {
      trackingThreshold: 0,
      liveTolerance: 15,
    },
    sources: [
      {
        src: source,
        type: 'application/x-mpegURL',
      },
    ],
  };

  const handlePlayerReady = (player, videojs) => {
    playerRef.current = player;
    setSavedVolume();
    setupAirplay(player, videojs);
    setupUnmuteButton(player, videojs);

    // You can handle player events here, for example:
    player.on('waiting', () => {
      console.debug('player is waiting');
    });

    player.on('dispose', () => {
      console.debug('player will dispose');
      ping.stop();
    });

    player.on('playing', () => {
      console.debug('player is playing');
      ping.start();
      setVideoPlaying(true);
    });

    player.on('pause', () => {
      console.debug('player is paused');
      ping.stop();
      setVideoPlaying(false);
    });

    player.on('ended', () => {
      console.debug('player is ended');
      ping.stop();
      setVideoPlaying(false);
    });

    videojs.hookOnce();

    player.on('volumechange', handleVolume);

    playbackMetrics = new PlaybackMetrics(player, videojs);
    playbackMetrics.setClockSkew(clockSkew);

    createSettings(player, videojs);
  };

  useEffect(() => {
    if (playbackMetrics) {
      playbackMetrics.setClockSkew(clockSkew);
    }
  }, [clockSkew]);

  useEffect(
    () => () => {
      stopLatencyCompensator();
      playbackMetrics?.stop();
    },
    [],
  );

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentError
          componentName="OwncastPlayer"
          message={getErrorMessage(error)}
          retryFunction={resetErrorBoundary}
        />
      )}
    >
      <div className={classNames(styles.container, className)} id="player">
        {online && (
          <div className={styles.player}>
            <VideoJS options={videoJsOptions} onReady={handlePlayerReady} aria-label={title} />
          </div>
        )}
        <div className={styles.poster}>
          {!videoPlaying && (
            <VideoPoster online={online} initialSrc="/thumbnail.jpg" src="/thumbnail.jpg" />
          )}
        </div>
      </div>
    </ErrorBoundary>
  );
};

import React, { useContext, useState, useEffect } from 'react';
import { Typography, Slider } from 'antd';
import { ServerStatusContext } from '../../utils/server-status-context';
import { AlertMessageContext } from '../../utils/alert-message-context';
import {
  API_VIDEO_SEGMENTS,
  RESET_TIMEOUT,
  postConfigUpdateToAPI,
} from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import FormStatusIndicator from './form-status-indicator';

const { Title } = Typography;

const SLIDER_MARKS = {
  1: 'Low',
  2: '',
  3: '',
  4: '',
  5: '',
  6: 'High',
};

const SLIDER_COMMENTS = {
  1: 'Lowest latency, lowest error tolerance',
  2: 'Low latency, low error tolerance',
  3: 'Lower latency, lower error tolerance',
  4: 'Medium latency, medium error tolerance (Default)',
  5: 'Higher latency, higher error tolerance',
  6: 'Highest latency, highest error tolerance',
};

export default function VideoLatency() {
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  const [selectedOption, setSelectedOption] = useState(null);

  const serverStatusData = useContext(ServerStatusContext);
  const { setMessage } = useContext(AlertMessageContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { videoSettings } = serverConfig || {};

  let resetTimer = null;

  if (!videoSettings) {
    return null;
  }

  useEffect(() => {
    setSelectedOption(videoSettings.latencyLevel);
  }, [videoSettings]);

  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  // posts all the variants at once as an array obj
  const postUpdateToAPI = async (postValue: any) => {
    setSubmitStatus(createInputStatus(STATUS_PROCESSING));

    await postConfigUpdateToAPI({
      apiPath: API_VIDEO_SEGMENTS,
      data: { value: postValue },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'latencyLevel',
          value: postValue,
          path: 'videoSettings',
        });
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Latency buffer level updated.'));

        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
        if (serverStatusData.online) {
          setMessage(
            'Your latency buffer setting will take effect the next time you begin a live stream.',
          );
        }
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));

        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  const handleChange = value => {
    postUpdateToAPI(value);
  };

  return (
    <div className="config-video-latency-container">
      <Title level={3} className="section-title">
        Latency Buffer
      </Title>
      <p className="description">
        While it&apos;s natural to want to keep your latency as low as possible, you may experience
        reduced error tolerance and stability in some environments the lower you go.
      </p>
      <p className="description">
        For interactive live streams you may want to experiment with a lower latency, for
        non-interactive broadcasts you may want to increase it.{' '}
        <a
          href="https://owncast.online/docs/encoding#latency-buffer?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          Read to learn more.
        </a>
      </p>

      <div className="segment-slider-container">
        <Slider
          tipFormatter={value => SLIDER_COMMENTS[value]}
          onChange={handleChange}
          min={1}
          max={6}
          marks={SLIDER_MARKS}
          defaultValue={selectedOption}
          value={selectedOption}
        />
        <p className="selected-value-note">{SLIDER_COMMENTS[selectedOption]}</p>
        <FormStatusIndicator status={submitStatus} />
      </div>
    </div>
  );
}

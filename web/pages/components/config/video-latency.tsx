import React, { useContext, useState, useEffect } from 'react';
import { Typography, Slider } from 'antd';
import { ServerStatusContext } from '../../../utils/server-status-context';
import { API_VIDEO_SEGMENTS, RESET_TIMEOUT, postConfigUpdateToAPI } from './constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../../utils/input-statuses';
import FormStatusIndicator from './form-status-indicator';

const { Title } = Typography;

const SLIDER_MARKS = {
  1: 'low',
  2: '',
  3: '',
  4: '',
  5: '',
  6: 'high',
};

const SLIDER_COMMENTS = {
  1: 'Lowest latency, but least reliability',
  2: 'Low latency, some reliability',
  3: 'Lower latency, some reliability',
  4: 'Optimal latency and reliability (Default)',
  5: 'High latency, better reliability',
  6: 'Highest latency, higher reliability',
};

interface SegmentToolTipProps {
  value: string;
}

function SegmentToolTip({ value }: SegmentToolTipProps) {
  return <span className="segment-tip">{value}</span>;
}

export default function VideoLatency() {
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);

  // const [submitStatus, setSubmitStatus] = useState(null);
  // const [submitStatusMessage, setSubmitStatusMessage] = useState('');
  const [selectedOption, setSelectedOption] = useState(null);

  const serverStatusData = useContext(ServerStatusContext);
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
    // setSubmitStatusMessage('');
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
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Variants updated.'));

        // setSubmitStatus('success');
        // setSubmitStatusMessage('Variants updated.');
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));

        // setSubmitStatus('error');
        // setSubmitStatusMessage(message);
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  const handleChange = value => {
    postUpdateToAPI(value);
  };

  return (
    <div className="config-video-segements-conatiner">
      <Title level={3}>Latency Buffer</Title>
      <p>
        There are trade-offs when cosidering video latency and reliability. Blah blah .. better
        wording here needed.
      </p>
      <br />
      <br />
      <div className="segment-slider-container">
        <Slider
          tipFormatter={value => <SegmentToolTip value={SLIDER_COMMENTS[value]} />}
          onChange={handleChange}
          min={1}
          max={6}
          marks={SLIDER_MARKS}
          defaultValue={selectedOption}
          value={selectedOption}
        />
      </div>
      <FormStatusIndicator status={submitStatus} />
    </div>
  );
}

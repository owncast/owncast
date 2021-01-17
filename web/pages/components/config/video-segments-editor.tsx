import React, { useContext, useState } from 'react';
import { Typography, Slider, } from 'antd';
import { ServerStatusContext } from '../../../utils/server-status-context';
import { API_VIDEO_SEGMENTS, SUCCESS_STATES, RESET_TIMEOUT,postConfigUpdateToAPI } from './constants';

const { Title } = Typography;

// numberOfSegments
// secondsPerSegment
/*

2 segments, 2 seconds
3 segments, 3 seconds
3 segments, 4 seconds
4 segments, 4 seconds <- default
8 segments, 4 seconds
10 segments, 4 seconds

Lowest latancy, less reliability 
-> Highest latancy, higher reliability
*/
const DEFAULT_OPTION = 3;

const SLIDER_OPTIONS = [
  {
    numberOfSegments: 2,
    secondsPerSegment: 2,
  },
  {
    numberOfSegments: 3,
    secondsPerSegment: 3,
  },
  {
    numberOfSegments: 3,
    secondsPerSegment: 4,
  },
  {
    numberOfSegments: 4,
    secondsPerSegment: 4,
  },
  {
    numberOfSegments: 8,
    secondsPerSegment: 4,
  },
  {
    numberOfSegments: 10,
    secondsPerSegment: 4,
  },
];

const SLIDER_MARKS = {
  0: '1',
  1: '2',
  2: '3',
  3: '4',
  4: '5',
  5: '6',
};

const SLIDER_COMMENTS = {
  0: 'Lowest latency, but least reliability',
  1: 'Low latency, some reliability',
  2: 'Lower latency, some reliability',
  3: 'Optimal latency and reliability (Default)',
  4: 'High latency, better reliability',
  5: 'Highest latency, higher reliability',
};

interface SegmentToolTipProps {
  value: string;
}

function SegmentToolTip({ value }: SegmentToolTipProps) {
  return (
    <span className="segment-tip">{value}</span>
  );
}

function findSelectedOption(videoSettings) {
  const { numberOfSegments, secondsPerSegment } = videoSettings;
  let selected = DEFAULT_OPTION;
  SLIDER_OPTIONS.map((option, index) => {
    if (
      (option.numberOfSegments === numberOfSegments &&
      option.secondsPerSegment === secondsPerSegment) ||
      option.numberOfSegments === numberOfSegments
    ) {
      selected = index;
    }
    return option;
  });
  return selected;
}

export default function VideoSegmentsEditor() {
  const [submitStatus, setSubmitStatus] = useState(null);
  const [submitStatusMessage, setSubmitStatusMessage] = useState('');

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { videoSettings } = serverConfig || {};

  let resetTimer = null;

  if (!videoSettings) {
    return null;
  }

  const selectedOption = findSelectedOption(videoSettings);

  const resetStates = () => {
    setSubmitStatus(null);
    setSubmitStatusMessage('');
    resetTimer = null;
    clearTimeout(resetTimer);
  }
  
  // posts all the variants at once as an array obj
  const postUpdateToAPI = async (postValue: any) => {
    await postConfigUpdateToAPI({
      apiPath: API_VIDEO_SEGMENTS,
      data: { value: postValue },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'numberOfSegments', 
          value: postValue.numberOfSegments, 
          path: 'videoSettings'
        });
        setFieldInConfigState({
          fieldName: 'secondsPerSegment',
          value: postValue.secondsPerSegment,
          path: 'videoSettings',
        });

        setSubmitStatus('success');
        setSubmitStatusMessage('Variants updated.');
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus('error');
        setSubmitStatusMessage(message);
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  const {
    icon: newStatusIcon = null,
    message: newStatusMessage = '',
  } = SUCCESS_STATES[submitStatus] || {};
  
  const statusMessage = (
    <div className={`status-message ${submitStatus || ''}`}>
      {newStatusIcon} {newStatusMessage} {submitStatusMessage}
    </div>
  );
  
  const handleSegmentChange = value => {
    const postData = SLIDER_OPTIONS[value];
    postUpdateToAPI(postData);
  };

  return (
    <div className="module-container config-video-segements-conatiner">
      <Title level={3}>Video Tolerance</Title>
      <p>
        There are trade-offs when cosidering video latency and reliability. Blah blah .. better wording here needed.
      </p>
      <br /><br />
      <div className="segment-slider">
        <Slider 
          tooltipVisible
          tipFormatter={value => <SegmentToolTip value={SLIDER_COMMENTS[value]} />}
          onChange={handleSegmentChange}
          min={0}
          max={SLIDER_OPTIONS.length - 1}
          marks={SLIDER_MARKS}
          defaultValue={DEFAULT_OPTION}
          value={selectedOption}
        />
      </div>
      {statusMessage}
    </div>
  );
}
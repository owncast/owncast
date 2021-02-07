import React, { useContext, useState, useEffect } from 'react';
import { Typography, Slider } from 'antd';
import { ServerStatusContext } from '../../utils/server-status-context';

const { Title } = Typography;

const SLIDER_MARKS = {
  1: 'lowest',
  2: '',
  3: '',
  4: '',
  5: 'highest',
};

const TOOLTIPS = {
  1: 'lowest',
  2: 'low',
  3: 'medium',
  4: 'high',
  5: 'highest',
};

export default function CPUUsageSelector({ defaultValue, onChange }) {
  const [selectedOption, setSelectedOption] = useState(null);

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};
  const { videoSettings } = serverConfig || {};

  if (!videoSettings) {
    return null;
  }

  useEffect(() => {
    setSelectedOption(defaultValue);
  }, [videoSettings]);

  const handleChange = value => {
    setSelectedOption(value);
    onChange(value);
  };

  return (
    <div className="config-video-segements-conatiner">
      <Title level={3}>CPU Usage</Title>
      <p>There are trade-offs when considering CPU usage blah blah more wording here.</p>
      <br />
      <br />
      <div className="segment-slider-container">
        <Slider
          tipFormatter={value => TOOLTIPS[value]}
          onChange={handleChange}
          min={1}
          max={Object.keys(SLIDER_MARKS).length}
          marks={SLIDER_MARKS}
          defaultValue={selectedOption}
          value={selectedOption}
        />
      </div>
    </div>
  );
}

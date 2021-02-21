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
  1: 'Lowest CPU usage - lowest quality video',
  2: 'Low CPU usage - low quality video',
  3: 'Medium CPU usage - average quality video',
  4: 'High CPU usage - high quality video',
  5: 'Highest CPU usage - higher quality video',
};
interface Props {
  defaultValue: number;
  disabled: boolean;
  onChange: (arg: number) => void;
}
export default function CPUUsageSelector({ defaultValue, disabled, onChange }: Props) {
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

  const cpuUsageNote = () => {
    if (disabled) {
      return 'CPU usage selection is disabled when Video Passthrough is enabled.';
    }

    return TOOLTIPS[selectedOption]
  }

  return (
    <div className="config-video-cpu-container">
      <Title level={3}>CPU Usage</Title>
      <p className="description">
        Reduce to improve server performance, or increase it to improve video quality.
      </p>
      <div className="segment-slider-container">
        <Slider
          tipFormatter={value => TOOLTIPS[value]}
          onChange={handleChange}
          min={1}
          max={Object.keys(SLIDER_MARKS).length}
          marks={SLIDER_MARKS}
          defaultValue={selectedOption}
          value={selectedOption}
          disabled={disabled}
        />
        <p className="selected-value-note">{cpuUsageNote()}</p>
      </div>
    </div>
  );
}

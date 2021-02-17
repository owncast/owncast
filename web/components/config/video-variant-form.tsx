// This content populates the video variant modal, which is spawned from the variants table.
import React from 'react';
import { Row, Col, Slider, Collapse, Typography } from 'antd';
import { FieldUpdaterFunc, VideoVariant, UpdateArgs } from '../../types/config-section';
import TextField from './form-textfield';
import { DEFAULT_VARIANT_STATE } from '../../utils/config-constants';
import CPUUsageSelector from './cpu-usage';
import ToggleSwitch from './form-toggleswitch';

const { Panel } = Collapse;

const VIDEO_VARIANT_DEFAULTS = {
  framerate: {
    min: 24,
    max: 120,
    defaultValue: 24,
    unit: 'fps',
    incrementBy: null,
    tip:
      'Reducing your framerate will decrease the amount of video that needs to be encoded and sent to your viewers, saving CPU and bandwidth at the expense of smoothness.  A lower value is generally is fine for most content.',
  },
  videoBitrate: {
    min: 600,
    max: 6000,
    defaultValue: 1200,
    unit: 'kbps',
    incrementBy: 100,
    tip: 'The overall quality of your stream is generally impacted most by bitrate.',
  },
  audioBitrate: {
    min: 600,
    max: 1200,
    defaultValue: 800,
    unit: 'kbps',
    incrementBy: 100,
    tip: 'nothing to see here',
  },
  videoPassthrough: {
    tip: 'If No is selected, then you should set your desired Video Bitrate.',
  },
  audioPassthrough: {
    tip: 'If No is selected, then you should set your desired Audio Bitrate.',
  },
  scaledWidth: {
    fieldName: 'scaledWidth',
    label: 'Resized Width',
    maxLength: 4,
    placeholder: '1080',
    tip: "Optionally resize this content's width.",
  },
  scaledHeight: {
    fieldName: 'scaledHeight',
    label: 'Resized Height',
    maxLength: 4,
    placeholder: '720',
    tip: "Optionally resize this content's height.",
  },
};
interface VideoVariantFormProps {
  dataState: VideoVariant;
  onUpdateField: FieldUpdaterFunc;
}

export default function VideoVariantForm({
  dataState = DEFAULT_VARIANT_STATE,
  onUpdateField,
}: VideoVariantFormProps) {
  const handleFramerateChange = (value: number) => {
    onUpdateField({ fieldName: 'framerate', value });
  };
  const handleVideoBitrateChange = (value: number) => {
    onUpdateField({ fieldName: 'videoBitrate', value });
  };
  const handleVideoPassChange = (value: boolean) => {
    onUpdateField({ fieldName: 'videoPassthrough', value });
  };
  const handleVideoCpuUsageLevelChange = (value: number) => {
    onUpdateField({ fieldName: 'cpuUsageLevel', value });
  };
  const handleScaledWidthChanged = (args: UpdateArgs) => {
    const value = Number(args.value);
    // eslint-disable-next-line no-restricted-globals
    if (isNaN(value)) {
      return;
    }
    onUpdateField({ fieldName: 'scaledWidth', value: value || '' });
  };
  const handleScaledHeightChanged = (args: UpdateArgs) => {
    const value = Number(args.value);
    // eslint-disable-next-line no-restricted-globals
    if (isNaN(value)) {
      return;
    }

    onUpdateField({ fieldName: 'scaledHeight', value: value || '' });
  };
  const framerateDefaults = VIDEO_VARIANT_DEFAULTS.framerate;
  const framerateMin = framerateDefaults.min;
  const framerateMax = framerateDefaults.max;
  const framerateUnit = framerateDefaults.unit;
  const framerateMarks = {
    [framerateMin]: `${framerateMin} ${framerateUnit}`,
    30: '',
    60: '',
    90: '',
    [framerateMax]: `${framerateMax} ${framerateUnit}`,
  };

  const videoBitrateDefaults = VIDEO_VARIANT_DEFAULTS.videoBitrate;
  const videoBRMin = videoBitrateDefaults.min;
  const videoBRMax = videoBitrateDefaults.max;
  const videoBRUnit = videoBitrateDefaults.unit;
  const videoBRMarks = {
    [videoBRMin]: `${videoBRMin} ${videoBRUnit}`,
    3000: 3000,
    4500: 4500,
    [videoBRMax]: `${videoBRMax} ${videoBRUnit}`,
  };

  const selectedVideoBRnote = () => {
    let note = `${dataState.videoBitrate}${videoBRUnit}`;
    if (dataState.videoBitrate < 2000) {
      note = `${note} - Good for low bandwidth environments.`;
    } else if (dataState.videoBitrate < 3500) {
      note = `${note} - Good for most bandwidth environments.`;
    } else {
      note = `${note} - Good for high bandwidth environments.`;
    }
    return note;
  };
  const selectedFramerateNote = () => {
    let note = `Selected: ${dataState.framerate}${framerateUnit}`;
    switch (dataState.framerate) {
      case 24:
        note = `${note} - Good for film, presentations, music, low power/bandwidth servers.`;
        break;
      case 30:
        note = `${note} - Good for slow/casual games, chat, general purpose.`;
        break;
      case 60:
        note = `${note} - Good for fast/action games, sports, HD video.`;
        break;
      case 90:
        note = `${note} - Good for newer fast games and hardware.`;
        break;
      case 120:
        note = `${note} - Experimental, use at your own risk!`;
        break;
      default:
        note = '';
    }
    return note;
  };

  return (
    <div className="config-variant-form">
      <p className="description">
        <a href="https://owncast.online/docs/video">Learn more</a> about how each of these settings
        can impact the performance of your server.
      </p>

      <Row gutter={16}>
        <Col sm={24} md={12}>
          {/* ENCODER PRESET FIELD */}
          <div className="form-module cpu-usage-container">
            <CPUUsageSelector
              defaultValue={dataState.cpuUsageLevel}
              onChange={handleVideoCpuUsageLevelChange}
            />
            <p className="read-more-subtext">
              <a href="https://owncast.online/docs/video/#cpu-usage">Read more about CPU usage.</a>
            </p>
          </div>

          {/* VIDEO PASSTHROUGH FIELD - currently disabled */}
          <div style={{ display: 'none' }} className="form-module">
            <ToggleSwitch
              label="Use Video Passthrough?"
              fieldName="video-passthrough"
              tip={VIDEO_VARIANT_DEFAULTS.videoPassthrough.tip}
              checked={dataState.videoPassthrough}
              onChange={handleVideoPassChange}
            />
          </div>
        </Col>

        <Col sm={24} md={12}>
          {/* VIDEO BITRATE FIELD */}
          <div
            className={`form-module bitrate-container ${
              dataState.videoPassthrough ? 'disabled' : ''
            }`}
          >
            <Typography.Title level={3}>Video Bitrate</Typography.Title>
            <p className="description">{VIDEO_VARIANT_DEFAULTS.videoBitrate.tip}</p>
            <div className="segment-slider-container">
              <Slider
                tipFormatter={value => `${value} ${videoBRUnit}`}
                disabled={dataState.videoPassthrough === true}
                defaultValue={dataState.videoBitrate}
                value={dataState.videoBitrate}
                onChange={handleVideoBitrateChange}
                step={videoBitrateDefaults.incrementBy}
                min={videoBRMin}
                max={videoBRMax}
                marks={videoBRMarks}
              />
              <p className="selected-value-note">{selectedVideoBRnote()}</p>
            </div>
            <p className="read-more-subtext">
              <a href="https://owncast.online/docs/video/#bitrate">Read more about bitrates.</a>
            </p>
          </div>
        </Col>
      </Row>
      <Collapse className="advanced-settings">
        <Panel header="Advanced Settings" key="1">
          <p className="description">
            Resizing your content will take additional resources on your server. If you wish to
            optionally resize your content for this stream output then you should either set the
            width <strong>or</strong> the height to keep your aspect ratio.{' '}
            <a href="https://owncast.online/docs/video/#resolution">Read more about resolutions.</a>
          </p>

          <TextField
            type="number"
            {...VIDEO_VARIANT_DEFAULTS.scaledWidth}
            value={dataState.scaledWidth}
            onChange={handleScaledWidthChanged}
          />
          <TextField
            type="number"
            {...VIDEO_VARIANT_DEFAULTS.scaledHeight}
            value={dataState.scaledHeight}
            onChange={handleScaledHeightChanged}
          />

          {/* FRAME RATE FIELD */}
          <div className="form-module frame-rate">
            <Typography.Title level={3}>Frame rate</Typography.Title>
            <p className="description">{VIDEO_VARIANT_DEFAULTS.framerate.tip}</p>
            <div className="segment-slider-container">
              <Slider
                tipFormatter={value => `${value} ${framerateUnit}`}
                defaultValue={dataState.framerate}
                value={dataState.framerate}
                onChange={handleFramerateChange}
                step={framerateDefaults.incrementBy}
                min={framerateMin}
                max={framerateMax}
                marks={framerateMarks}
              />
              <p className="selected-value-note">{selectedFramerateNote()}</p>
            </div>
            <p className="read-more-subtext">
              <a href="https://owncast.online/docs/video/#framerate">Read more about framerates.</a>
            </p>
          </div>
        </Panel>
      </Collapse>
    </div>
  );
}

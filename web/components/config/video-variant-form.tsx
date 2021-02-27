// This content populates the video variant modal, which is spawned from the variants table. This relies on the `dataState` prop fed in by the table.
import React from 'react';
import { Popconfirm, Row, Col, Slider, Collapse, Typography } from 'antd';
import { ExclamationCircleFilled } from '@ant-design/icons';
import classNames from 'classnames';
import { FieldUpdaterFunc, VideoVariant, UpdateArgs } from '../../types/config-section';
import TextField from './form-textfield';
import {
  DEFAULT_VARIANT_STATE,
  VIDEO_VARIANT_SETTING_DEFAULTS,
  ENCODER_PRESET_SLIDER_MARKS,
  ENCODER_PRESET_TOOLTIPS,
  VIDEO_BITRATE_DEFAULTS,
  VIDEO_BITRATE_SLIDER_MARKS,
  FRAMERATE_SLIDER_MARKS,
  FRAMERATE_DEFAULTS,
  FRAMERATE_TOOLTIPS,
} from '../../utils/config-constants';
import ToggleSwitch from './form-toggleswitch';

const { Panel } = Collapse;

interface VideoVariantFormProps {
  dataState: VideoVariant;
  onUpdateField: FieldUpdaterFunc;
}

export default function VideoVariantForm({
  dataState = DEFAULT_VARIANT_STATE,
  onUpdateField,
}: VideoVariantFormProps) {
  const videoPassthroughEnabled = dataState.videoPassthrough;

  const handleFramerateChange = (value: number) => {
    onUpdateField({ fieldName: 'framerate', value });
  };
  const handleVideoBitrateChange = (value: number) => {
    onUpdateField({ fieldName: 'videoBitrate', value });
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

  // Video passthrough handling
  const handleVideoPassConfirm = () => {
    onUpdateField({ fieldName: 'videoPassthrough', value: true });
  };
  // If passthrough is currently on, set it back to false on toggle.
  // Else let the Popconfirm turn it on.
  const handleVideoPassthroughToggle = (value: boolean) => {
    if (videoPassthroughEnabled) {
      onUpdateField({ fieldName: 'videoPassthrough', value });
    }
  };

  // Slider notes
  const selectedVideoBRnote = () => {
    if (videoPassthroughEnabled) {
      return 'Bitrate selection is disabled when Video Passthrough is enabled.';
    }
    let note = `${dataState.videoBitrate}${VIDEO_BITRATE_DEFAULTS.unit}`;
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
    if (videoPassthroughEnabled) {
      return 'Framerate selection is disabled when Video Passthrough is enabled.';
    }
    return FRAMERATE_TOOLTIPS[dataState.framerate] || '';
  };
  const cpuUsageNote = () => {
    if (videoPassthroughEnabled) {
      return 'CPU usage selection is disabled when Video Passthrough is enabled.';
    }
    return ENCODER_PRESET_TOOLTIPS[dataState.cpuUsageLevel] || '';
  };

  const classes = classNames({
    'config-variant-form': true,
    'video-passthrough-enabled': videoPassthroughEnabled,
  });
  return (
    <div className={classes}>
      <p className="description">
        <a href="https://owncast.online/docs/video" target="_blank" rel="noopener noreferrer">
          Learn more
        </a>{' '}
        about how each of these settings can impact the performance of your server.
      </p>

      {videoPassthroughEnabled && (
        <p className="passthrough-warning">
          NOTE: Video Passthrough for this output stream variant is <em>enabled</em>, disabling the
          below video encoding settings.
        </p>
      )}

      <Row gutter={16}>
        <Col sm={24} md={12}>
          {/* ENCODER PRESET (CPU USAGE) FIELD */}
          <div className="form-module cpu-usage-container">
            <Typography.Title level={3}>CPU Usage</Typography.Title>
            <p className="description">
              Reduce to improve server performance, or increase it to improve video quality.
            </p>
            <div className="segment-slider-container">
              <Slider
                tipFormatter={value => ENCODER_PRESET_TOOLTIPS[value]}
                onChange={handleVideoCpuUsageLevelChange}
                min={1}
                max={Object.keys(ENCODER_PRESET_SLIDER_MARKS).length}
                marks={ENCODER_PRESET_SLIDER_MARKS}
                defaultValue={dataState.cpuUsageLevel}
                value={dataState.cpuUsageLevel}
                disabled={dataState.videoPassthrough}
              />
              <p className="selected-value-note">{cpuUsageNote()}</p>
            </div>
            <p className="read-more-subtext">
              <a
                href="https://owncast.online/docs/video/#cpu-usage"
                target="_blank"
                rel="noopener noreferrer"
              >
                Read more about CPU usage.
              </a>
            </p>
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
            <p className="description">{VIDEO_BITRATE_DEFAULTS.tip}</p>
            <div className="segment-slider-container">
              <Slider
                tipFormatter={value => `${value} ${VIDEO_BITRATE_DEFAULTS.unit}`}
                disabled={dataState.videoPassthrough}
                defaultValue={dataState.videoBitrate}
                value={dataState.videoBitrate}
                onChange={handleVideoBitrateChange}
                step={VIDEO_BITRATE_DEFAULTS.incrementBy}
                min={VIDEO_BITRATE_DEFAULTS.min}
                max={VIDEO_BITRATE_DEFAULTS.max}
                marks={VIDEO_BITRATE_SLIDER_MARKS}
              />
              <p className="selected-value-note">{selectedVideoBRnote()}</p>
            </div>
            <p className="read-more-subtext">
              <a
                href="https://owncast.online/docs/video/#bitrate"
                target="_blank"
                rel="noopener noreferrer"
              >
                Read more about bitrates.
              </a>
            </p>
          </div>
        </Col>
      </Row>
      <Collapse className="advanced-settings">
        <Panel header="Advanced Settings" key="1">
          <Row gutter={16}>
            <Col sm={24} md={12}>
              <div className="form-module resolution-module">
                <Typography.Title level={3}>Resolution</Typography.Title>
                <p className="description">
                  Resizing your content will take additional resources on your server. If you wish
                  to optionally resize your content for this stream output then you should either
                  set the width <strong>or</strong> the height to keep your aspect ratio.{' '}
                  <a
                    href="https://owncast.online/docs/video/#resolution"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    Read more about resolutions.
                  </a>
                </p>
                <br />
                <TextField
                  type="number"
                  {...VIDEO_VARIANT_SETTING_DEFAULTS.scaledWidth}
                  value={dataState.scaledWidth}
                  onChange={handleScaledWidthChanged}
                  disabled={dataState.videoPassthrough}
                />
                <TextField
                  type="number"
                  {...VIDEO_VARIANT_SETTING_DEFAULTS.scaledHeight}
                  value={dataState.scaledHeight}
                  onChange={handleScaledHeightChanged}
                  disabled={dataState.videoPassthrough}
                />
              </div>
            </Col>
            <Col sm={24} md={12}>
              {/* VIDEO PASSTHROUGH FIELD */}
              <div className="form-module video-passthrough-module">
                <Typography.Title level={3}>Video Passthrough</Typography.Title>
                <p className="description">
                  <p>
                    Enabling video passthrough may allow for less hardware utilization, but may also
                    make your stream <strong>unplayable</strong>.
                  </p>
                  <p>
                    All other settings for this stream output will be disabled if passthrough is
                    used.
                  </p>
                  <p>
                    <a
                      href="https://owncast.online/docs/video/#video-passthrough"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      Read the documentation before enabling, as it impacts your stream.
                    </a>
                  </p>
                </p>
                <Popconfirm
                  disabled={dataState.videoPassthrough === true}
                  title="Did you read the documentation about video passthrough and understand the risks involved with enabling it?"
                  icon={<ExclamationCircleFilled />}
                  onConfirm={handleVideoPassConfirm}
                  okText="Yes"
                  cancelText="No"
                >
                  {/* adding an <a> tag to force Popcofirm to register click on toggle */}
                  {/* eslint-disable-next-line jsx-a11y/anchor-is-valid */}
                  <a href="#">
                    <ToggleSwitch
                      label="Use Video Passthrough?"
                      fieldName="video-passthrough"
                      tip={VIDEO_VARIANT_SETTING_DEFAULTS.videoPassthrough.tip}
                      checked={dataState.videoPassthrough}
                      onChange={handleVideoPassthroughToggle}
                    />
                  </a>
                </Popconfirm>
              </div>
            </Col>
          </Row>

          {/* FRAME RATE FIELD */}
          <div className="form-module frame-rate-module">
            <Typography.Title level={3}>Frame rate</Typography.Title>
            <p className="description">{FRAMERATE_DEFAULTS.tip}</p>
            <div className="segment-slider-container">
              <Slider
                tipFormatter={value => `${value} ${FRAMERATE_DEFAULTS.unit}`}
                defaultValue={dataState.framerate}
                value={dataState.framerate}
                onChange={handleFramerateChange}
                step={FRAMERATE_DEFAULTS.incrementBy}
                min={FRAMERATE_DEFAULTS.min}
                max={FRAMERATE_DEFAULTS.max}
                marks={FRAMERATE_SLIDER_MARKS}
                disabled={dataState.videoPassthrough}
              />
              <p className="selected-value-note">{selectedFramerateNote()}</p>
            </div>
            <p className="read-more-subtext">
              <a
                href="https://owncast.online/docs/video/#framerate"
                target="_blank"
                rel="noopener noreferrer"
              >
                Read more about framerates.
              </a>
            </p>
          </div>
        </Panel>
      </Collapse>
    </div>
  );
}

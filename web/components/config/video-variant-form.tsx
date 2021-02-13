// This content populates the video variant modal, which is spawned from the variants table.
import React from 'react';
import { Slider, Switch, Collapse, Typography } from 'antd';
import { FieldUpdaterFunc, VideoVariant, UpdateArgs } from '../../types/config-section';
import TextField from './form-textfield';
import { DEFAULT_VARIANT_STATE } from '../../utils/config-constants';
import InfoTip from '../info-tip';
import CPUUsageSelector from './cpu-usage';
// import ToggleSwitch from './form-toggleswitch-with-submit';

const { Panel } = Collapse;

const VIDEO_VARIANT_DEFAULTS = {
  framerate: {
    min: 10,
    max: 90,
    defaultValue: 24,
    unit: 'fps',
    incrementBy: 1,
    tip: 'You prob wont need to touch this unless youre a hardcore gamer and need all the bitties',
  },
  videoBitrate: {
    min: 600,
    max: 6000,
    defaultValue: 1200,
    unit: 'kbps',
    incrementBy: 100,
    tip: 'This is importatnt yo',
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

  const videoBitrateDefaults = VIDEO_VARIANT_DEFAULTS.videoBitrate;
  const videoBRMin = videoBitrateDefaults.min;
  const videoBRMax = videoBitrateDefaults.max;
  const videoBRUnit = videoBitrateDefaults.unit;

  const selectedVideoBRnote = `Selected: ${dataState.videoBitrate}${videoBRUnit} - it sucks`;
  const selectedFramerateNote = `Selected: ${dataState.framerate}${framerateUnit} - whoa there`;
  const selectedPresetNote = '';

  return (
    <div className="config-variant-form">
      <p className="description">
        Say a thing here about how this all works. Read more{' '}
        <a href="https://owncast.online/docs/configuration/">here</a>.
      </p>

      <div className="row">
        <div>
          {/* ENCODER PRESET FIELD */}
          <div className="form-module cpu-usage-container">
            <CPUUsageSelector
              defaultValue={dataState.cpuUsageLevel}
              onChange={handleVideoCpuUsageLevelChange}
            />
            {selectedPresetNote && (
              <span className="selected-value-note">{selectedPresetNote}</span>
            )}
          </div>

          {/* VIDEO PASSTHROUGH FIELD */}
          <div style={{ display: 'none' }} className="form-module">
            <p className="label">
              <InfoTip tip={VIDEO_VARIANT_DEFAULTS.videoPassthrough.tip} />
              Use Video Passthrough?
            </p>
            <div className="form-component">
              {/* todo: change to ToggleSwitch for layout */}
              <Switch
                defaultChecked={dataState.videoPassthrough}
                checked={dataState.videoPassthrough}
                onChange={handleVideoPassChange}
                // label="Use Video Passthrough"
                checkedChildren="Yes"
                unCheckedChildren="No"
              />
            </div>
          </div>

          {/* VIDEO BITRATE FIELD */}
          <div className={`form-module ${dataState.videoPassthrough ? 'disabled' : ''}`}>
            <Typography.Title level={3} className="section-title">
              Video Bitrate
            </Typography.Title>
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
                marks={{
                  [videoBRMin]: `${videoBRMin} ${videoBRUnit}`,
                  [videoBRMax]: `${videoBRMax} ${videoBRUnit}`,
                }}
              />
              {selectedVideoBRnote && (
                <span className="selected-value-note">{selectedVideoBRnote}</span>
              )}
            </div>
          </div>
        </div>
        <Collapse className="advanced-settings">
          <Panel header="Advanced Settings" key="1">
            <div className="section-intro">
              Resizing your content will take additional resources on your server. If you wish to
              optionally resize your output for this stream variant then you should either set the
              width <strong>or</strong> the height to keep your aspect ratio.
            </div>
            <div className="field">
              <TextField
                type="number"
                {...VIDEO_VARIANT_DEFAULTS.scaledWidth}
                value={dataState.scaledWidth}
                onChange={handleScaledWidthChanged}
              />
            </div>
            <div className="field">
              <TextField
                type="number"
                {...VIDEO_VARIANT_DEFAULTS.scaledHeight}
                value={dataState.scaledHeight}
                onChange={handleScaledHeightChanged}
              />
            </div>

            {/* FRAME RATE FIELD */}
            <div className="field">
              <p className="label">
                <InfoTip tip={VIDEO_VARIANT_DEFAULTS.framerate.tip} />
                Frame rate:
              </p>
              <div className="segment-slider-container form-component">
                <Slider
                  // tooltipVisible
                  tipFormatter={value => `${value} ${framerateUnit}`}
                  defaultValue={dataState.framerate}
                  value={dataState.framerate}
                  onChange={handleFramerateChange}
                  step={framerateDefaults.incrementBy}
                  min={framerateMin}
                  max={framerateMax}
                  marks={{
                    [framerateMin]: `${framerateMin} ${framerateUnit}`,
                    [framerateMax]: `${framerateMax} ${framerateUnit}`,
                  }}
                />
                {selectedFramerateNote ? (
                  <span className="selected-value-note">{selectedFramerateNote}</span>
                ) : null}
              </div>
            </div>
          </Panel>
        </Collapse>
      </div>
    </div>
  );
}

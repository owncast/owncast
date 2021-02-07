// This content populates the video variant modal, which is spawned from the variants table.
import React from 'react';
import { Slider, Switch, Collapse } from 'antd';
import { FieldUpdaterFunc, VideoVariant } from '../../types/config-section';
import { DEFAULT_VARIANT_STATE } from '../../utils/config-constants';
import InfoTip from '../info-tip';
import CPUUsageSelector from './cpu-usage';

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
      <div className="section-intro">
        Say a thing here about how this all works. Read more{' '}
        <a href="https://owncast.online/docs/configuration/">here</a>.
        <br />
        <br />
      </div>

      {/* ENCODER PRESET FIELD */}
      <div className="field">
        <div className="form-component">
          <CPUUsageSelector
            defaultValue={dataState.cpuUsageLevel}
            onChange={handleVideoCpuUsageLevelChange}
          />
          {selectedPresetNote ? (
            <span className="selected-value-note">{selectedPresetNote}</span>
          ) : null}
        </div>
      </div>

      {/* VIDEO PASSTHROUGH FIELD */}
      <div style={{ display: 'none' }}>
        <div className="field">
          <p className="label">
            <InfoTip tip={VIDEO_VARIANT_DEFAULTS.videoPassthrough.tip} />
            Use Video Passthrough?
          </p>
          <div className="form-component">
            <Switch
              defaultChecked={dataState.videoPassthrough}
              checked={dataState.videoPassthrough}
              onChange={handleVideoPassChange}
              checkedChildren="Yes"
              unCheckedChildren="No"
            />
          </div>
        </div>
      </div>

      {/* VIDEO BITRATE FIELD */}
      <div className={`field ${dataState.videoPassthrough ? 'disabled' : ''}`}>
        <p className="label">
          <InfoTip tip={VIDEO_VARIANT_DEFAULTS.videoBitrate.tip} />
          Video Bitrate:
        </p>
        <div className="form-component">
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
          {selectedVideoBRnote ? (
            <span className="selected-value-note">{selectedVideoBRnote}</span>
          ) : null}
        </div>
      </div>

      <br />
      <br />
      <br />
      <br />

      <Collapse>
        <Panel header="Advanced Settings" key="1">
          <div className="section-intro">Touch if you dare.</div>

          {/* FRAME RATE FIELD */}
          <div className="field">
            <p className="label">
              <InfoTip tip={VIDEO_VARIANT_DEFAULTS.framerate.tip} />
              Frame rate:
            </p>
            <div className="form-component">
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
  );
}

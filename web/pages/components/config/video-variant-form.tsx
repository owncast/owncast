import React from 'react';
import { Slider, Select, Switch, Divider, Collapse } from 'antd';
import { PRESETS, VideoVariant } from '../../../types/config-section';
import { ENCODER_PRESETS, DEFAULT_VARIANT_STATE } from './constants';
import InfoTip from '../info-tip';

const { Option } = Select;
const { Panel } = Collapse;

const VIDEO_VARIANT_DEFAULTS = {
  framerate: {
    min: 10,
    max: 80,
    defaultValue: 24,
    unit: 'fps',
    incrementBy: 1,
    tip: 'You prob wont need to touch this unless youre a hardcore gamer and need all the bitties',
  },
  videoBitrate: {
    min: 600,
    max: 1200,
    defaultValue: 800,
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
    tip: 'nothing to see here'
  },
  encoderPreset: {
    defaultValue: ENCODER_PRESETS[2],
    tip: 'Info and stuff.'
  },
  videoPassthrough: {
    tip: 'If No is selected, then you should set your desired Video Bitrate.'
  },
  audioPassthrough: {
    tip: 'If No is selected, then you should set your desired Audio Bitrate.'
  },
};

interface VideoVariantFormProps {
  dataState: VideoVariant;
  onUpdateField: () => void;
}

export default function VideoVariantForm({ dataState = DEFAULT_VARIANT_STATE, onUpdateField }: VideoVariantFormProps) {
  // const [dataState, setDataState] = useState(initialValues);

  const handleFramerateChange = (value: number) => {
    onUpdateField({ fieldname: 'framerate', value });
  };
  const handleVideoBitrateChange = (value: number) => {
    onUpdateField({ fieldname: 'videoBitrate', value });
  };
  const handleAudioBitrateChange = (value: number) => {
    onUpdateField({ fieldname: 'audioBitrate', value });
  };
  const handleEncoderPresetChange = (value: PRESETS) => {
    onUpdateField({ fieldname: 'encoderPreset', value });
  };
  const handleAudioPassChange = (value: boolean) => {
    onUpdateField({ fieldname: 'audioPassthrough', value });
  };
  const handleVideoPassChange = (value: boolean) => {
    onUpdateField({ fieldname: 'videoPassthrough', value });
  };


  const framerateDefaults = VIDEO_VARIANT_DEFAULTS.framerate;
  const framerateMin = framerateDefaults.min;
  const framerateMax = framerateDefaults.max;
  const framerateUnit = framerateDefaults.unit;

  const encoderDefaults = VIDEO_VARIANT_DEFAULTS.encoderPreset;
  const videoBitrateDefaults = VIDEO_VARIANT_DEFAULTS.videoBitrate;
  const videoBRMin = videoBitrateDefaults.min;
  const videoBRMax = videoBitrateDefaults.max;
  const videoBRUnit = videoBitrateDefaults.unit;

  const audioBitrateDefaults = VIDEO_VARIANT_DEFAULTS.audioBitrate;
  const audioBRMin = audioBitrateDefaults.min;
  const audioBRMax = audioBitrateDefaults.max;
  const audioBRUnit = audioBitrateDefaults.unit;

  const selectedVideoBRnote = `Selected: ${dataState.videoBitrate}${videoBRUnit} - it sucks`;
  const selectedAudioBRnote = `Selected: ${dataState.audioBitrate}${audioBRUnit} - too slow`;
  const selectedFramerateNote = `Selected: ${dataState.framerate}${framerateUnit} - whoa there`;
  const selectedPresetNote = '';

  return (
    <div className="variant-form">
      <div className="section-intro">
        Say a thing here about how this all works.

        Read more <a href="https://owncast.online/docs/configuration/">here</a>.
        <br /><br />
      </div>

      {/* ENCODER PRESET FIELD */}
      <div className="field">
        <p className="label">
          <InfoTip tip={encoderDefaults.tip} />
          Encoder Preset:
        </p>
        <div className="form-component">
          <Select defaultValue={encoderDefaults.defaultValue} style={{ width: 200 }} onChange={handleEncoderPresetChange}>
            {
              ENCODER_PRESETS.map(preset => (
                <Option
                  key={`option-${preset}`}
                  value={preset}
                >{preset}</Option>
              ))
            }
          </Select>
          {selectedPresetNote ? <span className="selected-value-note">{selectedPresetNote}</span> : null }
        </div>
      </div>



      {/* VIDEO PASSTHROUGH FIELD */}
      <div style={{ display: 'none'}}>
        <div className="field">
          <p className="label">
            <InfoTip tip={VIDEO_VARIANT_DEFAULTS.videoPassthrough.tip} />
            Use Video Passthrough?
          </p>
          <div className="form-component">
            <Switch
              defaultChecked={dataState.videoPassthrough}
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
            // tooltipVisible={dataState.videoPassthrough !== true}
            tipFormatter={value => `${value} ${videoBRUnit}`}
            disabled={dataState.videoPassthrough === true}
            defaultValue={dataState.videoBitrate}
            onChange={handleVideoBitrateChange}
            step={videoBitrateDefaults.incrementBy}
            min={videoBRMin}
            max={videoBRMax}
            marks={{
              [videoBRMin]: `${videoBRMin} ${videoBRUnit}`,
              [videoBRMax]: `${videoBRMax} ${videoBRUnit}`,
            }}
          />
          {selectedVideoBRnote ? <span className="selected-value-note">{selectedVideoBRnote}</span> : null }

        </div>
      </div>


      <br /><br /><br /><br />
      
      <Collapse>
        <Panel header="Advanced Settings" key="1">
          <div className="section-intro">
            Touch if you dare.
          </div>

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
                onChange={handleFramerateChange}
                step={framerateDefaults.incrementBy}
                min={framerateMin}
                max={framerateMax}
                marks={{
                  [framerateMin]: `${framerateMin} ${framerateUnit}`,
                  [framerateMax]: `${framerateMax} ${framerateUnit}`,
                }}
              />
              {selectedFramerateNote ? <span className="selected-value-note">{selectedFramerateNote}</span> : null }

            </div>
          </div>

          <Divider />

          {/* AUDIO PASSTHROUGH FIELD */}
          <div className="field">
            <p className="label">
              <InfoTip tip={VIDEO_VARIANT_DEFAULTS.audioPassthrough.tip} />
              Use Audio Passthrough?
            </p>
            <div className="form-component">
              <Switch
                defaultChecked={dataState.audioPassthrough}
                onChange={handleAudioPassChange}
                checkedChildren="Yes"
                unCheckedChildren="No"
              />
              {dataState.audioPassthrough ? <span className="note">Same as source</span>: null}
            </div>
          </div>

          {/* AUDIO BITRATE FIELD */}
          <div className={`field ${dataState.audioPassthrough ? 'disabled' : ''}`}>
            <p className="label">
              <InfoTip tip={VIDEO_VARIANT_DEFAULTS.audioBitrate.tip} />
              Audio Bitrate:
            </p>
            <div className="form-component">
              <Slider
                // tooltipVisible={dataState.audioPassthrough !== true}
                tipFormatter={value => `${value} ${audioBRUnit}`}
                disabled={dataState.audioPassthrough === true}
                defaultValue={dataState.audioBitrate}
                onChange={handleAudioBitrateChange}
                step={audioBitrateDefaults.incrementBy}
                min={audioBRMin}
                max={audioBRMax}
                marks={{
                  [audioBRMin]: `${audioBRMin} ${audioBRUnit}`,
                  [audioBRMax]: `${audioBRMax} ${audioBRUnit}`,
                }}
              />

              {selectedAudioBRnote ? <span className="selected-value-note">{selectedAudioBRnote}</span> : null }
            </div>
          </div>


        </Panel>
      </Collapse>

    </div>
  );

}
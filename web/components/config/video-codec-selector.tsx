import React, { useContext, useState, useEffect } from 'react';
import { Typography, Select } from 'antd';
import { ServerStatusContext } from '../../utils/server-status-context';
import { AlertMessageContext } from '../../utils/alert-message-context';
import {
  API_VIDEO_CODEC,
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


export default function CodecSelector() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { videoCodec, supportedCodecs } = serverConfig || {};
  const { Title } = Typography;
  const { Option } = Select;
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  const { setMessage } = useContext(AlertMessageContext);
  const [selectedCodec, setSelectedCodec] = useState(videoCodec);

  let resetTimer = null;

  useEffect(() => {
    setSelectedCodec(videoCodec);
  }, [videoCodec]);

  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  async function handleChange(value) {
    console.log(`selected ${value}`);
    setSelectedCodec(value);

    await postConfigUpdateToAPI({
      apiPath: API_VIDEO_CODEC,
      data: { value: value },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'videoCodec',
          value: value,
          path: 'videoSettings',
        });
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Video codec updated.'));

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
  }

  const items = supportedCodecs.map(function (codec) {
    var title = codec;
    if (title === 'libx264') {
      title = 'Default (libx264)';
    } else if (title === 'h264_nvenc') {
      title = 'NVIDIA GPU acceleration';
    } else if (title === 'h264_vaapi') {
      title = 'VAAPI hardware encoding';
    } else if (title === 'h264_qsv') {
      title = 'Intel QuickSync';
    } else if (title === 'h264_v4l2m2m') {
      title = 'Video4Linux hardware encoding';
    }

    return (
      <Option key={codec} value={codec}>
        {title}
      </Option>
    );
  })

  var description = '';
  if (selectedCodec === 'libx264') {
    description = 'libx264 is default codec and generally the choice you will want to use unless you have access to more specialized options. It is also likely the only option for running on shared VPS environments.';
  } else if (selectedCodec === 'h264_nvenc') {
    description = 'You can use your NVIDIA GPU for encoding if you have a modern NVIDIA card with encoding cores.';
  } else if (selectedCodec === 'h264_vaapi') {
    description = 'VAAPI may be supported by NVIDIA proprietary drivers, Mesa open-source drivers for AMD or Intel graphics cards.';
  } else if (selectedCodec === 'h264_qsv') {
    description = "Quick Sync Video is Intel's brand for its dedicated video encoding and decoding hardware core. May be an option if you have a modern Intel CPU with integrated graphics.";
  } else if (selectedCodec === 'h264_v4l2m2m') {
    description = "Video4Linux is an interface to multiple different hardware encoding platforms such as Intel and AMD."
  }

  return (
    <>
      <Title level={3} className="section-title">
        Video Codec
      </Title>
      <p className="description">
        Blurb about codecs go here.
      </p>
      <div className="segment-slider-container">
        <Select defaultValue={selectedCodec} value={selectedCodec} style={{ width: '100%' }} onChange={handleChange} >
          {items}
        </Select>
        <FormStatusIndicator status={submitStatus} />
        <p className="selected-value-note">
          {description}
        </p>
      </div>
    </>
  );
}

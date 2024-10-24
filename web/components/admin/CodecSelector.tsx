import { Popconfirm, Select, Typography } from 'antd';
import React, { FC, useContext, useEffect, useState } from 'react';
import { AlertMessageContext } from '../../utils/alert-message-context';
import {
  API_VIDEO_CODEC,
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
} from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import { ServerStatusContext } from '../../utils/server-status-context';
import { FormStatusIndicator } from './FormStatusIndicator';

export type CodecSelectorProps = {};

export const CodecSelector: FC<CodecSelectorProps> = () => {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { videoCodec, supportedCodecs } = serverConfig || {};
  const { Title } = Typography;
  const { Option } = Select;
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  const { setMessage } = useContext(AlertMessageContext);
  const [selectedCodec, setSelectedCodec] = useState(videoCodec);
  const [pendingSaveCodec, setPendingSavecodec] = useState(videoCodec);
  const [confirmPopupOpen, setConfirmPopupOpen] = React.useState(false);

  let resetTimer = null;

  useEffect(() => {
    setSelectedCodec(videoCodec);
  }, [videoCodec]);

  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  function handleChange(value) {
    setPendingSavecodec(value);
    setConfirmPopupOpen(true);
  }

  async function save() {
    setSelectedCodec(pendingSaveCodec);
    setPendingSavecodec('');
    setConfirmPopupOpen(false);

    await postConfigUpdateToAPI({
      apiPath: API_VIDEO_CODEC,
      data: { value: pendingSaveCodec },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'videoCodec',
          value: pendingSaveCodec,
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

  const items = supportedCodecs.map(codec => {
    let title = codec;
    if (title === 'libx264') {
      title = 'Default (libx264)';
    } else if (title === 'h264_nvenc') {
      title = 'NVIDIA GPU acceleration';
    } else if (title === 'h264_vaapi') {
      title = 'VA-API hardware encoding';
    } else if (title === 'h264_qsv') {
      title = 'Intel QuickSync';
    } else if (title === 'h264_v4l2m2m') {
      title = 'Video4Linux hardware encoding';
    } else if (title === 'h264_omx') {
      title = 'OpenMax (omx) for Raspberry Pi';
    } else if (title === 'h264_videotoolbox') {
      title = 'Apple VideoToolbox (hardware)';
    } else if (title === 'h264_rkmpp') {
      title = 'Rockchip MPP (hardware)';
    }

    return (
      <Option key={codec} value={codec}>
        {title}
      </Option>
    );
  });

  let description = '';
  if (selectedCodec === 'libx264') {
    description =
      'libx264 is the default codec and generally the only working choice for shared VPS environments. This is likely what you should be using unless you know you have set up other options.';
  } else if (selectedCodec === 'h264_nvenc') {
    description =
      'You can use your NVIDIA GPU for encoding if you have a modern NVIDIA card with encoding cores.';
  } else if (selectedCodec === 'h264_vaapi') {
    description =
      'VA-API may be supported by your NVIDIA proprietary drivers, Mesa open-source drivers for AMD or Intel graphics.';
  } else if (selectedCodec === 'h264_qsv') {
    description =
      "Quick Sync Video is Intel's brand for its dedicated video encoding and decoding hardware. It may be an option if you have a modern Intel CPU with integrated graphics.";
  } else if (selectedCodec === 'h264_v4l2m2m') {
    description =
      'Video4Linux is an interface to multiple different hardware encoding platforms such as Intel and AMD.';
  } else if (selectedCodec === 'h264_omx') {
    description = 'OpenMax is a codec most often used with a Raspberry Pi.';
  } else if (selectedCodec === 'h264_videotoolbox') {
    description =
      'Apple VideoToolbox is a low-level framework that provides direct access to hardware encoders and decoders.';
  } else if (selectedCodec === 'h264_rkmpp') {
    description =
      'Rockchip MPP is a hardware-accelerate encoder for rockchip ARM SoC chips, like orange pi etc.';
  }

  return (
    <>
      <Title level={3} className="section-title">
        Video Codec
      </Title>
      <div className="description">
        If you have access to specific hardware with the drivers and software installed for them,
        you may be able to improve your video encoding performance.
        <p>
          <a
            href="https://owncast.online/docs/codecs?source=admin"
            target="_blank"
            rel="noopener noreferrer"
          >
            Read the documentation about this setting before changing it or you may make your stream
            unplayable.
          </a>
        </p>
      </div>
      <div className="segment-slider-container">
        <Popconfirm
          title={`Are you sure you want to change your video codec to ${pendingSaveCodec} and understand what this means?`}
          open={confirmPopupOpen}
          placement="leftBottom"
          onConfirm={save}
          onCancel={() => setConfirmPopupOpen(false)}
          okText="Yes"
          cancelText="No"
        >
          <Select
            defaultValue={selectedCodec}
            value={selectedCodec}
            style={{ width: '100%' }}
            onChange={handleChange}
          >
            {items}
          </Select>
        </Popconfirm>
        <FormStatusIndicator status={submitStatus} />

        <p id="selected-codec-note" className="selected-value-note">
          {description}
        </p>
      </div>
    </>
  );
};

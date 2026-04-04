import { Popconfirm, Select, Typography } from 'antd';
import React, { FC, useContext, useEffect, useState } from 'react';
import { useTranslation } from 'next-export-i18n';
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
import { Localization } from '../../types/localization';
import { FormStatusIndicator } from './FormStatusIndicator';

export type CodecSelectorProps = {};

export const CodecSelector: FC<CodecSelectorProps> = () => {
  const { t } = useTranslation();
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { videoEncoder, supportedVideoEncoders } = serverConfig || {};
  const { Title } = Typography;
  const { Option } = Select;
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  const { setMessage } = useContext(AlertMessageContext);
  const [selectedCodec, setSelectedCodec] = useState(videoEncoder);
  const [pendingSaveCodec, setPendingSaveCodec] = useState(videoEncoder);
  const [confirmPopupOpen, setConfirmPopupOpen] = React.useState(false);

  let resetTimer = null;

  useEffect(() => {
    setSelectedCodec(videoEncoder);
  }, [videoEncoder]);

  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  function handleChange(value) {
    setPendingSaveCodec(value);
    setConfirmPopupOpen(true);
  }

  async function save() {
    setSelectedCodec(pendingSaveCodec);
    setPendingSaveCodec('');
    setConfirmPopupOpen(false);

    await postConfigUpdateToAPI({
      apiPath: API_VIDEO_CODEC,
      data: { value: pendingSaveCodec },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'videoEncoder',
          value: pendingSaveCodec,
        });
        setSubmitStatus(
          createInputStatus(STATUS_SUCCESS, t(Localization.Admin.StatusMessages.videoCodecUpdated)),
        );

        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
        if (serverStatusData.online) {
          setMessage(
            'Your encoder setting will take effect the next time you begin a live stream.',
          );
        }
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));

        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  }

  const items = (supportedVideoEncoders || []).map(encoder => (
    <Option key={encoder.encoderType} value={encoder.encoderType}>
      {encoder.encoderDisplayName}
    </Option>
  ));

  let description = '';
  if (selectedCodec === 'software') {
    description =
      'Software encoding uses your CPU. This is the default and generally the only working choice for shared VPS environments. This is likely what you should be using unless you know you have set up other options.';
  } else if (selectedCodec === 'nvenc') {
    description =
      'You can use your NVIDIA GPU for encoding if you have a modern NVIDIA card with encoding cores.';
  } else if (selectedCodec === 'vaapi') {
    description =
      'VA-API may be supported by your NVIDIA proprietary drivers, Mesa open-source drivers for AMD or Intel graphics.';
  } else if (selectedCodec === 'qsv') {
    description =
      "Quick Sync Video is Intel's brand for its dedicated video encoding and decoding hardware. It may be an option if you have a modern Intel CPU with integrated graphics.";
  } else if (selectedCodec === 'v4l2m2m') {
    description =
      'Video4Linux is an interface to multiple different hardware encoding platforms such as Intel and AMD.';
  } else if (selectedCodec === 'omx') {
    description = 'OpenMax is an encoder most often used with a Raspberry Pi.';
  } else if (selectedCodec === 'videotoolbox') {
    description =
      'Apple VideoToolbox is a low-level framework that provides direct access to hardware encoders and decoders.';
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
          title={`Are you sure you want to change your video encoder to ${pendingSaveCodec} and understand what this means?`}
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

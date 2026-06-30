import { Select, Typography } from 'antd';
import React, { FC, useContext, useEffect, useState } from 'react';
import { API_AUTOPLAY, postConfigUpdateToAPI, RESET_TIMEOUT } from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import { ServerStatusContext } from '../../utils/server-status-context';
import { FormStatusIndicator } from './FormStatusIndicator';

// The three viewer-facing autoplay behaviors. `value` is what's persisted and
// what the player maps to a video.js autoplay option (off -> false,
// always -> 'any', sound-only -> 'play'). `description` is shown live under the
// dropdown and updates as the selection changes.
const AUTOPLAY_OPTIONS = [
  {
    value: 'off',
    label: 'Never',
    description:
      'Viewers press play to start watching. Uses the least bandwidth, since video only loads when someone chooses to watch.',
  },
  {
    value: 'always',
    label: 'Always',
    description:
      'The stream starts on its own the moment someone opens the page. It plays with sound where the browser allows it, and starts muted with a large one-tap unmute button everywhere else, so it always plays.',
  },
  {
    value: 'sound-only',
    label: 'Only with sound',
    description:
      "The stream starts on its own with sound, but only when the viewer's browser allows it. Otherwise the viewer presses play. It never starts silently.",
  },
];

export type AutoplaySelectorProps = {};

export const AutoplaySelector: FC<AutoplaySelectorProps> = () => {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { instanceDetails } = serverConfig || {};
  const { autoplay } = instanceDetails || {};
  const { Title } = Typography;
  const { Option } = Select;
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  const [selectedAutoplay, setSelectedAutoplay] = useState(autoplay);

  let resetTimer = null;

  useEffect(() => {
    setSelectedAutoplay(autoplay);
  }, [autoplay]);

  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  async function save(value: string) {
    setSelectedAutoplay(value);
    await postConfigUpdateToAPI({
      apiPath: API_AUTOPLAY,
      data: { value },
      onSuccess: () => {
        setFieldInConfigState({ fieldName: 'autoplay', value, path: 'instanceDetails' });
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Autoplay setting updated.'));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  }

  const description =
    AUTOPLAY_OPTIONS.find(option => option.value === selectedAutoplay)?.description || '';

  return (
    <>
      <Title level={3} className="section-title">
        Autoplay
      </Title>
      <div className="description">
        Choose whether the stream starts playing on its own when a viewer opens your page.
      </div>
      <div className="segment-slider-container">
        <Select
          defaultValue={selectedAutoplay}
          value={selectedAutoplay}
          style={{ width: '100%' }}
          onChange={save}
        >
          {AUTOPLAY_OPTIONS.map(option => (
            <Option key={option.value} value={option.value}>
              {option.label}
            </Option>
          ))}
        </Select>
        <FormStatusIndicator status={submitStatus} />
        <p className="selected-value-note">{description}</p>
      </div>
    </>
  );
};

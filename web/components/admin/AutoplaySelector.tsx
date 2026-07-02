import { Select, Typography } from 'antd';
import React, { FC, useContext, useEffect, useState } from 'react';
import { useTranslation } from 'next-export-i18n';
import { API_AUTOPLAY, postConfigUpdateToAPI, RESET_TIMEOUT } from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import { ServerStatusContext } from '../../utils/server-status-context';
import { Localization } from '../../types/localization';
import { Translation } from '../ui/Translation/Translation';
import { FormStatusIndicator } from './FormStatusIndicator';

// The three viewer-facing autoplay behaviors. `value` is what's persisted and
// what the player maps to a video.js autoplay option (off -> false,
// always -> 'any', sound-only -> 'play'). The description is shown live under
// the dropdown and updates as the selection changes.
const AUTOPLAY_OPTIONS = [
  {
    value: 'off',
    labelKey: Localization.Admin.Autoplay.optionOffLabel,
    labelDefault: 'Never',
    descriptionKey: Localization.Admin.Autoplay.optionOffDescription,
    descriptionDefault: 'Viewers will press play to start watching.',
  },
  {
    value: 'always',
    labelKey: Localization.Admin.Autoplay.optionAlwaysLabel,
    labelDefault: 'Always',
    descriptionKey: Localization.Admin.Autoplay.optionAlwaysDescription,
    descriptionDefault:
      'The stream always starts on its own the moment someone opens the page. It plays with sound where the browser allows it, and starts muted everywhere else.',
  },
  {
    value: 'sound-only',
    labelKey: Localization.Admin.Autoplay.optionSoundOnlyLabel,
    labelDefault: 'Only if sound is available',
    descriptionKey: Localization.Admin.Autoplay.optionSoundOnlyDescription,
    descriptionDefault:
      "The always stream starts on its own, but only when the viewer's browser allows sound to play as well. Otherwise the viewer presses play. It never starts silently. Different browsers handle this differently.",
  },
];

export type AutoplaySelectorProps = {};

export const AutoplaySelector: FC<AutoplaySelectorProps> = () => {
  const { t } = useTranslation();
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
        setSubmitStatus(
          createInputStatus(STATUS_SUCCESS, t(Localization.Admin.StatusMessages.autoplayUpdated)),
        );
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  }

  const selectedOption = AUTOPLAY_OPTIONS.find(option => option.value === selectedAutoplay);

  return (
    <>
      <Title level={3} className="section-title">
        <Translation translationKey={Localization.Admin.Autoplay.title} defaultText="Autoplay" />
      </Title>
      <div className="description">
        <Translation
          translationKey={Localization.Admin.Autoplay.description}
          defaultText="Choose whether the stream starts playing on its own when a viewer opens your page."
        />
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
              <Translation translationKey={option.labelKey} defaultText={option.labelDefault} />
            </Option>
          ))}
        </Select>
        <FormStatusIndicator status={submitStatus} />
        <p className="selected-value-note">
          {selectedOption && (
            <Translation
              translationKey={selectedOption.descriptionKey}
              defaultText={selectedOption.descriptionDefault}
            />
          )}
        </p>
      </div>
    </>
  );
};

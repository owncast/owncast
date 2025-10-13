import React, { CSSProperties, FC, useState } from 'react';
import { useRecoilValue } from 'recoil';
import { Input, Button, Select, Form } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { MessageType } from '../../../interfaces/socket-events';
import WebsocketService from '../../../services/websocket-service';
import { websocketServiceAtom, currentUserAtom } from '../../stores/ClientConfigStore';
import { validateDisplayName } from '../../../utils/displayNameValidation';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';
import styles from './NameChangeModal.module.scss';

const { Option } = Select;

export type UserColorProps = {
  color: number;
};

const UserColor: FC<UserColorProps> = ({ color }) => {
  const style: CSSProperties = {
    textAlign: 'center',
    backgroundColor: `var(--theme-color-users-${color})`,
    width: '100%',
    height: '100%',
  };
  return <div style={style} />;
};

type NameChangeModalProps = {
  closeModal: () => void;
};

export const NameChangeModal: FC<NameChangeModalProps> = ({ closeModal }) => {
  const { t } = useTranslation();
  const currentUser = useRecoilValue(currentUserAtom);
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
  const [newName, setNewName] = useState<string>(currentUser?.displayName || '');

  const characterLimit = 30;

  if (!currentUser) {
    return null;
  }

  const { displayName, displayColor } = currentUser;

  const saveEnabled = () => {
    if (!newName || !displayName) return false;
    const validation = validateDisplayName(newName, displayName, characterLimit);
    return validation.isValid && websocketService?.isConnected();
  };

  const handleNameChange = () => {
    if (!newName || !displayName) return;
    const validation = validateDisplayName(newName, displayName, characterLimit);

    if (!validation.isValid || !websocketService?.isConnected()) {
      return;
    }

    const nameChange = {
      type: MessageType.NAME_CHANGE,
      newName: validation.trimmedName,
    };
    websocketService.send(nameChange);
    closeModal();
  };

  const handleColorChange = (color: string) => {
    const colorChange = {
      type: MessageType.COLOR_CHANGE,
      newColor: Number(color),
    };
    websocketService.send(colorChange);
  };

  const showCount = info =>
    info.count > characterLimit ? (
      <Translation
        translationKey={Localization.Frontend.NameChangeModal.overLimit}
        defaultText="Over limit"
      />
    ) : (
      ''
    );

  const maxColor = 8; // 0...n
  const colorOptions = [...Array(maxColor)].map((_, i) => i);

  const placeholderText =
    t(Localization.Frontend.NameChangeModal.placeholder) || 'Your chat display name';

  const validation = validateDisplayName(newName, displayName, characterLimit);

  const saveButton = (
    <Button
      type="primary"
      id="name-change-submit"
      onClick={handleNameChange}
      disabled={!saveEnabled()}
    >
      <Translation
        translationKey={Localization.Frontend.NameChangeModal.buttonText}
        defaultText="Change name"
      />
    </Button>
  );

  return (
    <div>
      <div id="owncast-name-change-description-text">
        <Translation
          translationKey={Localization.Frontend.NameChangeModal.description}
          defaultText="Your chat display name is what people see when you send chat messages."
        />
      </div>
      <Form onSubmitCapture={handleNameChange} className={styles.form}>
        <Input.Search
          enterButton={saveButton}
          id="name-change-field"
          value={newName}
          onChange={e => setNewName(e.target.value)}
          placeholder={placeholderText}
          aria-label={placeholderText}
          showCount={{ formatter: showCount }}
          defaultValue={displayName}
          className={styles.inputGroup}
        />
        {!validation.isValid && validation.errorMessage && (
          <div className={styles.error}>{validation.errorMessage}</div>
        )}
      </Form>
      <Form.Item
        label={
          <Translation
            translationKey={Localization.Frontend.NameChangeModal.colorLabel}
            defaultText="Your Color"
          />
        }
        className={styles.colorChange}
      >
        <Select
          style={{ width: 120 }}
          onChange={handleColorChange}
          defaultValue={displayColor.toString()}
          className={styles.colorDropdown}
        >
          {colorOptions.map(e => (
            <Option key={e.toString()} title={e}>
              <UserColor color={e} aria-label={e.toString()} />
            </Option>
          ))}
        </Select>
      </Form.Item>
      <div id="owncast-name-change-auth-info-text">
        <Translation
          translationKey={Localization.Frontend.NameChangeModal.authInfo}
          defaultText='You can also authenticate an IndieAuth or Fediverse account via the "Authenticate" menu.'
        />
      </div>
    </div>
  );
};

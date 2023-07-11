import React, { CSSProperties, FC, useState } from 'react';
import { useRecoilValue } from 'recoil';
import { Input, Button, Select, Form } from 'antd';
import { MessageType } from '../../../interfaces/socket-events';
import WebsocketService from '../../../services/websocket-service';
import { websocketServiceAtom, currentUserAtom } from '../../stores/ClientConfigStore';
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
  const currentUser = useRecoilValue(currentUserAtom);
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
  const [newName, setNewName] = useState<string>(currentUser?.displayName);

  const characterLimit = 30;

  if (!currentUser) {
    return null;
  }

  const { displayName, displayColor } = currentUser;

  const saveEnabled = () => {
    const count = newName !== undefined ? Array.from(newName).length : 0;
    return (
      newName !== displayName &&
      count > 0 &&
      count <= characterLimit &&
      websocketService?.isConnected()
    );
  };

  const handleNameChange = () => {
    if (!saveEnabled()) return;

    const nameChange = {
      type: MessageType.NAME_CHANGE,
      newName,
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

  const showCount = info => (info.count > characterLimit ? 'Over limit' : '');

  const maxColor = 8; // 0...n
  const colorOptions = [...Array(maxColor)].map((_, i) => i);

  const saveButton = (
    <Button
      type="primary"
      id="name-change-submit"
      onClick={handleNameChange}
      disabled={!saveEnabled()}
    >
      Change name
    </Button>
  );

  return (
    <div>
      Your chat display name is what people see when you send chat messages.
      <Form onSubmitCapture={handleNameChange} className={styles.form}>
        <Input.Search
          enterButton={saveButton}
          id="name-change-field"
          value={newName}
          onChange={e => setNewName(e.target.value)}
          placeholder="Your chat display name"
          aria-label="Your chat display name"
          showCount={{ formatter: showCount }}
          defaultValue={displayName}
          className={styles.inputGroup}
        />
      </Form>
      <Form.Item label="Your Color" className={styles.colorChange}>
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
      You can also authenticate an IndieAuth or Fediverse account via the &quot;Authenticate&quot;
      menu.
    </div>
  );
};

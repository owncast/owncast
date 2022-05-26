import { useState } from 'react';
import { useRecoilValue } from 'recoil';
import { Input, Button } from 'antd';
import { MessageType } from '../../interfaces/socket-events';
import WebsocketService from '../../services/websocket-service';
import { websocketServiceAtom, chatDisplayNameAtom } from '../stores/ClientConfigStore';

/* eslint-disable @typescript-eslint/no-unused-vars */
interface Props {}

export default function NameChangeModal(props: Props) {
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
  const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  const [newName, setNewName] = useState<any>(chatDisplayName);

  const handleNameChange = () => {
    const nameChange = {
      type: MessageType.NAME_CHANGE,
      newName,
    };
    websocketService.send(nameChange);
  };

  const saveEnabled =
    newName !== chatDisplayName && newName !== '' && websocketService?.isConnected();

  return (
    <div>
      Your chat display name is what people see when you send chat messages. Other information can
      go here to mention auth, and stuff.
      <Input
        value={newName}
        onChange={e => setNewName(e.target.value)}
        placeholder="Your chat display name"
        maxLength={10}
        showCount
        defaultValue={chatDisplayName}
      />
      <Button disabled={!saveEnabled} onClick={handleNameChange}>
        Change name
      </Button>
    </div>
  );
}

import { SmileOutlined } from '@ant-design/icons';
import { Button, Input } from 'antd';
import React, { useRef, useState } from 'react';
import { useRecoilValue } from 'recoil';
import ContentEditable from 'react-contenteditable';
import WebsocketService from '../../../services/websocket-service';
import { websocketServiceAtom } from '../../stores/ClientConfigStore';
import { getCaretPosition, convertToText, convertOnPaste } from '../chat';
import { MessageType } from '../../../interfaces/socket-events';
import s from './ChatTextField.module.scss';

interface Props {
  value?: string;
}

export default function ChatTextField(props: Props) {
  const { value: originalValue } = props;
  const [value, setValue] = useState(originalValue);
  const [showEmojis, setShowEmojis] = useState(false);
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);

  const text = useRef(value);

  // large is 40px
  const size = 'small';

  const sendMessage = () => {
    if (!websocketService) {
      console.log('websocketService is not defined');
      return;
    }

    const message = convertToText(value);
    websocketService.send({ type: MessageType.CHAT, body: message });
    setValue('');
  };

  const handleChange = evt => {
    text.current = evt.target.value;
    setValue(evt.target.value);
  };

  const handleKeyDown = event => {
    const key = event && event.key;

    if (key === 'Enter') {
    }
  };

  return (
    <div>
      <Input.Group compact style={{ display: 'flex', width: '100%', position: 'absolute' }}>
        <ContentEditable
          style={{ width: '60%', maxHeight: '50px', padding: '5px' }}
          html={text.current}
          onChange={handleChange}
          onKeyDown={e => {
            handleKeyDown(e);
          }}
        />
        <Button type="default" ghost title="Emoji" onClick={() => setShowEmojis(!showEmojis)}>
          <SmileOutlined style={{ color: 'rgba(0,0,0,.45)' }} />
        </Button>
        <Button size={size} type="primary" onClick={sendMessage}>
          Submit
        </Button>
      </Input.Group>
    </div>
  );
}

ChatTextField.defaultProps = {
  value: '',
};

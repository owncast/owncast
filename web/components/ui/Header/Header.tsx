import s from './Header.module.scss';
import { Layout } from 'antd';
import { ServerStatusStore, serverStatusState } from '../../stores/ServerStatusStore';
import {
  ClientConfigStore,
  clientConfigState,
  chatCurrentlyVisible,
} from '../../stores/ClientConfigStore';
import { ClientConfig } from '../../../interfaces/client-config.model';
import { useRecoilState, useRecoilValue } from 'recoil';
import { useEffect } from 'react';

const { Header } = Layout;

export default function HeaderComponent() {
  const clientConfig = useRecoilValue<ClientConfig>(clientConfigState);
  const [chatOpen, setChatOpen] = useRecoilState(chatCurrentlyVisible);

  const { name } = clientConfig;

  useEffect(() => {
    console.log({ chatOpen });
  }, [chatOpen]);

  return (
    <Header className={`${s.header}`}>
      {name}
      <button onClick={() => setChatOpen(!chatOpen)}>Toggle Chat</button>
    </Header>
  );
}

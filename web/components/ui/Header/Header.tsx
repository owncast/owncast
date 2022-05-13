import { Layout } from 'antd';
import { useRecoilValue } from 'recoil';
import { ChatState } from '../../../interfaces/application-state';
import { OwncastLogo, UserDropdown } from '../../common';
import { chatStateAtom } from '../../stores/ClientConfigStore';
import s from './Header.module.scss';

const { Header } = Layout;

interface Props {
  name: string;
}

export default function HeaderComponent({ name = 'Your stream title' }: Props) {
  const chatState = useRecoilValue<ChatState>(chatStateAtom);

  return (
    <Header className={`${s.header}`}>
      <div className={`${s.logo}`}>
        <OwncastLogo variant="contrast" />
        <span>{name}</span>
      </div>
      <UserDropdown chatState={chatState} />
    </Header>
  );
}

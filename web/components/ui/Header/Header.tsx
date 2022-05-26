import { Layout } from 'antd';
import { OwncastLogo, UserDropdown } from '../../common';
import s from './Header.module.scss';

const { Header } = Layout;

interface Props {
  name: string;
}

export default function HeaderComponent({ name = 'Your stream title' }: Props) {
  return (
    <Header className={`${s.header}`}>
      <div className={`${s.logo}`}>
        <OwncastLogo variant="contrast" />
        <span>{name}</span>
      </div>
      <UserDropdown />
    </Header>
  );
}

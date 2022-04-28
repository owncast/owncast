import { Layout } from 'antd';
import UserDropdown from '../../UserDropdownMenu';
import s from './Header.module.scss';

const { Header } = Layout;

export default function HeaderComponent(props) {
  const { name } = props;

  return (
    <Header className={`${s.header}`}>
      Logo goes here
      {name}
      <UserDropdown />
    </Header>
  );
}

import { Tabs } from 'antd';
import { useRecoilValue } from 'recoil';
import { FC } from 'react';
import IndieAuthModal from '../IndieAuthModal/IndieAuthModal';
import FediAuthModal from '../FediAuthModal/FediAuthModal';

import FediverseIcon from '../../../assets/images/fediverse-black.png';
import IndieAuthIcon from '../../../assets/images/indieauth.png';

import s from './AuthModal.module.scss';
import {
  chatDisplayNameAtom,
  chatAuthenticatedAtom,
  accessTokenAtom,
} from '../../stores/ClientConfigStore';

const { TabPane } = Tabs;

export const AuthModal: FC = () => {
  const chatDisplayName = useRecoilValue<string>(chatDisplayNameAtom);
  const authenticated = useRecoilValue<boolean>(chatAuthenticatedAtom);
  const accessToken = useRecoilValue<string>(accessTokenAtom);
  const federationEnabled = false;

  return (
    <div>
      <Tabs
        defaultActiveKey="1"
        type="card"
        size="small"
        renderTabBar={federationEnabled ? null : () => null}
      >
        <TabPane
          tab={
            <span className={s.tabContent}>
              <img className={s.icon} src={IndieAuthIcon.src} alt="IndieAuth" />
              IndieAuth
            </span>
          }
          key="1"
        >
          <IndieAuthModal
            authenticated={authenticated}
            displayName={chatDisplayName}
            accessToken={accessToken}
          />
        </TabPane>
        <TabPane
          tab={
            <span className={s.tabContent}>
              <img className={s.icon} src={FediverseIcon.src} alt="Fediverse auth" />
              FediAuth
            </span>
          }
          key="2"
        >
          <FediAuthModal />
        </TabPane>
      </Tabs>
    </div>
  );
};
export default AuthModal;

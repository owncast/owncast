/* eslint-disable react/no-unescaped-entities */
import { FC } from 'react';
import { Space } from 'antd';
import styles from './FollowModal.module.scss';
import { FollowForm } from './FollowForm';
import { Translation } from '../../ui/Translation/Translation';
import { Localization } from '../../../types/localization';

export type FollowModalProps = {
  handleClose: () => void;
  account: string;
  name: string;
};

export const FollowModal: FC<FollowModalProps> = ({ handleClose, account, name }) => (
  <Space direction="vertical" id="follow-modal">
    <div className={styles.header}>
      <Translation
        translationKey={Localization.Frontend.FollowModal.description}
        defaultText="By following this stream you'll get notified on the Fediverse when it goes live. Now is a great time to"
      />
      <a href="https://owncast.online/join-fediverse" target="_blank" rel="noreferrer">
        &nbsp;
        <Translation
          translationKey={Localization.Frontend.FollowModal.learnFediverse}
          defaultText="learn about the Fediverse"
        />
        &nbsp;
      </a>
      <Translation
        translationKey={Localization.Frontend.FollowModal.newToYou}
        defaultText="if it's new to you."
      />
    </div>
    <div className={styles.account}>
      <img src="/logo" alt="logo" className={styles.logo} />
      <div className={styles.username}>
        <div className={styles.name}>{name}</div>
        <div>{account}</div>
      </div>
    </div>

    <FollowForm handleClose={handleClose} />
  </Space>
);

/* eslint-disable react/no-danger */

import { FC, useEffect, useState } from 'react';
import classNames from 'classnames';
import Head from 'next/head';
import { Button, Spin, Alert, Typography } from 'antd';
import styles from './OfflineEmbed.module.scss';
import { Modal } from '../Modal/Modal';
import { FollowForm } from '../../modals/FollowModal/FollowForm';

const { Title } = Typography;

export type OfflineEmbedProps = {
  streamName: string;
  subtitle?: string;
  image: string;
  supportsFollows: boolean;
};

enum EmbedMode {
  CannotFollow = 1,
  CanFollow,
  // FollowPrompt,
  // InProgress,
}

export const OfflineEmbed: FC<OfflineEmbedProps> = ({
  streamName,
  subtitle,
  image,
  supportsFollows,
}) => {
  const [currentMode, setCurrentMode] = useState(EmbedMode.CanFollow);
  // const [loading, setLoading] = useState(false);
  // const [errorMessage, setErrorMessage] = useState(null);
	const [showFollowModal, setShowFollowModal] = useState(false);

  useEffect(() => {
    if (!supportsFollows) {
      setCurrentMode(EmbedMode.CannotFollow);
    } else if (currentMode === EmbedMode.CannotFollow) {
      setCurrentMode(EmbedMode.CanFollow);
    }
  }, [supportsFollows]);

  const followButtonPressed = async () => {
    // setCurrentMode(EmbedMode.FollowPrompt);
		setShowFollowModal(true);
  };

  // const handleErrorClose = () => {
  //   setErrorMessage('');
  //   setCurrentMode(EmbedMode.FollowPrompt);
  // };

  return (
    <div>
      <Head>
        <title>{streamName}</title>
      </Head>
      <div className={classNames(styles.offlineContainer)}>
				<div className={classNames(styles.content, {
					[styles.followable]: supportsFollows,
				})}>
					<Title level={1} className={styles.headerContainer}>
						<div className={styles.pageLogo} style={{ backgroundImage: `url(${image})` }} />
						<div className={styles.streamName}>{streamName}</div>
					</Title>

					<div className={styles.messageContainer}>
						<Title level={2} className={styles.offlineTitle}>
							This stream is not currently live.
						</Title>
						<div className={styles.message} dangerouslySetInnerHTML={{ __html: subtitle }} />
					</div>

					{/* {errorMessage && (
						<Alert
							message="Follow Error"
							description={errorMessage}
							type="error"
							showIcon
							closable
							onClose={handleErrorClose}
						/>
					)} */}

					{currentMode === EmbedMode.CanFollow && (
						<>
							<Button className={styles.submitButton} type="primary" onClick={followButtonPressed}>
								Follow Server
							</Button>
							<Modal
								title={`Follow ${name}`}
								open={showFollowModal}
								handleCancel={() => setShowFollowModal(false)}
							>
								<FollowForm  />
							</Modal>
						</>
					)}

					{/* {currentMode === EmbedMode.InProgress && (
						<Title level={4} className={styles.heading}>
							Follow the instructions on your Fediverse server to complete the follow.
						</Title>
					)} */}
				</div>
      </div>
    </div>
  );
};

import { Modal, Button } from 'antd';

import { FC } from 'react';
import dynamic from 'next/dynamic';
import { USER_SET_MODERATOR, fetchData } from '../../utils/apis';
import { User } from '../../types/chat';

// Lazy loaded components

const ExclamationCircleFilled = dynamic(() => import('@ant-design/icons/ExclamationCircleFilled'), {
  ssr: false,
});

const QuestionCircleFilled = dynamic(() => import('@ant-design/icons/QuestionCircleFilled'), {
  ssr: false,
});

const StopTwoTone = dynamic(() => import('@ant-design/icons/StopTwoTone'), {
  ssr: false,
});

const SafetyCertificateTwoTone = dynamic(
  () => import('@ant-design/icons/SafetyCertificateTwoTone'),
  {
    ssr: false,
  },
);

export type ModeratorUserButtonProps = {
  user: User;
  onClick?: () => void;
};

export const ModeratorUserButton: FC<ModeratorUserButtonProps> = ({ user, onClick }) => {
  async function buttonClicked({ id }, setAsModerator: Boolean): Promise<Boolean> {
    const data = {
      userId: id,
      isModerator: setAsModerator,
    };
    try {
      const result = await fetchData(USER_SET_MODERATOR, {
        data,
        method: 'POST',
        auth: true,
      });
      return result.success;
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error(e);
    }
    return false;
  }

  const isModerator = user.scopes?.includes('MODERATOR');
  const actionString = isModerator ? 'remove moderator' : 'add moderator';
  const icon = isModerator ? (
    <ExclamationCircleFilled style={{ color: 'var(--ant-error)' }} />
  ) : (
    <QuestionCircleFilled style={{ color: 'var(--ant-warning)' }} />
  );

  const content = (
    <>
      Are you sure you want to {actionString} <strong>{user.displayName}</strong>?
    </>
  );

  const confirmBlockAction = () => {
    Modal.confirm({
      title: `Confirm ${actionString}`,
      content,
      onCancel: () => {},
      onOk: () =>
        new Promise((resolve, reject) => {
          const result = buttonClicked(user, !isModerator);
          if (result) {
            // wait a bit before closing so the user/client tables repopulate
            // GW: TODO: put users/clients data in global app context instead, then call a function here to update that state. (current in another branch)
            setTimeout(() => {
              resolve(result);
              onClick?.();
            }, 3000);
          } else {
            reject();
          }
        }),
      okType: 'danger',
      okText: isModerator ? 'Yup!' : null,
      icon,
    });
  };

  return (
    <Button
      type="primary"
      onClick={confirmBlockAction}
      size="small"
      icon={
        isModerator ? (
          <StopTwoTone twoToneColor="#ff4d4f" />
        ) : (
          <SafetyCertificateTwoTone twoToneColor="#22bb44" />
        )
      }
      className="block-user-button"
    >
      {actionString}
    </Button>
  );
};

ModeratorUserButton.defaultProps = {
  onClick: null,
};

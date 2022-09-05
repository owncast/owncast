import { Modal, Button } from 'antd';
import { ExclamationCircleFilled, QuestionCircleFilled, StopTwoTone } from '@ant-design/icons';
import { USER_ENABLED_TOGGLE, fetchData } from '../utils/apis';
import { User } from '../types/chat';

interface BanUserButtonProps {
  user: User;
  isEnabled: Boolean; // = this user's current status
  label?: string;
  onClick?: () => void;
}
export default function BanUserButton({ user, isEnabled, label, onClick }: BanUserButtonProps) {
  async function buttonClicked({ id }): Promise<Boolean> {
    const data = {
      userId: id,
      enabled: !isEnabled, // set user to this value
    };
    try {
      const result = await fetchData(USER_ENABLED_TOGGLE, {
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

  const actionString = isEnabled ? 'ban' : 'unban';
  const icon = isEnabled ? (
    <ExclamationCircleFilled style={{ color: 'var(--ant-error)' }} />
  ) : (
    <QuestionCircleFilled style={{ color: 'var(--ant-warning)' }} />
  );

  const content = (
    <>
      Are you sure you want to {actionString} <strong>{user.displayName}</strong>
      {isEnabled ? ' and remove their messages?' : '?'}
    </>
  );

  const confirmBlockAction = () => {
    Modal.confirm({
      title: `Confirm ${actionString}`,
      content,
      onCancel: () => {},
      onOk: () =>
        new Promise((resolve, reject) => {
          const result = buttonClicked(user);
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
      okText: isEnabled ? 'Absolutely' : null,
      icon,
    });
  };

  return (
    <Button
      onClick={confirmBlockAction}
      size="small"
      icon={isEnabled ? <StopTwoTone twoToneColor="#ff4d4f" /> : null}
      className="block-user-button"
    >
      {label || actionString}
    </Button>
  );
}
BanUserButton.defaultProps = {
  label: '',
  onClick: null,
};

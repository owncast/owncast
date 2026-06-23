import { Button, Dropdown, Menu, Modal } from 'antd';
import { FC, ReactNode } from 'react';
import dynamic from 'next/dynamic';
import { USER_ENABLED_TOGGLE, USER_SET_MODERATOR, DELETE_USER, fetchData } from '../../utils/apis';
import { User } from '../../types/chat';

// Lazy loaded components

const DownOutlined = dynamic(() => import('@ant-design/icons/DownOutlined'), {
  ssr: false,
});

const ExclamationCircleFilled = dynamic(() => import('@ant-design/icons/ExclamationCircleFilled'), {
  ssr: false,
});

export type UserActionsDropdownProps = {
  user: User;
  onChanged?: () => void;
};

// UserActionsDropdown collapses the per-user moderation actions (moderator
// toggle, ban/unban, delete) into a single dropdown. Every action confirms
// before it runs, and calls onChanged so the caller can refresh its data.
export const UserActionsDropdown: FC<UserActionsDropdownProps> = ({ user, onChanged }) => {
  const isEnabled = !user.disabledAt;
  const isModerator = user.scopes?.includes('MODERATOR');

  async function runAction(url: string, data: Record<string, unknown>) {
    try {
      const result = await fetchData(url, { data, method: 'POST', auth: true });
      if (result?.success) {
        onChanged?.();
      }
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error(e);
    }
  }

  function showConfirm(options: {
    title: string;
    content: ReactNode;
    okText: string;
    danger?: boolean;
    onOk: () => Promise<void>;
  }) {
    Modal.confirm({
      title: options.title,
      icon: (
        <ExclamationCircleFilled
          style={{ color: options.danger ? 'var(--ant-error)' : 'var(--ant-warning)' }}
        />
      ),
      content: options.content,
      okText: options.okText,
      okType: options.danger ? 'danger' : 'primary',
      onOk: options.onOk,
    });
  }

  const onModerator = () =>
    showConfirm({
      title: isModerator ? 'Remove moderator' : 'Add moderator',
      content: isModerator ? (
        <>
          Remove moderator access from <strong>{user.displayName}</strong>?
        </>
      ) : (
        <>
          Make <strong>{user.displayName}</strong> a moderator?
        </>
      ),
      okText: isModerator ? 'Remove' : 'Add',
      onOk: () => runAction(USER_SET_MODERATOR, { userId: user.id, isModerator: !isModerator }),
    });

  const onBan = () =>
    showConfirm({
      title: isEnabled ? 'Ban user' : 'Unban user',
      danger: true,
      content: isEnabled ? (
        <>
          Ban <strong>{user.displayName}</strong> and remove their messages?
        </>
      ) : (
        <>
          Unban <strong>{user.displayName}</strong>?
        </>
      ),
      okText: isEnabled ? 'Ban' : 'Unban',
      onOk: () => runAction(USER_ENABLED_TOGGLE, { userId: user.id, enabled: !isEnabled }),
    });

  const onDelete = () =>
    showConfirm({
      title: 'Delete user',
      danger: true,
      content: (
        <>
          Permanently delete <strong>{user.displayName}</strong> and all of their data? This cannot
          be undone.
        </>
      ),
      okText: 'Delete',
      onOk: () => runAction(DELETE_USER, { userId: user.id }),
    });

  const menu = (
    <Menu
      items={[
        {
          key: 'moderator',
          label: isModerator ? 'Remove moderator' : 'Add moderator',
          onClick: onModerator,
        },
        {
          key: 'ban',
          label: isEnabled ? 'Ban user' : 'Unban user',
          onClick: onBan,
        },
        { type: 'divider' },
        {
          key: 'delete',
          label: 'Delete user',
          danger: true,
          onClick: onDelete,
        },
      ]}
    />
  );

  return (
    <Dropdown overlay={menu} trigger={['click']}>
      <Button size="small">
        Actions <DownOutlined />
      </Button>
    </Dropdown>
  );
};

UserActionsDropdown.defaultProps = {
  onChanged: null,
};

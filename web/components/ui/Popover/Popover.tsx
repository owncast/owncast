import React, { FC, ReactNode } from 'react';
import styles from './Popover.module.scss';

export type PopoverProps = {
  open: boolean;
  title: ReactNode;
  content: ReactNode;
  children?: ReactNode;
};

/// Replace Popover from antd with a custom popover implementation.
///
/// The popover is positioned relative to its parent (unlike antd's popover,
/// which uses absolute positioning on the page).
//
// TODO the following properties of antd's popover have not been implemented
//  placement ("topRight" is used as default)
//  defaultOpen
//  destroyTooltipOnHide
//  overlayInnerStyle
//  color (it uses var(--theme-color-components-primary-button-background))

export const Popover: FC<PopoverProps> = ({ open, title, content, children }) => (
  <div style={{ width: 'max-content', height: 'max-content' }}>
    {open && (
      <div className={styles.anchor}>
        <div className={styles.popover}>
          <div className={styles.title}>{title}</div>
          <hr style={{ color: 'var(--color-owncast-palette-4)' }} />
          <div className={styles.content}>{content}</div>
        </div>
      </div>
    )}
    {children}
  </div>
);

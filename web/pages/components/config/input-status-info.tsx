import React from 'react';
import classNames from 'classnames';

import { StatusState } from '../../../utils/input-statuses';

interface InputStatusInfoProps {
  status: StatusState;
}
export default function InputStatusInfo({ status }: InputStatusInfoProps) {
  const { type, icon, message } = status || {};
  const classes = classNames({
    'status-container': true,
    [`status-${type}`]: type,
    empty: !message,
  });
  return (
    <div className={classes}>
      {icon ? <span className="status-icon">{icon}</span> : null}
      {message ? <span className="status-message">{message}</span> : null}
    </div>
  );
}

import React, { FC } from 'react';
import classNames from 'classnames';

import { StatusState } from '../../utils/input-statuses';

export type FormStatusIndicatorProps = {
  status: StatusState;
};

export const FormStatusIndicator: FC<FormStatusIndicatorProps> = ({ status }) => {
  const { type, icon, message } = status || {};
  const classes = classNames({
    'status-container': true,
    [`status-${type}`]: type,
    empty: !message,
  });
  return (
    <span className={classes}>
      {icon ? <span className="status-icon">{icon}</span> : null}
      {message ? <span className="status-message">{message}</span> : null}
    </span>
  );
};
export default FormStatusIndicator;

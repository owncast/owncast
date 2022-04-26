import {
  CheckCircleFilled,
  ExclamationCircleFilled,
  LoadingOutlined,
  WarningOutlined,
} from '@ant-design/icons';

export const STATUS_RESET_TIMEOUT = 3000;

export const STATUS_ERROR = 'error';
export const STATUS_INVALID = 'invalid';
export const STATUS_PROCESSING = 'proessing';
export const STATUS_SUCCESS = 'success';
export const STATUS_WARNING = 'warning';

export type InputStatusTypes = 'error' | 'invalid' | 'proessing' | 'success' | 'warning';

export interface StatusState {
  type: InputStatusTypes;
  icon: any; // Element type of sorts?
  message: string;
}
interface InputStates {
  [key: string]: StatusState;
}

export const INPUT_STATES: InputStates = {
  [STATUS_SUCCESS]: {
    type: STATUS_SUCCESS,
    icon: <CheckCircleFilled style={{ color: 'green' }} />,
    message: 'Success!',
  },
  [STATUS_ERROR]: {
    type: STATUS_ERROR,
    icon: <ExclamationCircleFilled style={{ color: 'red' }} />,
    message: 'An error occurred.',
  },
  [STATUS_INVALID]: {
    type: STATUS_INVALID,
    icon: <ExclamationCircleFilled style={{ color: 'red' }} />,
    message: 'An error occurred.',
  },
  [STATUS_PROCESSING]: {
    type: STATUS_PROCESSING,
    icon: <LoadingOutlined />,
    message: '',
  },
  [STATUS_WARNING]: {
    type: STATUS_WARNING,
    icon: <WarningOutlined style={{ color: '#fc0' }} />,
    message: '',
  },
};

// Don't like any of the default messages in INPUT_STATES? Create a state with custom message by providing an icon style with your message.
export function createInputStatus(type: InputStatusTypes, message?: string): StatusState {
  if (!type || !INPUT_STATES[type]) {
    return null;
  }
  if (!message) {
    return INPUT_STATES[type];
  }
  return {
    type,
    icon: INPUT_STATES[type].icon,
    message,
  };
}

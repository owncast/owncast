import { CheckCircleFilled, ExclamationCircleFilled, LoadingOutlined, WarningOutlined } from '@ant-design/icons';

export const STATUS_RESET_TIMEOUT = 3000;

export const STATUS_ERROR = 'error';
export const STATUS_INVALID = 'invalid';
export const STATUS_PROCESSING = 'proessing';
export const STATUS_SUCCESS = 'success';
export const STATUS_WARNING = 'warning';

export type InputStatusTypes = 
  typeof STATUS_ERROR | 
  typeof STATUS_INVALID | 
  typeof STATUS_PROCESSING | 
  typeof STATUS_SUCCESS | 
  typeof STATUS_WARNING;

export type StatusState = {
  icon: any; // Element type of sorts?
  message: string;
};

export const INPUT_STATES = {
  [STATUS_SUCCESS]: {
    icon: <CheckCircleFilled style={{ color: 'green' }} />,
    message: 'Success!',
  },
  [STATUS_ERROR]: {
    icon: <ExclamationCircleFilled style={{ color: 'red' }} />,
    message: 'An error occurred.',  
  },
  [STATUS_INVALID]: {
    icon: <ExclamationCircleFilled style={{ color: 'red' }} />,
    message: 'An error occurred.',  
  },
  [STATUS_PROCESSING]: {
    icon: <LoadingOutlined />,
    message: '',
  },
  [STATUS_WARNING]: {
    icon: <WarningOutlined style={{ color: '#fc0' }} />,
    message: '',  
  },
};

// Don't like any of the default messages in INPUT_STATES? Create a state with custom message by providing an icon style with your message.
export function createInputStatus(type: InputStatusTypes, message?: string): StatusState {
  if (!type || !INPUT_STATES[type]) {
    return null;
  }
  if (message === null) {
    return INPUT_STATES[type];
  }
  return {
    icon: INPUT_STATES[type].icon,
    message,
  };
}

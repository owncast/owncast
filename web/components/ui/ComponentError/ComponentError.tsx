import { Alert, Button } from 'antd';
import { FC } from 'react';

export type ComponentErrorProps = {
  message?: string;
  componentName: string;
  details?: string;
  retryFunction?: () => void;
};

const openBugReport = () => {
  window.open(
    'https://github.com/owncast/owncast/issues/new?assignees=&labels=&template=bug-report-feature-request.yml',
    '_blank',
  );
};

const ErrorContent = ({
  message,
  componentName,
  details,
  canRetry,
}: {
  message: string;
  componentName: string;
  details: string;
  canRetry: boolean;
}) => (
  <div>
    <p>
      There was an unexpected error. It would be appreciated if you would report this so it can be
      fixed in the future.
    </p>
    {!!canRetry && (
      <p>You may optionally retry, however functionality might not work as expected.</p>
    )}
    <code>
      <div>{message && `Error: ${message}`}</div>
      <div>Component: {componentName}</div>
      <div>{details && details}</div>
    </code>
  </div>
);

export const ComponentError: FC<ComponentErrorProps> = ({
  message,
  componentName,
  details,
  retryFunction,
}: ComponentErrorProps) => (
  <Alert
    message="Error"
    showIcon
    description=<ErrorContent
      message={message}
      details={details}
      componentName={componentName}
      canRetry={!!retryFunction}
    />
    type="error"
    action={
      <>
        {retryFunction && (
          <Button ghost size="small" onClick={retryFunction}>
            Retry
          </Button>
        )}
        <Button ghost size="small" danger onClick={openBugReport}>
          Report Error
        </Button>
      </>
    }
  />
);

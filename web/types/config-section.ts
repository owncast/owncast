// TS types for elements on the Config pages

export interface TextFieldProps {
  handleResetValue: ({ fieldName }) => void;
  fieldName: string;
  initialValues: any;
  type: string;
}

export interface UpdateArgs {
  fieldName: string;
  value: string;
  path?: string;
}

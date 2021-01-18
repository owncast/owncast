// TS types for elements on the Config pages

export interface TextFieldProps {
  handleResetValue?: (fieldName: string) => void;
  fieldName: string;
  initialValues?: any;
  type?: string;
  configPath?: string;
  required?: boolean;
  disabled?: boolean;
  onSubmit?: () => void;
  onBlur?: () => void;
  onChange?: () => void;
}

export interface ToggleSwitchProps {
  fieldName: string;
  initialValues?: any;
  configPath?: string;
  disabled?: boolean;
}

export type FieldUpdaterFunc = (args: UpdateArgs) => void;

export interface UpdateArgs {
  fieldName: string;
  value: any;
  path?: string;
}

export interface ApiPostArgs {
  apiPath: string,
  data: object,
  onSuccess?: (arg: any) => void,
  onError?: (arg: any) => void,
}

export interface ConfigDirectoryFields {
  enabled: boolean;
  instanceUrl: string,
}

export interface ConfigInstanceDetailsFields {
  extraPageContent: string;
  logo: string;
  name: string;
  nsfw: boolean;
  streamTitle: string;
  summary: string;
  tags: string[];
  title: string;
}


export type PRESET = 'fast' | 'faster' | 'veryfast' | 'superfast' | 'ultrafast';

export interface VideoVariant {
  key?: number; // unique identifier generated on client side just for ant table rendering
  encoderPreset: PRESET,
  framerate: number;

  audioPassthrough: boolean;
  audioBitrate: number;
  videoPassthrough: boolean;
  videoBitrate: number;
}
export interface VideoSettingsFields {
  latencyLevel: number;
  videoQualityVariants: VideoVariant[],
}


export interface ConfigDetails {
  ffmpegPath: string;
  instanceDetails: ConfigInstanceDetailsFields;
  rtmpServerPort: string;
  s3: any; // tbd
  streamKey: string;
  webServerPort: string;
  yp: ConfigDirectoryFields;
  videoSettings: VideoSettingsFields;
}

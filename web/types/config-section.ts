// TS types for elements on the Config pages

export interface TextFieldProps {
  apiPath: string;
  fieldName: string;

  configPath?: string;
  disabled?: boolean;
  initialValue?: string;
  label?: string;
  maxLength?: number;
  placeholder?: string;
  required?: boolean;
  tip?: string;
  type?: string;
  value?: string;
  onSubmit?: () => void;
  onBlur?: () => void;
  onChange?: () => void;
}

export interface ToggleSwitchProps {
  apiPath: string;
  fieldName: string;

  checked?: boolean;
  configPath?: string;
  disabled?: boolean;
  label?: string;
  tip?: string;
}

// for dropdown
export interface SocialHandleDropdownItem {
  icon: string;
  platform: string;
  key: string;
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
  socialHandles: SocialHandle[],
  streamTitle: string;
  summary: string;
  tags: string[];
  title: string;
}


export type PRESET = 'fast' | 'faster' | 'veryfast' | 'superfast' | 'ultrafast';

// from data
export interface SocialHandle {
  platform: string;
  url: string;
}

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

export as namespace EmojiButton;

export = EmojiButton;

declare namespace EmojiButton {
  export class EmojiButton {
    constructor(options?: EmojiButton.Options);
    on(event: Event, callback: (selection: EmojiSelection) => void): void;
    off(event: Event, callback: (selection: EmojiSelection) => void): void;
    hidePicker(): void;
    destroyPicker(): void;
    showPicker(referenceEl: HTMLElement): void;
    togglePicker(referenceEl: HTMLElement): void;
    isPickerVisible(): boolean;
    setTheme(theme: EmojiTheme): void;
  }

  export interface Options {
    position?: Placement | FixedPosition;
    autoHide?: boolean;
    autoFocusSearch?: boolean;
    showAnimation?: boolean;
    showPreview?: boolean;
    showSearch?: boolean;
    showRecents?: boolean;
    showVariants?: boolean;
    showCategoryButtons?: boolean;
    recentsCount?: number;
    emojiVersion?: EmojiVersion;
    i18n?: I18NStrings;
    zIndex?: number;
    theme?: EmojiTheme;
    categories?: Category[];
    style?: EmojiStyle;
    emojisPerRow?: number;
    rows?: number;
    emojiSize?: string;
    initialCategory?: Category | 'recents';
    custom?: CustomEmoji[];
    plugins?: Plugin[];
    icons?: Icons;
    rootElement?: HTMLElement;
  }

  export interface FixedPosition {
    top?: string;
    bottom?: string;
    left?: string;
    right?: string;
  }

  export interface Plugin {
    render(picker: EmojiButton): HTMLElement;
    destroy?(): void;
  }

  export interface EmojiSelection {
    name: string;
    custom?: boolean;
    emoji?: string;
    url?: string;
  }

  export interface CustomEmoji {
    name: string;
    emoji: string;
  }

  export type EmojiStyle = 'native' | 'twemoji';

  export type EmojiTheme = 'dark' | 'light' | 'auto';

  export type Event = 'emoji' | 'hidden';

  export type Placement =
    | 'auto'
    | 'auto-start'
    | 'auto-end'
    | 'top'
    | 'top-start'
    | 'top-end'
    | 'bottom'
    | 'bottom-start'
    | 'bottom-end'
    | 'right'
    | 'right-start'
    | 'right-end'
    | 'left'
    | 'left-start'
    | 'left-end';

  export type EmojiVersion =
    | '1.0'
    | '2.0'
    | '3.0'
    | '4.0'
    | '5.0'
    | '11.0'
    | '12.0'
    | '12.1';

  export type Category =
    | 'smileys'
    | 'people'
    | 'animals'
    | 'food'
    | 'activities'
    | 'travel'
    | 'objects'
    | 'symbols'
    | 'flags';

  export type I18NCategory =
    | 'recents'
    | 'smileys'
    | 'people'
    | 'animals'
    | 'food'
    | 'activities'
    | 'travel'
    | 'objects'
    | 'symbols'
    | 'flags'
    | 'custom';

  export interface I18NStrings {
    search: string;
    categories: {
      [key in I18NCategory]: string;
    };
    notFound: string;
  }

  export interface Icons {
    search?: string;
    clearSearch?: string;
    categories?: {
      [key in I18NCategory]?: string;
    };
    notFound?: string;
  }
}

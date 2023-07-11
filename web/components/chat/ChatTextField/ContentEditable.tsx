import * as React from 'react';

export interface ContentEditableProps extends React.HTMLAttributes<HTMLElement> {
  onRootRef?: Function;
  onContentChange?: Function;
  tagName?: string;
  html: string;
  disabled: boolean;
}

export default class ContentEditable extends React.Component<ContentEditableProps> {
  private root: HTMLElement;

  private mutationObserver: MutationObserver;

  private innerHTMLBuffer: string;

  public componentDidMount() {
    this.mutationObserver = new MutationObserver(this.onContentChange);
    this.mutationObserver.observe(this.root, {
      childList: true,
      subtree: true,
      characterData: true,
    });
  }

  private onContentChange = (mutations: MutationRecord[]) => {
    mutations.forEach(() => {
      const { innerHTML } = this.root;

      if (this.innerHTMLBuffer === undefined || this.innerHTMLBuffer !== innerHTML) {
        this.innerHTMLBuffer = innerHTML;

        if (this.props.onContentChange) {
          this.props.onContentChange({
            target: {
              value: innerHTML,
            },
          });
        }
      }
    });
  };

  private onRootRef = (elt: HTMLElement) => {
    this.root = elt;
    if (this.props.onRootRef) {
      this.props.onRootRef(this.root);
    }
  };

  public render() {
    const { tagName, html, ...newProps } = this.props;

    delete newProps.onRootRef;
    delete newProps.onContentChange;

    return React.createElement(tagName || 'div', {
      ...newProps,
      ref: this.onRootRef,
      contentEditable: !this.props.disabled,
      dangerouslySetInnerHTML: { __html: html },
    });
  }
}

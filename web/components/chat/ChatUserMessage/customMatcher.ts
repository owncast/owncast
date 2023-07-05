/* eslint-disable class-methods-use-this */
import { ChildrenNode, Matcher, MatchResponse, Node } from 'interweave';
import React from 'react';

export interface CustomProps {
  children: React.ReactNode;
  key: string;
}

interface options {
  highlightString: string;
}

export class ChatMessageHighlightMatcher extends Matcher {
  match(str: string): MatchResponse<{}> | null {
    const { highlightString } = this.options as options;

    if (!highlightString) {
      return null;
    }

    const highlightRegex = new RegExp(highlightString.replace(/\s/g, '\\s'), 'ui');

    const result = str.match(highlightRegex);

    if (!result) {
      return null;
    }

    return {
      index: result.index!,
      length: result[0].length,
      match: result[0],
      valid: true,
    };
  }

  replaceWith(children: ChildrenNode, props: CustomProps): Node {
    const { key } = props;
    return React.createElement('mark', { key }, children);
  }

  asTag(): string {
    return 'mark';
  }
}

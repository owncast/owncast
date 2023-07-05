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

    const escapedString = highlightString
      .replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
      .replace(/\s/g, '\\s');

    const normalizedString = escapedString.normalize('NFD').replace(/[\u0300-\u036f]/g, '');

    let highlightRegex = escapedString;
    if (escapedString !== normalizedString) {
      highlightRegex = `(?:${escapedString})|(?:${normalizedString})`;
    }

    const result = str.match(new RegExp(highlightRegex, 'ui'));

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

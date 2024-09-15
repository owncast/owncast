/* eslint-disable class-methods-use-this */
import { ChildrenNode, Matcher, MatchResponse, Node } from 'interweave';
import rewritePattern from 'regexpu-core';
import React from 'react';

export interface ChatMessageEmojiProps {
  key: string;
}

interface options {
  className: string;
}

const emojiPattern = '\\p{RGI_Emoji}';

const regexSupportsUnicodeSets = (() => {
  // Using a variable for regexpFlags to avoid eslint error about the flag
  // being invalid. It's not invalid, it's just new.
  const regexpFlags = 'v';

  // A bit more working around eslint - just calling new RegExp throws an
  // error about not saving the value / just using side effects.
  let regexp = null;
  try {
    regexp = new RegExp(emojiPattern, regexpFlags);
  } catch {
    return false;
  }

  // We have to use the variable somehow. Since we didn't catch an error
  // this line always returns true.
  return regexp !== null;
})();

const emojiRegex = (() => {
  const rewriteFlags = {
    unicodeSetsFlag: regexSupportsUnicodeSets ? 'parse' : 'transform',
  };
  const regexFlag = regexSupportsUnicodeSets ? 'v' : 'u';

  const regexPattern = rewritePattern(emojiPattern, 'v', rewriteFlags);

  return new RegExp(regexPattern, regexFlag);
})();

export class ChatMessageEmojiMatcher extends Matcher {
  match(str: string): MatchResponse<{}> | null {
    const result = str.match(emojiRegex);

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

  replaceWith(children: ChildrenNode, props: ChatMessageEmojiProps): Node {
    const { key } = props;
    const { className } = this.options as options;
    return React.createElement('span', { key, className }, children);
  }

  asTag(): string {
    return 'span';
  }
}

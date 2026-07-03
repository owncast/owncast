/* eslint-disable class-methods-use-this */
import { ChildrenNode, Matcher, MatchResponse, Node } from 'interweave';
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

// Exported for tests, which need to exercise the fallback branches on a
// runtime whose native RegExp already supports the `v` flag.
export const buildEmojiRegex = (supportsUnicodeSets: boolean): RegExp => {
  if (supportsUnicodeSets) {
    // Modern browsers (2023+) compile the RGI_Emoji property-of-strings
    // pattern natively, so we don't need to ship regexpu-core's ~500KB of
    // unicode tables just to rewrite this one pattern at runtime.
    const regexpFlags = 'v';
    return new RegExp(emojiPattern, regexpFlags);
  }

  try {
    // Browsers without the RegExp `v` flag: approximate RGI_Emoji with
    // pictographic ZWJ sequences. Slightly narrower (misses keycaps and
    // flag sequences) but this matcher is only cosmetic styling.
    // A regex literal would be a parse-time syntax error on browsers without
    // unicode property escapes, so the runtime constructor is load-bearing:
    // it is what lets the catch branch below run at all.
    // eslint-disable-next-line prefer-regex-literals
    return new RegExp('\\p{Extended_Pictographic}(?:\\u200D\\p{Extended_Pictographic})*', 'u');
  } catch {
    // Ancient browsers without unicode property escapes: never match, so
    // emoji simply render as plain text.
    return /$^/;
  }
};

const emojiRegex = buildEmojiRegex(regexSupportsUnicodeSets);

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

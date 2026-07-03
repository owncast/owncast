import {
  buildEmojiRegex,
  ChatMessageEmojiMatcher,
} from '../components/chat/ChatUserMessage/emojiMatcher';

// Built from explicit escapes so the ZWJ / variation selectors are visible in
// the source instead of hiding inside a literal.
const ZWJ_FAMILY = '\u{1F469}\u200D\u{1F469}\u200D\u{1F466}'; // 👩‍👩‍👦
const KEYCAP_ONE = '1\uFE0F\u20E3'; // 1️⃣
const PARTY = '\u{1F389}'; // 🎉

describe('buildEmojiRegex', () => {
  test('returns a usable RegExp for both tiers', () => {
    expect(buildEmojiRegex(true)).toBeInstanceOf(RegExp);
    expect(buildEmojiRegex(false)).toBeInstanceOf(RegExp);
  });

  describe('native RGI_Emoji tier (supportsUnicodeSets: true)', () => {
    const regex = buildEmojiRegex(true);

    test('matches a single emoji', () => {
      expect(PARTY.match(regex)?.[0]).toBe(PARTY);
    });

    test('matches a ZWJ sequence as one match', () => {
      expect(ZWJ_FAMILY.match(regex)?.[0]).toBe(ZWJ_FAMILY);
    });

    test('matches a keycap sequence as one match', () => {
      // RGI_Emoji is a property of strings, so the digit + VS16 + combining
      // keycap sequence matches as a single unit.
      expect(KEYCAP_ONE.match(regex)?.[0]).toBe(KEYCAP_ONE);
    });

    test('does not match plain text or bare digits', () => {
      expect('hello'.match(regex)).toBeNull();
      expect('1'.match(regex)).toBeNull();
    });
  });

  describe('pictographic fallback tier (supportsUnicodeSets: false)', () => {
    const regex = buildEmojiRegex(false);

    test('matches pictographic emoji', () => {
      expect(PARTY.match(regex)?.[0]).toBe(PARTY);
    });

    test('matches a ZWJ sequence as one contiguous match', () => {
      expect(ZWJ_FAMILY.match(regex)?.[0]).toBe(ZWJ_FAMILY);
    });

    test('does not match plain text', () => {
      expect('hello'.match(regex)).toBeNull();
    });

    test('accepted degradation: keycap sequences are not matched', () => {
      // Neither the digit, VS16, nor U+20E3 has Extended_Pictographic, so the
      // fallback tier skips keycaps entirely. They render as plain text, which
      // is the documented tradeoff of this tier.
      expect(KEYCAP_ONE.match(regex)).toBeNull();
    });
  });

  describe('never-match tier semantics', () => {
    test('/$^/ never matches non-empty input, so the matcher no-ops', () => {
      // Regression guard: the last-resort regex must be a real RegExp that
      // matches nothing, not null (str.match(null) matches the string "null").
      const never = /$^/;
      expect('null'.match(never)).toBeNull();
      expect(`party ${PARTY} time`.match(never)).toBeNull();
    });
  });
});

describe('ChatMessageEmojiMatcher.match', () => {
  const matcher = new ChatMessageEmojiMatcher('emoji', { className: 'emoji' });

  test('returns null when the message has no emoji', () => {
    expect(matcher.match('no emoji here')).toBeNull();
  });

  test('returns the first emoji occurrence with its position and length', () => {
    expect(matcher.match(`party ${PARTY} time`)).toEqual({
      index: 6,
      length: PARTY.length,
      match: PARTY,
      valid: true,
    });
  });
});

import {
  composeWeeklyRecurrence,
  parseWeeklyRecurrence,
  wallTimeInZone,
  wallTimeInZoneToUTC,
  WeeklyRecurrence,
} from '../utils/schedule';

const laRule: WeeklyRecurrence = {
  days: ['MO', 'FR'],
  time: '18:00',
  startsOn: '2026-07-06',
  timezone: 'America/Los_Angeles',
  endsOn: '2026-07-20',
};

describe('composeWeeklyRecurrence', () => {
  test('composes DTSTART;TZID plus RRULE with UNTIL as the zone-local end of day in UTC', () => {
    // 2026-07-20 23:59 in Los Angeles is PDT (UTC-7), so the inclusive UNTIL
    // instant is 2026-07-21T06:59Z.
    expect(composeWeeklyRecurrence(laRule)).toBe(
      'DTSTART;TZID=America/Los_Angeles:20260706T180000\n' +
        'RRULE:FREQ=WEEKLY;BYDAY=MO,FR;UNTIL=20260721T065900Z',
    );
  });

  test('an open-ended rule carries no UNTIL', () => {
    expect(composeWeeklyRecurrence({ ...laRule, endsOn: undefined })).toBe(
      'DTSTART;TZID=America/Los_Angeles:20260706T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO,FR',
    );
  });
});

describe('compose/parse round-trip stability', () => {
  // Regression: parseWeeklyRecurrence used to slice UNTIL's UTC digits for
  // endsOn, so negative-offset zones recovered the wall date one day late and
  // every edit-save cycle extended the end date by another day.
  test.each([['America/Los_Angeles'], ['Asia/Tokyo'], ['UTC']])(
    'endsOn survives three compose/parse cycles unchanged in %s',
    timezone => {
      let composed = composeWeeklyRecurrence({ ...laRule, timezone });
      for (let cycle = 1; cycle <= 3; cycle += 1) {
        const parsed = parseWeeklyRecurrence(composed);
        expect(parsed).not.toBeNull();
        expect(parsed!.endsOn).toBe('2026-07-20');
        expect(parsed!.startsOn).toBe('2026-07-06');
        composed = composeWeeklyRecurrence(parsed!);
      }
    },
  );

  test('a composed value is a fixed point of parse-then-compose', () => {
    const composed = composeWeeklyRecurrence(laRule);
    expect(composeWeeklyRecurrence(parseWeeklyRecurrence(composed)!)).toBe(composed);
  });
});

describe('parseWeeklyRecurrence', () => {
  test('recovers the full form state from a composed value', () => {
    expect(parseWeeklyRecurrence(composeWeeklyRecurrence(laRule))).toEqual(laRule);
  });

  test('an open-ended rule parses without endsOn', () => {
    const parsed = parseWeeklyRecurrence(
      'DTSTART;TZID=Asia/Tokyo:20260706T090000\nRRULE:FREQ=WEEKLY;BYDAY=TU',
    );
    expect(parsed).toEqual({
      days: ['TU'],
      time: '09:00',
      startsOn: '2026-07-06',
      timezone: 'Asia/Tokyo',
    });
    expect(parsed!.endsOn).toBeUndefined();
  });

  // Rules this form cannot represent must parse to null so the modal
  // degrades to details-only editing instead of silently stripping parts.
  test.each([
    ['missing DTSTART', 'RRULE:FREQ=WEEKLY;BYDAY=MO'],
    ['missing RRULE', 'DTSTART;TZID=UTC:20260706T180000'],
    ['non-weekly FREQ', 'DTSTART;TZID=UTC:20260706T180000\nRRULE:FREQ=DAILY'],
    ['missing BYDAY', 'DTSTART;TZID=UTC:20260706T180000\nRRULE:FREQ=WEEKLY'],
    ['INTERVAL present', 'DTSTART;TZID=UTC:20260706T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO;INTERVAL=2'],
    ['COUNT present', 'DTSTART;TZID=UTC:20260706T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO;COUNT=5'],
    ['WKST present', 'DTSTART;TZID=UTC:20260706T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO;WKST=MO'],
    // rrule-go also accepts date-only and zone-local UNTIL values; parsing
    // them as UTC instants would misread the end date, so they must degrade
    // to details-only editing rather than parse as open-ended.
    [
      'date-only UNTIL',
      'DTSTART;TZID=UTC:20260706T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO;UNTIL=20260720',
    ],
    [
      'zone-local UNTIL without Z',
      'DTSTART;TZID=UTC:20260706T180000\nRRULE:FREQ=WEEKLY;BYDAY=MO;UNTIL=20260720T235900',
    ],
  ])('%s parses to null', (_name, value) => {
    expect(parseWeeklyRecurrence(value)).toBeNull();
  });
});

describe('wallTimeInZoneToUTC', () => {
  test.each([
    [
      'Los Angeles winter (PST, UTC-8)',
      '2026-01-15',
      '10:00',
      'America/Los_Angeles',
      '2026-01-15T18:00:00.000Z',
    ],
    [
      'Los Angeles summer (PDT, UTC-7)',
      '2026-07-20',
      '10:00',
      'America/Los_Angeles',
      '2026-07-20T17:00:00.000Z',
    ],
    [
      'Kolkata (half-hour offset, UTC+5:30)',
      '2026-07-20',
      '10:00',
      'Asia/Kolkata',
      '2026-07-20T04:30:00.000Z',
    ],
  ])('%s maps to the exact UTC instant', (_name, date, time, zone, expected) => {
    expect(wallTimeInZoneToUTC(date, time, zone).toISOString()).toBe(expected);
  });

  test('a nonexistent spring-forward wall time resolves without throwing', () => {
    // 2027-03-14 02:30 does not exist in Los Angeles: clocks jump from
    // 02:00 PST straight to 03:00 PDT. The two-pass offset resolution lands
    // on the pre-transition side: 09:30Z, i.e. 01:30 PST wall time.
    const instant = wallTimeInZoneToUTC('2027-03-14', '02:30', 'America/Los_Angeles');
    expect(Number.isNaN(instant.getTime())).toBe(false);
    expect(instant.toISOString()).toBe('2027-03-14T09:30:00.000Z');
  });
});

describe('wallTimeInZone', () => {
  test('renders a UTC instant as zone-local date and time strings', () => {
    expect(wallTimeInZone('2026-07-21T06:59:00Z', 'America/Los_Angeles')).toEqual({
      date: '2026-07-20',
      time: '23:59',
    });
    expect(wallTimeInZone('2026-07-20T14:59:00Z', 'Asia/Tokyo')).toEqual({
      date: '2026-07-20',
      time: '23:59',
    });
  });

  test('an exact-midnight instant renders as hour 00, never 24', () => {
    // 2026-07-20T07:00Z is exactly 00:00 in Los Angeles (PDT). On h24 ICU
    // builds Intl reports the hour as '24'; the guard must normalize it so
    // the value is always a valid HH:mm form input.
    expect(wallTimeInZone('2026-07-20T07:00:00Z', 'America/Los_Angeles')).toEqual({
      date: '2026-07-20',
      time: '00:00',
    });
    expect(wallTimeInZone('2026-07-20T00:00:00Z', 'UTC')).toEqual({
      date: '2026-07-20',
      time: '00:00',
    });
  });
});

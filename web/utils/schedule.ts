// Helpers for the scheduled streams admin UI: composing and parsing the
// RFC 5545 recurrence values the backend stores, and converting a wall-clock
// time in an arbitrary IANA timezone to a UTC instant without a timezone
// library.

// WEEKDAY_CODES are the RRULE BYDAY codes in display order. 2024-01-01 was
// a Monday, so day i of that week renders code i's localized label.
const WEEKDAY_CODES = ['MO', 'TU', 'WE', 'TH', 'FR', 'SA', 'SU'];

// WEEKDAYS pairs each RRULE code with a short weekday label in the
// browser's own locale, derived through Intl so the chips localize without
// a translation table.
export const WEEKDAYS = WEEKDAY_CODES.map((code, index) => ({
  code,
  label: new Intl.DateTimeFormat(undefined, { weekday: 'short', timeZone: 'UTC' }).format(
    new Date(Date.UTC(2024, 0, 1 + index)),
  ),
}));

const WEEKDAY_LABELS: Record<string, string> = {
  MO: 'Mon',
  TU: 'Tue',
  WE: 'Wed',
  TH: 'Thu',
  FR: 'Fri',
  SA: 'Sat',
  SU: 'Sun',
};

export interface WeeklyRecurrence {
  // RRULE BYDAY codes, e.g. ['MO', 'FR'].
  days: string[];
  // Wall-clock start time as HH:mm.
  time: string;
  // First date the rule applies from, as YYYY-MM-DD.
  startsOn: string;
  // IANA timezone name the wall times live in.
  timezone: string;
  // Optional last day of the rule, as YYYY-MM-DD.
  endsOn?: string;
}

// composeWeeklyRecurrence builds the DTSTART+RRULE value the backend stores.
// The timezone travels inside the value (TZID) so expansion stays
// DST-correct server-side.
export function composeWeeklyRecurrence(rule: WeeklyRecurrence): string {
  const date = rule.startsOn.replaceAll('-', '');
  const time = `${rule.time.replace(':', '')}00`;
  let rrule = `FREQ=WEEKLY;BYDAY=${rule.days.join(',')}`;
  if (rule.endsOn) {
    // UNTIL is an inclusive UTC instant; the end of the chosen day in the
    // rule's own timezone keeps that day's occurrence included.
    const endOfDay = wallTimeInZoneToUTC(rule.endsOn, '23:59', rule.timezone);
    rrule += `;UNTIL=${endOfDay
      .toISOString()
      .replaceAll(/[-:]/g, '')
      .replace(/\.\d+Z$/, 'Z')}`;
  }
  return `DTSTART;TZID=${rule.timezone}:${date}T${time}\nRRULE:${rrule}`;
}

// SUPPORTED_RRULE_PARTS are the only RRULE components this form can
// represent. Anything else (INTERVAL, COUNT, WKST...) means the rule was
// created elsewhere and editing it here would silently strip parts, so
// parsing refuses and the modal degrades to details-only editing.
const SUPPORTED_RRULE_PARTS: Record<string, true> = {
  FREQ: true,
  BYDAY: true,
  UNTIL: true,
};

// parseWeeklyRecurrence parses a stored recurrence value back into form
// state. Returns null for anything this UI does not produce (non-weekly
// rules or rules carrying unsupported RRULE parts), so callers can degrade
// to details-only editing.
export function parseWeeklyRecurrence(recurrence: string): WeeklyRecurrence | null {
  const dtstart = recurrence.match(/DTSTART;TZID=([^:]+):(\d{8})T(\d{4})\d{2}/);
  const rrule = recurrence.match(/RRULE:(.+)/);
  if (!dtstart || !rrule) {
    return null;
  }

  const parts: Record<string, string> = {};
  rrule[1].split(';').forEach(part => {
    const [key, value] = part.split('=');
    parts[key] = value;
  });
  if (parts.FREQ !== 'WEEKLY' || !parts.BYDAY) {
    return null;
  }
  if (Object.keys(parts).some(key => !SUPPORTED_RRULE_PARTS[key])) {
    return null;
  }

  const [, timezone, date, time] = dtstart;
  const rule: WeeklyRecurrence = {
    days: parts.BYDAY.split(','),
    time: `${time.slice(0, 2)}:${time.slice(2, 4)}`,
    startsOn: `${date.slice(0, 4)}-${date.slice(4, 6)}-${date.slice(6, 8)}`,
    timezone,
  };

  if (parts.UNTIL) {
    // Only the UTC-instant format this form itself produces is editable.
    // rrule-go also accepts date-only and zone-local UNTIL values; parsing
    // those here would misread or silently drop the end date on re-save,
    // so such rules degrade to details-only editing instead.
    const match = parts.UNTIL.match(/^(\d{4})(\d{2})(\d{2})T(\d{2})(\d{2})\d{2}Z$/);
    if (!match) {
      return null;
    }
    // UNTIL is a UTC instant; recover the wall DATE in the rule's own zone.
    // Slicing the UTC digits would land one day late in negative-offset
    // zones, and each edit-save cycle would extend the end date by a day.
    const instant = new Date(
      Date.UTC(
        Number(match[1]),
        Number(match[2]) - 1,
        Number(match[3]),
        Number(match[4]),
        Number(match[5]),
      ),
    );
    rule.endsOn = wallTimeInZone(instant.toISOString(), timezone).date;
  }
  return rule;
}

// describeRecurrence renders a stored recurrence value as a short human
// summary for the admin table, e.g. "Weekly on Mon, Fri at 18:00
// (America/Los_Angeles)". Falls back to the raw value for rules this UI
// does not understand.
export function describeRecurrence(recurrence: string): string {
  const rule = parseWeeklyRecurrence(recurrence);
  if (!rule) {
    return recurrence;
  }
  const days = rule.days.map(code => WEEKDAY_LABELS[code] || code).join(', ');
  const until = rule.endsOn ? ` until ${rule.endsOn}` : '';
  return `Weekly on ${days} at ${rule.time} (${rule.timezone})${until}`;
}

// timezoneOffsetAt returns the UTC offset in milliseconds of an IANA zone at
// a given instant, derived through Intl so no timezone database ships with
// the app.
function timezoneOffsetAt(instantMs: number, timezone: string): number {
  const formatter = new Intl.DateTimeFormat('en-US', {
    timeZone: timezone,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  });
  const parts = formatter.formatToParts(new Date(instantMs));
  const get = (type: string) => Number(parts.find(part => part.type === type)?.value);
  const asUTC = Date.UTC(
    get('year'),
    get('month') - 1,
    get('day'),
    get('hour') % 24,
    get('minute'),
    get('second'),
  );
  return asUTC - instantMs;
}

// wallTimeInZoneToUTC converts "this date and time on the wall clock of this
// IANA zone" into the corresponding UTC instant. Two-pass offset resolution
// handles instants near DST transitions.
export function wallTimeInZoneToUTC(date: string, time: string, timezone: string): Date {
  const [year, month, day] = date.split('-').map(Number);
  const [hour, minute] = time.split(':').map(Number);
  const wallMs = Date.UTC(year, month - 1, day, hour, minute, 0);

  const firstGuess = timezoneOffsetAt(wallMs, timezone);
  let utcMs = wallMs - firstGuess;
  const secondGuess = timezoneOffsetAt(utcMs, timezone);
  if (secondGuess !== firstGuess) {
    utcMs = wallMs - secondGuess;
  }
  return new Date(utcMs);
}

// wallTimeInZone renders a UTC instant as its wall-clock date and time in an
// IANA zone, e.g. for prefilling date/time inputs when editing.
export function wallTimeInZone(instant: string, timezone: string): { date: string; time: string } {
  const formatter = new Intl.DateTimeFormat('en-CA', {
    timeZone: timezone,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });
  const parts = formatter.formatToParts(new Date(instant));
  const get = (type: string) => parts.find(part => part.type === type)?.value || '';
  return {
    date: `${get('year')}-${get('month')}-${get('day')}`,
    time: `${get('hour') === '24' ? '00' : get('hour')}:${get('minute')}`,
  };
}

// timezoneChoices lists IANA zones for the timezone selector, preferring the
// runtime's own database and degrading to a minimal list on old browsers.
export function timezoneChoices(): string[] {
  if (typeof Intl.supportedValuesOf === 'function') {
    const zones = Intl.supportedValuesOf('timeZone');
    // ICU's list omits plain UTC, a natural choice for a stream schedule.
    return zones.includes('UTC') ? zones : ['UTC', ...zones];
  }
  const fallback = Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';
  return fallback === 'UTC' ? ['UTC'] : ['UTC', fallback];
}

var request = require('supertest');
request = request(process.env.OWNCAST_TEST_URL || 'http://127.0.0.1:8080');

const randomString = require('./lib/rand').randomString;

const adminPassword = 'abc123';

// One-off event a day out so it is upcoming but chat is not yet open.
const oneOffName = 'One-off Test Stream ' + randomString(8);
const updatedOneOffName = 'Renamed One-off ' + randomString(8);
const oneOffDescription = 'A single scheduled stream.';
const oneOffTimezone = 'America/New_York';
const oneOffDuration = 90;
const updatedOneOffDuration = 45;
const keepFieldsOneOffName = 'Keep-fields One-off ' + randomString(8);
const oneOffStart = new Date(Date.now() + 24 * 60 * 60 * 1000);
oneOffStart.setMilliseconds(0);

// Weekly series starting two days out, anchored to a wall time in LA.
const seriesName = 'Weekly Test Series ' + randomString(8);
const seriesDescription = 'A recurring scheduled stream.';
const seriesDuration = 120;
const keepFieldsSeriesName = 'Keep-fields Series ' + randomString(8);
const seriesAnchor = new Date(Date.now() + 2 * 24 * 60 * 60 * 1000);
const pad = (n) => String(n).padStart(2, '0');
const seriesDtStart = `${seriesAnchor.getUTCFullYear()}${pad(
	seriesAnchor.getUTCMonth() + 1,
)}${pad(seriesAnchor.getUTCDate())}T190000`;
const seriesRecurrence = `DTSTART;TZID=America/Los_Angeles:${seriesDtStart}\nRRULE:FREQ=WEEKLY`;

const reminderMessage = 'Starting soon: ' + randomString(8);

// State captured as the sequential story progresses.
let oneOffId;
let seriesId;
let cancelledOccurrenceId;

async function postScheduleEvent(payload, expectedStatus = 200) {
	return request
		.post('/api/admin/schedule/event')
		.auth('admin', adminPassword)
		.send(payload)
		.expect(expectedStatus);
}

async function getAdminSchedule() {
	return request
		.get('/api/admin/schedule')
		.auth('admin', adminPassword)
		.expect(200);
}

test('schedule feature is disabled by default', async () => {
	const res = await request.get('/api/config').expect(200);
	expect(res.body.schedule.enabled).toBe(false);
	expect(res.body.schedule.showCountdown).toBe(false);
	expect(res.body.schedule.chatOpenMinutesBefore).toBe(0);
});

test('public schedule is empty while the feature is disabled', async () => {
	const res = await request.get('/api/schedule').expect(200);
	expect(res.body).toEqual([]);
});

test('status carries no scheduled event while the feature is disabled', async () => {
	const res = await request.get('/api/status').expect(200);
	expect(res.body).not.toHaveProperty('scheduledEvent');
});

test('admin schedule requires authentication', async () => {
	const res = await request.get('/api/admin/schedule').expect(401);
	expect(res.headers['www-authenticate']).toContain('Basic realm=');
	expect(res.text).toContain('Unauthorized');
});

test('creating an event requires authentication', async () => {
	const res = await request
		.post('/api/admin/schedule/event')
		.send({ name: oneOffName, start: oneOffStart.toISOString() })
		.expect(401);
	expect(res.text).toContain('Unauthorized');
});

test('enable the schedule feature', async () => {
	const res = await request
		.post('/api/admin/config/schedule/enabled')
		.auth('admin', adminPassword)
		.send({ value: true })
		.expect(200);
	expect(res.body.success).toBe(true);

	const config = await request.get('/api/config').expect(200);
	expect(config.body.schedule.enabled).toBe(true);
});

test('enable the event countdown', async () => {
	const res = await request
		.post('/api/admin/config/schedule/showcountdown')
		.auth('admin', adminPassword)
		.send({ value: true })
		.expect(200);
	expect(res.body.success).toBe(true);

	const config = await request.get('/api/config').expect(200);
	expect(config.body.schedule.showCountdown).toBe(true);
});

test('set the chat open lead time', async () => {
	const res = await request
		.post('/api/admin/config/schedule/chatopenminutes')
		.auth('admin', adminPassword)
		.send({ value: 30 })
		.expect(200);
	expect(res.body.success).toBe(true);

	const config = await request.get('/api/config').expect(200);
	expect(config.body.schedule.chatOpenMinutesBefore).toBe(30);

	const invalid = await request
		.post('/api/admin/config/schedule/chatopenminutes')
		.auth('admin', adminPassword)
		.send({ value: 15 })
		.expect(400);
	expect(invalid.text).toContain('0, 5, 10, 30, or 60');
});

test('set the schedule reminder message', async () => {
	const res = await request
		.post('/api/admin/config/schedule/remindermessage')
		.auth('admin', adminPassword)
		.send({ value: reminderMessage })
		.expect(200);
	expect(res.body.success).toBe(true);
});

test('admin schedule starts with no series and no events', async () => {
	const res = await getAdminSchedule();
	expect(res.body.series).toEqual([]);
	expect(res.body.events).toEqual([]);
});

test('create a one-off event', async () => {
	const res = await postScheduleEvent({
		name: oneOffName,
		description: oneOffDescription,
		start: oneOffStart.toISOString(),
		timezone: oneOffTimezone,
		durationMinutes: oneOffDuration,
	});

	const event = res.body.events.find((e) => e.name === oneOffName);
	expect(event).toBeDefined();
	expect(event.description).toBe(oneOffDescription);
	expect(new Date(event.startTime).getTime()).toBe(oneOffStart.getTime());
	expect(event.durationMinutes).toBe(oneOffDuration);
	expect(event.timezone).toBe(oneOffTimezone);
	expect(event.status).toBe('scheduled');
	expect(event.seriesId).toBeUndefined();

	oneOffId = event.id;
	expect(oneOffId).toBeTruthy();
});

test('one-off event appears in the public schedule', async () => {
	const res = await request.get('/api/schedule').expect(200);
	expect(res.body).toHaveLength(1);

	const event = res.body[0];
	expect(event.id).toBe(oneOffId);
	expect(event.name).toBe(oneOffName);
	expect(new Date(event.startTime).getTime()).toBe(oneOffStart.getTime());
	expect(event.durationMinutes).toBe(oneOffDuration);
	expect(event.timezone).toBe(oneOffTimezone);
	expect(event.status).toBe('scheduled');
});

test('public iCalendar feed returns the current schedule', async () => {
	const res = await request
		.get('/api/schedule.ics')
		.expect('content-type', /text\/calendar/)
		.expect('cache-control', /no-cache/)
		.expect(200);

	expect(res.text).toContain('BEGIN:VCALENDAR\r\n');
	expect(res.text).toContain(`UID:${oneOffId}@owncast\r\n`);
	expect(res.text).toContain(`SUMMARY:${oneOffName}\r\n`);
	expect(res.text).toContain('DTSTART:');
	expect(res.text).toContain('END:VCALENDAR\r\n');
});

test('status shows the upcoming one-off with chat closed', async () => {
	const res = await request.get('/api/status').expect(200);
	expect(res.body.scheduledEvent).toBeDefined();
	expect(res.body.scheduledEvent.id).toBe(oneOffId);
	expect(res.body.scheduledEvent.name).toBe(oneOffName);
	expect(res.body.scheduledEvent.description).toBe(oneOffDescription);
	expect(res.body.scheduledEvent.durationMinutes).toBe(oneOffDuration);
	expect(new Date(res.body.scheduledEvent.startTime).getTime()).toBe(
		oneOffStart.getTime(),
	);
	expect(res.body.scheduledEvent.chatOpen).toBe(false);
});

test('create a weekly series', async () => {
	const res = await postScheduleEvent({
		name: seriesName,
		description: seriesDescription,
		recurrence: seriesRecurrence,
		durationMinutes: seriesDuration,
	});

	const series = res.body.series.find((s) => s.name === seriesName);
	expect(series).toBeDefined();
	expect(series.description).toBe(seriesDescription);
	expect(series.recurrence).toBe(seriesRecurrence);
	expect(series.durationMinutes).toBe(seriesDuration);
	expect(series.active).toBe(true);

	seriesId = series.id;
	expect(seriesId).toBeTruthy();

	const occurrences = res.body.events.filter((e) => e.seriesId === seriesId);
	expect(occurrences.length).toBeGreaterThanOrEqual(3);
	occurrences.forEach((o) => {
		expect(o.name).toBe(seriesName);
		expect(o.durationMinutes).toBe(seriesDuration);
		expect(o.status).toBe('scheduled');
		expect(new Date(o.startTime).getTime()).toBeGreaterThan(Date.now());
	});

	occurrences.sort(
		(a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime(),
	);
	cancelledOccurrenceId = occurrences[0].id;
});

test('series occurrences appear in the public schedule with a seriesId', async () => {
	const res = await request.get('/api/schedule').expect(200);

	const occurrences = res.body.filter((e) => e.seriesId === seriesId);
	expect(occurrences.length).toBeGreaterThanOrEqual(3);
	occurrences.forEach((o) => {
		expect(o.seriesId).toBe(seriesId);
		expect(o.name).toBe(seriesName);
		expect(o.status).toBe('scheduled');
	});

	// The one-off is still listed alongside the series occurrences.
	expect(res.body.find((e) => e.id === oneOffId)).toBeDefined();
});

test('reject an event with both a start and a recurrence', async () => {
	const res = await postScheduleEvent(
		{
			name: randomString(),
			start: oneOffStart.toISOString(),
			recurrence: seriesRecurrence,
		},
		400,
	);
	expect(res.text).toContain('mutually exclusive');
});

test('reject an invalid recurrence rule', async () => {
	const res = await postScheduleEvent(
		{
			name: randomString(),
			recurrence: `DTSTART;TZID=America/Los_Angeles:${seriesDtStart}\nRRULE:FREQ=BOGUS`,
		},
		400,
	);
	expect(res.text.length).toBeGreaterThan(0);
});

test('reject a recurrence with no DTSTART', async () => {
	const res = await postScheduleEvent(
		{
			name: randomString(),
			recurrence: 'RRULE:FREQ=WEEKLY',
		},
		400,
	);
	expect(res.text).toContain('DTSTART');
});

test('reject an event with no name', async () => {
	const res = await postScheduleEvent(
		{
			start: oneOffStart.toISOString(),
		},
		400,
	);
	expect(res.text).toContain('name is required');
});

test('reject an event starting in the past', async () => {
	const res = await postScheduleEvent(
		{
			name: randomString(),
			start: new Date(Date.now() - 60 * 60 * 1000).toISOString(),
		},
		400,
	);
	expect(res.text).toContain('start must be in the future');
});

test('reject an event with an unknown timezone', async () => {
	const res = await postScheduleEvent(
		{
			name: randomString(),
			start: oneOffStart.toISOString(),
			timezone: 'Mars/OlympusMons',
		},
		400,
	);
	expect(res.text).toContain('unknown timezone');
});

test('recurrence preview rejects an invalid rule', async () => {
	const res = await request
		.post('/api/admin/schedule/preview')
		.auth('admin', adminPassword)
		.send({ recurrence: 'RRULE:FREQ=WEEKLY' })
		.expect(400);
	expect(res.text).toContain('DTSTART');
});

test('recurrence preview returns at most five upcoming ISO occurrences', async () => {
	const res = await request
		.post('/api/admin/schedule/preview')
		.auth('admin', adminPassword)
		.send({ recurrence: seriesRecurrence })
		.expect(200);

	// A weekly rule expands to more than five occurrences in a year, so the
	// preview cap is exercised exactly.
	expect(res.body.occurrences).toHaveLength(5);

	let previous = Date.now();
	res.body.occurrences.forEach((occurrence) => {
		const parsed = Date.parse(occurrence);
		expect(Number.isFinite(parsed)).toBe(true);
		expect(parsed).toBeGreaterThan(previous);
		previous = parsed;
	});
});

test('edit the one-off name and duration', async () => {
	const res = await postScheduleEvent({
		id: oneOffId,
		name: updatedOneOffName,
		description: oneOffDescription,
		durationMinutes: updatedOneOffDuration,
	});

	const event = res.body.events.find((e) => e.id === oneOffId);
	expect(event).toBeDefined();
	expect(event.name).toBe(updatedOneOffName);
	expect(event.durationMinutes).toBe(updatedOneOffDuration);
	expect(new Date(event.startTime).getTime()).toBe(oneOffStart.getTime());
	expect(event.timezone).toBe(oneOffTimezone);
	expect(event.status).toBe('scheduled');

	const publicRes = await request.get('/api/schedule').expect(200);
	const publicEvent = publicRes.body.find((e) => e.id === oneOffId);
	expect(publicEvent.name).toBe(updatedOneOffName);
	expect(publicEvent.durationMinutes).toBe(updatedOneOffDuration);
});

test('cancel a single series occurrence', async () => {
	const res = await request
		.post('/api/admin/schedule/event/delete')
		.auth('admin', adminPassword)
		.send({ id: cancelledOccurrenceId, cancel: true })
		.expect(200);

	const cancelled = res.body.events.find((e) => e.id === cancelledOccurrenceId);
	expect(cancelled).toBeDefined();
	expect(cancelled.status).toBe('cancelled');
	expect(cancelled.seriesId).toBe(seriesId);
});

test('reject a series update that carries a start', async () => {
	const res = await postScheduleEvent(
		{
			id: seriesId,
			name: seriesName,
			start: oneOffStart.toISOString(),
		},
		400,
	);
	expect(res.text).toContain('has no single start');
});

test('reject turning a one-off event into a series', async () => {
	const res = await postScheduleEvent(
		{
			id: oneOffId,
			name: updatedOneOffName,
			recurrence: seriesRecurrence,
		},
		400,
	);
	expect(res.text).toContain('cannot become a recurring series');
});

test('reject changing the timezone on an update', async () => {
	const res = await postScheduleEvent(
		{
			id: oneOffId,
			name: updatedOneOffName,
			timezone: 'UTC',
		},
		400,
	);
	expect(res.text).toContain('timezone cannot be changed');
});

test('a rejected past-start update leaves the one-off untouched', async () => {
	const res = await postScheduleEvent(
		{
			id: oneOffId,
			name: 'Should Not Persist ' + randomString(8),
			description: 'should not persist',
			durationMinutes: 30,
			start: new Date(Date.now() - 60 * 60 * 1000).toISOString(),
		},
		400,
	);
	expect(res.text).toContain('start must be in the future');

	// Validation happens before the first write, so none of the details
	// riding along with the rejected start may stick.
	const schedule = await getAdminSchedule();
	const event = schedule.body.events.find((e) => e.id === oneOffId);
	expect(event).toBeDefined();
	expect(event.name).toBe(updatedOneOffName);
	expect(event.description).toBe(oneOffDescription);
	expect(event.durationMinutes).toBe(updatedOneOffDuration);
	expect(new Date(event.startTime).getTime()).toBe(oneOffStart.getTime());
});

test('updating a one-off with only a name keeps its other details', async () => {
	const res = await postScheduleEvent({
		id: oneOffId,
		name: keepFieldsOneOffName,
	});

	const event = res.body.events.find((e) => e.id === oneOffId);
	expect(event).toBeDefined();
	expect(event.name).toBe(keepFieldsOneOffName);
	expect(event.description).toBe(oneOffDescription);
	expect(event.durationMinutes).toBe(updatedOneOffDuration);
	expect(new Date(event.startTime).getTime()).toBe(oneOffStart.getTime());
	expect(event.timezone).toBe(oneOffTimezone);
});

test('updating a series with only a name keeps its other details', async () => {
	const res = await postScheduleEvent({
		id: seriesId,
		name: keepFieldsSeriesName,
	});

	const series = res.body.series.find((s) => s.id === seriesId);
	expect(series).toBeDefined();
	expect(series.name).toBe(keepFieldsSeriesName);
	expect(series.description).toBe(seriesDescription);
	expect(series.durationMinutes).toBe(seriesDuration);
	expect(series.recurrence).toBe(seriesRecurrence);

	// Regenerated occurrences pick up the new name while the cancelled
	// one is left alone.
	const occurrences = res.body.events.filter((e) => e.seriesId === seriesId);
	occurrences
		.filter((o) => o.status === 'scheduled')
		.forEach((o) => expect(o.name).toBe(keepFieldsSeriesName));
	const cancelled = occurrences.find((o) => o.id === cancelledOccurrenceId);
	expect(cancelled).toBeDefined();
	expect(cancelled.status).toBe('cancelled');
});

test('delete the series', async () => {
	const res = await request
		.post('/api/admin/schedule/event/delete')
		.auth('admin', adminPassword)
		.send({ id: seriesId })
		.expect(200);

	expect(res.body.series.find((s) => s.id === seriesId)).toBeUndefined();

	// Every scheduled occurrence of the series is gone. The cancelled one
	// stays for the record.
	const remaining = res.body.events.filter((e) => e.seriesId === seriesId);
	remaining.forEach((e) => expect(e.status).toBe('cancelled'));
	expect(remaining.map((e) => e.id)).toContain(cancelledOccurrenceId);

	// The one-off is unaffected.
	const oneOff = res.body.events.find((e) => e.id === oneOffId);
	expect(oneOff).toBeDefined();
	expect(oneOff.status).toBe('scheduled');
});

test('public schedule window far in the future is empty', async () => {
	const from = new Date(Date.now() + 200 * 24 * 60 * 60 * 1000).toISOString();
	const to = new Date(Date.now() + 210 * 24 * 60 * 60 * 1000).toISOString();
	const res = await request
		.get('/api/schedule')
		.query({ from, to })
		.expect(200);
	expect(res.body).toEqual([]);
});

test('public schedule rejects an inverted range', async () => {
	const from = new Date(Date.now() + 10 * 24 * 60 * 60 * 1000).toISOString();
	const to = new Date(Date.now() + 5 * 24 * 60 * 60 * 1000).toISOString();
	const res = await request
		.get('/api/schedule')
		.query({ from, to })
		.expect(400);
	expect(res.text).toContain('to must be after from');
});

test('public schedule rejects an oversized range', async () => {
	const from = new Date().toISOString();
	const to = new Date(Date.now() + 400 * 24 * 60 * 60 * 1000).toISOString();
	const res = await request
		.get('/api/schedule')
		.query({ from, to })
		.expect(400);
	expect(res.text).toContain('too large');
});

test('disable the schedule feature', async () => {
	const res = await request
		.post('/api/admin/config/schedule/enabled')
		.auth('admin', adminPassword)
		.send({ value: false })
		.expect(200);
	expect(res.body.success).toBe(true);

	const config = await request.get('/api/config').expect(200);
	expect(config.body.schedule.enabled).toBe(false);

	const publicSchedule = await request.get('/api/schedule').expect(200);
	expect(publicSchedule.body).toEqual([]);

	const status = await request.get('/api/status').expect(200);
	expect(status.body).not.toHaveProperty('scheduledEvent');
});

test('admin schedule data is preserved while the feature is disabled', async () => {
	const res = await getAdminSchedule();

	const oneOff = res.body.events.find((e) => e.id === oneOffId);
	expect(oneOff).toBeDefined();
	expect(oneOff.name).toBe(keepFieldsOneOffName);
	expect(oneOff.status).toBe('scheduled');
});

const { test, expect } = require('@jest/globals');
var request = require('supertest');
request = request('http://127.0.0.1:8080');

const registerChat = require('./lib/chat').registerChat;
const sendAdminPayload = require('./lib/admin').sendAdminPayload;
const getAdminResponse = require('./lib/admin').getAdminResponse;

var userId;
var displayName;

test('register a user', async () => {
	const registration = await registerChat();
	userId = registration.id;
	displayName = registration.displayName;
	expect(userId).toBeTruthy();
});

test('paginated user list returns total and results', async () => {
	const response = await getAdminResponse('users?offset=0&limit=50');

	expect(typeof response.body.total).toBe('number');
	expect(response.body.total).toBeGreaterThan(0);
	expect(Array.isArray(response.body.results)).toBe(true);

	// The freshly registered user is the most recently created, so it is on
	// the first page (results are ordered created_at DESC).
	const match = response.body.results.filter((user) => user.id === userId);
	expect(match).toHaveLength(1);
});

test('pagination limit is respected', async () => {
	const response = await getAdminResponse('users?offset=0&limit=1');

	// total reflects the full count, not just this page.
	expect(response.body.total).toBeGreaterThan(0);
	expect(response.body.results.length).toBeLessThanOrEqual(1);
});

test('search filters users by display name', async () => {
	const response = await getAdminResponse(
		'users?search=' + encodeURIComponent(displayName),
	);

	// Every returned user's display name contains the search term, and our
	// user is among them.
	for (const user of response.body.results) {
		expect(user.displayName.toLowerCase()).toContain(displayName.toLowerCase());
	}
	const match = response.body.results.filter((user) => user.id === userId);
	expect(match).toHaveLength(1);
});

test('authProviders is empty for an anonymous user', async () => {
	const response = await getAdminResponse(
		'users?search=' + encodeURIComponent(displayName),
	);
	const match = response.body.results.filter((user) => user.id === userId);
	expect(match).toHaveLength(1);
	expect(match[0].authProviders || []).toHaveLength(0);
});

test('sort orders the page by creation date in both directions', async () => {
	const asc = await getAdminResponse('users?offset=0&limit=50&sort=asc');
	const desc = await getAdminResponse('users?offset=0&limit=50&sort=desc');

	const ascTimes = asc.body.results.map((user) =>
		new Date(user.createdAt).getTime(),
	);
	const descTimes = desc.body.results.map((user) =>
		new Date(user.createdAt).getTime(),
	);

	// Each page is monotonically ordered by creation date.
	for (let i = 1; i < ascTimes.length; i++) {
		expect(ascTimes[i]).toBeGreaterThanOrEqual(ascTimes[i - 1]);
	}
	for (let i = 1; i < descTimes.length; i++) {
		expect(descTimes[i]).toBeLessThanOrEqual(descTimes[i - 1]);
	}

	// Omitting sort keeps the default newest-first order, i.e. sort=desc.
	const def = await getAdminResponse('users?offset=0&limit=50');
	expect(def.body.results.map((user) => user.id)).toEqual(
		desc.body.results.map((user) => user.id),
	);
});

// Moderators and banned users come from their own endpoints, which return the
// complete set rather than a page of the users API.
test('moderators list includes a promoted moderator', async () => {
	await sendAdminPayload('chat/users/setmoderator', {
		userId: userId,
		isModerator: true,
	});
	const response = await getAdminResponse('chat/users/moderators');
	const match = response.body.filter((user) => user.id === userId);
	expect(match).toHaveLength(1);
});

test('banned list includes a disabled user', async () => {
	await sendAdminPayload('chat/users/setenabled', {
		userId: userId,
		enabled: false,
	});
	const response = await getAdminResponse('chat/users/disabled');
	const match = response.body.filter((user) => user.id === userId);
	expect(match).toHaveLength(1);
});

test('status=bots excludes a normal user', async () => {
	const response = await getAdminResponse(
		'users?status=bots&search=' + encodeURIComponent(displayName),
	);
	const match = response.body.results.filter((user) => user.id === userId);
	expect(match).toHaveLength(0);
});

test('delete a user by admin', async () => {
	const res = await sendAdminPayload('users/delete', { userId: userId });
	expect(res.body.success).toBe(true);
});

test('deleted user no longer appears in the user list', async () => {
	const response = await getAdminResponse(
		'users?search=' + encodeURIComponent(displayName),
	);
	const match = response.body.results.filter((user) => user.id === userId);
	expect(match).toHaveLength(0);
});

test('deleting an unknown user fails', async () => {
	const res = await request
		.post('/api/admin/users/delete')
		.auth('admin', 'abc123')
		.send({ userId: 'this-user-does-not-exist' })
		.expect(400);
	expect(res.body.success).toBe(false);
});

test('deleting with no userId fails', async () => {
	const res = await request
		.post('/api/admin/users/delete')
		.auth('admin', 'abc123')
		.send({})
		.expect(400);
	expect(res.body.success).toBe(false);
});

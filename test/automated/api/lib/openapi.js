// Shared OpenAPI 3.1 response validation for the automated API suite.
// The repo spec is OpenAPI 3.1.0, which jest-openapi does not support. All
// $refs in openapi.yaml are internal ('#/components/...'), so no dereferencing
// library is needed: we register the whole spec document with Ajv (2020-12
// dialect, matching 3.1) and compile validators by JSON-pointer reference.
const fs = require('fs');
const path = require('path');
const yaml = require('js-yaml');
const Ajv2020 = require('ajv/dist/2020');

const SPEC_PATH = path.resolve(__dirname, '../../../../openapi.yaml');
const RESULTS_DIR = path.resolve(__dirname, '../results');
const EXERCISED_FILE = path.join(RESULTS_DIR, 'openapi-exercised.jsonl');
const SPEC_ID = 'owncast-openapi';

// Known spec drift where the live response is suspect and product code must
// not change. Validation is skipped for matching responses; the pair still
// counts as exercised. Keep entries narrow: exact spec path, method, status.
const EXCLUSIONS = [
	// (none yet)
];

const METHODS = ['get', 'post', 'put', 'delete', 'patch', 'head', 'options'];

// validateSchema:false because the registered document is a whole OpenAPI
// spec, not a JSON Schema; strict:false tolerates OpenAPI-only keywords.
const ajv = new Ajv2020({
	strict: false,
	validateFormats: false,
	allErrors: true,
	validateSchema: false,
});
const compiled = new Map(); // schema ref -> validate fn (or Error)

let routes = null;

function escapePointerSegment(seg) {
	return seg.replace(/~/g, '~0').replace(/\//g, '~1');
}

function resolvePointer(spec, pointer) {
	return pointer
		.split('/')
		.slice(1)
		.reduce(
			(node, seg) => node && node[seg.replace(/~1/g, '/').replace(/~0/g, '~')],
			spec,
		);
}

// Follows $ref chains, returning the resolved node and its JSON pointer
// within the spec (null pointer means the node sits at its original spot).
function deref(spec, node) {
	let pointer = null;
	while (node && typeof node === 'object' && typeof node.$ref === 'string') {
		pointer = node.$ref.slice(1);
		node = resolvePointer(spec, pointer);
	}
	return { node, pointer };
}

function loadRoutes() {
	if (routes) return routes;
	const spec = yaml.load(fs.readFileSync(SPEC_PATH, 'utf8'));
	ajv.addSchema(spec, SPEC_ID);
	routes = [];
	for (const [tmpl, item] of Object.entries(spec.paths || {})) {
		for (const method of METHODS) {
			if (!item[method]) continue;
			const responses = {};
			for (const [status, respRaw] of Object.entries(
				item[method].responses || {},
			)) {
				const { node: resp, pointer } = deref(spec, respRaw);
				if (!resp) continue;
				const basePointer =
					pointer ||
					`/paths/${escapePointerSegment(tmpl)}/${method}/responses/${status}`;
				const media = resp.content && resp.content['application/json'];
				responses[status] = {
					schemaRef:
						media && media.schema
							? `${SPEC_ID}#${basePointer}/content/application~1json/schema`
							: null,
				};
			}
			const escaped = tmpl
				.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
				.replace(/\\\{[^/]+?\\\}/g, '[^/]+');
			routes.push({
				method: method.toUpperCase(),
				path: tmpl,
				regex: new RegExp('^' + escaped + '$'),
				responses,
			});
		}
	}
	return routes;
}

// The spec's server URL is http://localhost:8080/api, so spec paths are
// relative to /api. Requests outside /api are not documented here.
function matchRoute(allRoutes, method, requestPath) {
	if (!requestPath.startsWith('/api/')) return null;
	const rel = requestPath.slice(4);
	const m = method.toUpperCase();
	let templated = null;
	for (const route of allRoutes) {
		if (route.method !== m) continue;
		if (route.path === rel) return route;
		if (!templated && route.regex.test(rel)) templated = route;
	}
	return templated;
}

function isExcluded(route, status) {
	return EXCLUSIONS.some(
		(e) =>
			e.path === route.path &&
			e.method === route.method &&
			(e.status === undefined || String(e.status) === String(status)),
	);
}

function validatorFor(schemaRef) {
	let v = compiled.get(schemaRef);
	if (v === undefined) {
		try {
			v = ajv.getSchema(schemaRef) || ajv.compile({ $ref: schemaRef });
		} catch (err) {
			v = err;
		}
		compiled.set(schemaRef, v);
	}
	return v;
}

function recordExercised(route) {
	try {
		fs.appendFileSync(EXERCISED_FILE, `${route.method} ${route.path}\n`);
	} catch (err) {
		// results dir missing (running outside jest); inventory is best-effort
	}
}

// Validates one supertest/superagent response against the spec. Pushes a
// human-readable string per problem into `failures`. Never throws.
function checkResponse(allRoutes, method, url, res, failures) {
	let requestPath;
	try {
		requestPath = new URL(url).pathname;
	} catch (err) {
		requestPath = String(url).split('?')[0];
	}
	const route = matchRoute(allRoutes, method, requestPath);
	if (!route) return;
	recordExercised(route);
	if (isExcluded(route, res.status)) return;
	// Operations documenting no responses at all (e.g. the /admin/prometheus
	// proxy passthrough) have no contract to check statuses against.
	if (Object.keys(route.responses).length === 0) return;

	const where = `${route.method} ${route.path} -> ${res.status}`;
	const resp = route.responses[String(res.status)] || route.responses.default;
	if (!resp) {
		failures.push(`${where}: status not documented in openapi.yaml`);
		return;
	}
	if (!resp.schemaRef) return;
	if (!/json/.test(res.type || '')) return;

	const validate = validatorFor(resp.schemaRef);
	if (validate instanceof Error) {
		failures.push(`${where}: schema failed to compile: ${validate.message}`);
		return;
	}
	if (!validate(res.body)) {
		const details = (validate.errors || [])
			.slice(0, 5)
			.map((e) => `${e.instancePath || '/'} ${e.message}`)
			.join('; ');
		failures.push(`${where}: body does not match schema: ${details}`);
	}
}

module.exports = {
	loadRoutes,
	matchRoute,
	checkResponse,
	EXCLUSIONS,
	EXERCISED_FILE,
	RESULTS_DIR,
};

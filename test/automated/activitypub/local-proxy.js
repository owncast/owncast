#!/usr/bin/env node

/**
 * Simple HTTPS reverse proxy for local ActivityPub testing.
 *
 * Routes requests based on hostname:
 * - owncast.local -> localhost:8080
 * - snac.local -> localhost:8443 (HTTPS)
 *
 * Usage:
 *   node local-proxy.js [options]
 *
 * Options:
 *   --port=N          Proxy listen port (default: 443)
 *   --cert=FILE       Path to certificate file
 *   --key=FILE        Path to key file
 *   --owncast-port=N  Owncast backend port (default: 8080)
 *   --snac-port=N     snac2 backend port (default: 8443)
 */

const https = require('https');
const http = require('http');
const fs = require('fs');
const { execSync } = require('child_process');
const path = require('path');

// Parse arguments
const args = process.argv.slice(2);
const getArg = (name, defaultVal) => {
	const arg = args.find((a) => a.startsWith(`--${name}=`));
	return arg ? arg.split('=')[1] : defaultVal;
};

const CONFIG = {
	port: parseInt(getArg('port', '443')),
	certFile: getArg('cert', ''),
	keyFile: getArg('key', ''),
	owncastPort: parseInt(getArg('owncast-port', '8080')),
	snacPort: parseInt(getArg('snac-port', '8443')),
	owncastHost: 'owncast.local',
	snacHost: 'snac.local',
};

// HTTP agents for connection pooling
const httpAgent = new http.Agent({
	keepAlive: true,
	maxSockets: 100,
	maxFreeSockets: 10,
	timeout: 30000,
});

const httpsAgent = new https.Agent({
	keepAlive: true,
	maxSockets: 100,
	maxFreeSockets: 10,
	timeout: 30000,
	rejectUnauthorized: false,
});

// Generate self-signed cert if not provided
function generateCert() {
	const tempDir = fs.mkdtempSync(path.join(require('os').tmpdir(), 'proxy-'));
	const keyPath = path.join(tempDir, 'key.pem');
	const certPath = path.join(tempDir, 'cert.pem');

	// Generate cert for both hostnames
	const cmd = `openssl req -x509 -newkey rsa:2048 \
		-keyout ${keyPath} -out ${certPath} \
		-days 1 -nodes \
		-subj "/CN=localhost" \
		-addext "subjectAltName=DNS:${CONFIG.owncastHost},DNS:${CONFIG.snacHost},DNS:localhost,IP:127.0.0.1" \
		2>/dev/null`;

	execSync(cmd);

	return { keyPath, certPath, tempDir };
}

// Proxy a request to a backend
function proxyRequest(clientReq, clientRes, targetHost, targetPort, useHttps) {
	const transport = useHttps ? https : http;
	const agent = useHttps ? httpsAgent : httpAgent;

	const options = {
		hostname: targetHost,
		port: targetPort,
		path: clientReq.url,
		method: clientReq.method,
		headers: {
			...clientReq.headers,
			// Preserve original Host header for HTTP signature verification
		},
		agent: agent,
		timeout: 30000,
	};

	// Remove connection-specific headers that would interfere with keep-alive
	delete options.headers['connection'];
	delete options.headers['keep-alive'];

	const proxyReq = transport.request(options, (proxyRes) => {
		clientRes.writeHead(proxyRes.statusCode, proxyRes.headers);
		proxyRes.pipe(clientRes);
	});

	proxyReq.on('timeout', () => {
		proxyReq.destroy(new Error('Request timeout'));
	});

	proxyReq.on('error', (err) => {
		console.error(`Proxy error to ${targetHost}:${targetPort}: ${err.message}`);
		if (!clientRes.headersSent) {
			// Return a JSON error response for ActivityPub clients
			clientRes.writeHead(502, { 'Content-Type': 'application/json' });
			clientRes.end(
				JSON.stringify({ error: 'Bad Gateway', message: err.message }),
			);
		} else {
			clientRes.end();
		}
	});

	clientReq.pipe(proxyReq);
}

function handleRequest(req, res) {
	const host = (req.headers.host || '').split(':')[0].toLowerCase();

	// Debug logging for federation endpoints only
	if (req.url.includes('/federation') || req.url.includes('/.well-known')) {
		console.log(`[Proxy] ${req.method} ${host}${req.url}`);
	}

	if (host === CONFIG.owncastHost || host === 'localhost') {
		// Route to Owncast (HTTP backend)
		proxyRequest(req, res, 'localhost', CONFIG.owncastPort, false);
	} else if (host === CONFIG.snacHost) {
		// Route to snac2 (HTTP backend - snac2 httpd runs HTTP by default)
		proxyRequest(req, res, 'localhost', CONFIG.snacPort, false);
	} else {
		console.log(`[Proxy] Unknown host: ${host}`);
		res.writeHead(404);
		res.end(`Unknown host: ${host}`);
	}
}

function main() {
	let certPath, keyPath, tempDir;

	if (CONFIG.certFile && CONFIG.keyFile) {
		certPath = CONFIG.certFile;
		keyPath = CONFIG.keyFile;
	} else {
		console.log('[Proxy] Generating self-signed certificate...');
		const certs = generateCert();
		certPath = certs.certPath;
		keyPath = certs.keyPath;
		tempDir = certs.tempDir;
	}

	const serverOptions = {
		key: fs.readFileSync(keyPath),
		cert: fs.readFileSync(certPath),
	};

	const server = https.createServer(serverOptions, handleRequest);

	server.on('error', (err) => {
		if (err.code === 'EACCES') {
			console.error(
				`[Proxy] Cannot bind to port ${CONFIG.port}. Try running with sudo or use --port=8443`,
			);
		} else {
			console.error(`[Proxy] Server error: ${err.message}`);
		}
		process.exit(1);
	});

	server.listen(CONFIG.port, () => {
		console.log(`[Proxy] HTTPS reverse proxy running on port ${CONFIG.port}`);
		console.log(`[Proxy] Routes:`);
		console.log(
			`  - https://${CONFIG.owncastHost}:${CONFIG.port}/ -> http://localhost:${CONFIG.owncastPort}/`,
		);
		console.log(
			`  - https://${CONFIG.snacHost}:${CONFIG.port}/ -> http://localhost:${CONFIG.snacPort}/`,
		);
		console.log('');
		console.log('[Proxy] Make sure /etc/hosts contains:');
		console.log(`  127.0.0.1 ${CONFIG.owncastHost} ${CONFIG.snacHost}`);
	});

	// Cleanup on exit
	process.on('SIGINT', () => {
		console.log('\n[Proxy] Shutting down...');
		server.close();
		if (tempDir) {
			fs.rmSync(tempDir, { recursive: true });
		}
		process.exit(0);
	});

	process.on('SIGTERM', () => {
		server.close();
		if (tempDir) {
			fs.rmSync(tempDir, { recursive: true });
		}
		process.exit(0);
	});
}

// Export for use as module
module.exports = {
	CONFIG,
	generateCert,
	start: main,
};

// Run if executed directly
if (require.main === module) {
	main();
}

vcl 4.0;

backend default {
	.host = "localhost";
	.port = "8080";
}

sub vcl_recv {
	# Implementing websocket support (https://www.varnish-cache.org/docs/4.0/users-guide/vcl-example-websockets.html)
	if (req.http.Upgrade ~ "(?i)websocket") {
		return (pipe);
	}
}

sub vcl_pipe {
	if (req.http.upgrade) {
		set bereq.http.upgrade = req.http.upgrade;
		set bereq.http.connection = req.http.connection;
	}
}

sub vcl_backend_response {
	# Set 1s ttl if origin response HTTP status code is anything other than 200
	if (beresp.status != 200) {
		set beresp.ttl = 1s;
		set beresp.uncacheable = true;
		return (deliver);
	}
	if (bereq.url ~ "m3u8") {
		# assuming chunks are 2 seconds long
		set beresp.ttl = 1s;
		set beresp.grace = 0s;
	}
	if (bereq.url ~ "ts") {
		set beresp.ttl = 10m;
		set beresp.grace = 5m;
	}
}

// Internal OAuth2 / OIDC provider for the OAuth auth-gate integration test.
//
// This is NOT a hand-rolled fake: it is node-oidc-provider, a certified,
// production-grade OpenID Connect provider. It enforces real client
// registration, redirect_uri matching, and client authentication, and issues
// real RS256-signed tokens — so the test exercises the plugin's OAuth *client*
// against a spec-compliant server, not a stub.
//
// The only test-shaped parts are the bits a hermetic test must control:
//   - a single pre-registered confidential client (the plugin),
//   - a fixed account, and
//   - a scripted interaction handler that auto-approves login + consent so no
//     human has to click through the (real) consent flow.
//
// Routes are mapped to the conventional /authorize, /token, /userinfo paths.
import { createServer } from "node:http";
import Provider from "oidc-provider";

const PORT = Number(process.env.FAKE_OAUTH_PORT || 9876);
const ISSUER = `http://localhost:${PORT}`;
const CLIENT_ID = "owncast-oauth-test";
const CLIENT_SECRET = "owncast-oauth-test-secret";
const REDIRECT_URI = "http://localhost:8080/plugins/oauth-gate-test/callback";
const ACCOUNT = "oauth-test-user";

const configuration = {
  clients: [
    {
      client_id: CLIENT_ID,
      client_secret: CLIENT_SECRET,
      grant_types: ["authorization_code"],
      response_types: ["code"],
      redirect_uris: [REDIRECT_URI],
      token_endpoint_auth_method: "client_secret_post",
    },
  ],
  pkce: { required: () => false },
  routes: { authorization: "/authorize", token: "/token", userinfo: "/userinfo" },
  features: { devInteractions: { enabled: false } },
  claims: { openid: ["sub"], profile: ["name", "preferred_username"] },
  cookies: { keys: ["oauth-gate-test-cookie-key"] },
  findAccount: (_ctx, id) => ({
    accountId: id,
    claims: async () => ({
      sub: id,
      name: "OAuth Test User",
      preferred_username: "oauthtester",
    }),
  }),
};

const provider = new Provider(ISSUER, configuration);
const oidc = provider.callback();

// Auto-complete the (real) interaction: log in the fixed account, then grant
// the requested scopes. node-oidc-provider drives this as two prompts (login,
// then consent); the cookie-following client walks both redirects.
async function handleInteraction(req, res) {
  const details = await provider.interactionDetails(req, res);
  const { prompt, params, grantId } = details;

  if (prompt.name === "login") {
    return provider.interactionFinished(
      req,
      res,
      { login: { accountId: ACCOUNT } },
      { mergeWithLastSubmission: false },
    );
  }

  const grant = grantId
    ? await provider.Grant.find(grantId)
    : new provider.Grant({ accountId: ACCOUNT, clientId: params.client_id });
  if (params.scope) grant.addOIDCScope(params.scope);
  const newGrantId = await grant.save();
  return provider.interactionFinished(
    req,
    res,
    { consent: { grantId: newGrantId } },
    { mergeWithLastSubmission: true },
  );
}

const server = createServer((req, res) => {
  const url = new URL(req.url, ISSUER);
  if (url.pathname.startsWith("/interaction/")) {
    handleInteraction(req, res).catch((e) => {
      res.statusCode = 500;
      res.end(String(e));
    });
    return;
  }
  oidc(req, res);
});

server.listen(PORT, "localhost", () => {
  console.log(`[fake-oauth] node-oidc-provider listening at ${ISSUER}`);
});

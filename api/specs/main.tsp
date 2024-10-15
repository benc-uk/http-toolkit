// ==================================================================
// TypeSpec for HTTP Toolkit
// Compile by running `generate.sh`
// ==================================================================

import "@typespec/http";
import "@typespec/openapi";
using TypeSpec.Http;
using TypeSpec.OpenAPI;

@service({
  title: "HTTP Toolkit",
})
@info({
  version: "1.4.0",
})

namespace HTTPToolkit;

@doc("Vanilla text/plain response")
model PlainText {
  @header contentType: "text/plain";
  @body body: string;
}

@doc("Simple OK response")
model OK {
  @header contentType: "text/plain";
  message: string = "OK";
}

@doc("System information")
model SystemInfo {
  hostname: string;
  os: string;
  architecture: string;
  cpuCount: integer;
  memory: string;
  goVersion: string;
  clientAddr: string;
  serverHost: string;
  uptime: string;
}

@doc("Details of an incoming HTTP request")
model RequestInfo {
  method: string;
  path: string;
  remoteAddr: string;
  headers: Record<string>;
  query?: Record<string>;
  body?: string;
  timestamp: string;
}

@tag("Base Routes")
interface Base {
  @route("/")
  @doc("Get a simple OK response")
  @get ok(): OK;

  @route("/health")
  @doc("Check the health of the service, always returns 200 OK")
  @get health(): PlainText;

  @route("/info")
  @doc("Get system information")
  @get info(): OK;
}

@tag("Inspect Routes")
@route("/inspect")
interface Inspect {
  @doc("Inspect the incoming HTTP request and return the details")
  @get inspect(): RequestInfo;
  
  @doc("Inspect the incoming HTTP request and return the details")
  @post inspectPost(): RequestInfo;
  
  @doc("Inspect the incoming HTTP request and return the details")
  @put inspectPut(): RequestInfo;
  
  @doc("Inspect the incoming HTTP request and return the details")
  @delete inspectDelete(): RequestInfo;
  
  @doc("Inspect the incoming HTTP request and return the details")
  @patch inspectPatch(): RequestInfo;
}

@tag("Wildcard Inspection Routes")
@route("/anything")
interface Wildcard {
  @doc("Any route or path that doesn't match other operations is automatically inspected, e.g. /blah/cheese/shoe will still return the details of the request")
  @get inspectAnything(@path extraPath: string = "foo"): RequestInfo;

  @doc("Any route or path that doesn't match other operations is automatically inspected, e.g. /blah/cheese/shoe will still return the details of the request")
  @post inspectAnythingPost(@path extraPath: string = "foo"): RequestInfo;

  @doc("Any route or path that doesn't match other operations is automatically inspected, e.g. /blah/cheese/shoe will still return the details of the request")
  @put inspectAnythingPut(@path extraPath: string = "foo"): RequestInfo;

  @doc("Any route or path that doesn't match other operations is automatically inspected, e.g. /blah/cheese/shoe will still return the details of the request")
  @delete inspectAnythingDelete(@path extraPath: string = "foo"): RequestInfo;

  @doc("Any route or path that doesn't match other operations is automatically inspected, e.g. /blah/cheese/shoe will still return the details of the request")
  @patch inspectAnythingPatch(@path extraPath: string = "foo"): RequestInfo;
}

@tag("Utility Routes")
interface Utils {
  @route("/status/{code}")
  @doc("Get a response with the specified status code")
  @get status(code: integer): PlainText;

  @route("/word")
  @doc("Get a random word")
  @get word(): PlainText;

  @route("/word/{count}")
  @doc("Get a number of random words")
  @get words(count: integer): PlainText;

  @route("/number")
  @doc("Get a random number")
  @get number(): PlainText;

  @route("/number/{max}")
  @doc("Get a random number up to the specified maximum")
  @get numberMax(max: integer): PlainText;

  @route("/uuid")
  @doc("Get a random UUID")
  @get uuid(): PlainText;

  @route("/uuid/{input}")
  @doc("Get a UUID generated from the input string")
  @get uuidFrom(input: string): PlainText;

  @route("/delay")
  @doc("Delay the response by 1 to 10 seconds randomly")
  @get delayRandom(): PlainText;

  @route("/delay/{seconds}")
  @doc("Delay the response by the specified number of seconds")
  @get delay(seconds: integer): PlainText;
}

@tag("Authenticated Routes")
@route("/auth")
interface Auth {
  @route("/basic")
  @doc("Route is protected by basic auth, see docs for the credentials")
  @useAuth(BasicAuth)
  @get basicAuth(): OK | ForbiddenResponse | UnauthorizedResponse;

  @route("/basic")
  @doc("Route is protected by basic auth, see docs for the credentials")
  @useAuth(BasicAuth)
  @post basicAuthPost(): OK | ForbiddenResponse | UnauthorizedResponse;

  @route("/jwt")
  @doc("Route is protected by SHA256 JWT auth, see docs for obtaining the token")
  @useAuth(BearerAuth)
  @get jwtAuth(): OK | ForbiddenResponse | UnauthorizedResponse;

  @route("/jwt")
  @doc("Route is protected by basic auth, see docs for the credentials")
  @useAuth(BasicAuth)
  @post jwtAuthPost(): OK | ForbiddenResponse | UnauthorizedResponse;
}
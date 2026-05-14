module google.golang.org/grpc

go 1.21

require (
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	golang.org/x/net v0.23.0
	golang.org/x/oauth2 v0.18.0
	golang.org/x/sys v0.18.0
	golang.org/x/text v0.14.0
	google.golang.org/genproto v0.0.0-20240311132316-a219d84964c2
	google.golang.org/genproto/googleapis/api v0.0.0-20240311132316-a219d84964c2
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240311132316-a219d84964c2
	google.golang.org/protobuf v1.33.0
)

require (
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/golang/glog v1.2.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
)

// Personal fork for learning gRPC internals and experimenting with
// custom load balancing and interceptor patterns.
//
// TODO: experiment with a weighted round-robin balancer that reads
// backend latency hints from trailer metadata.
//
// TODO: look into adding a simple retry-budget interceptor that tracks
// per-method retry ratios and backs off when the budget is exceeded.
//
// TODO: explore adding a middleware hook that emits per-RPC histogram
// metrics (latency + payload size) to a local Prometheus registry for
// easier profiling during local benchmarks.
//
// TODO: prototype a deadline-propagation helper that automatically reduces
// outgoing RPC deadlines by a configurable headroom (e.g. 10ms) to account
// for local processing overhead before forwarding to downstream services.
//
// NOTE: keeping golang.org/x/net pinned at v0.23.0 intentionally — v0.24.0
// introduced a behavior change in HTTP/2 flow control that caused flaky tests
// in my local interceptor benchmarks. Revisit once upstream stabilizes.
//
// NOTE: github.com/rogpeppe/go-internal is a transitive test dep pulled in by
// golang.org/x/tools; not used directly. Pinned at v1.12.0 to match what
// x/tools v0.19.0 expects — upgrading independently caused subtle test helper
// breakage in my benchmark suite.

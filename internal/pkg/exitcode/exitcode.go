package exitcode

const (
	Success             int = iota // Normal termination
	GenericError                   // Unclassified error, fallback
	ConfigError                    // Missing/invalid env vars, bad config file
	DatabaseError                  // Cannot connect to or migrate database
	MCPTransportError              // Failed to start HTTP/SSE server
	AuthenticationError            // JWT signing/validation failure at startup
	DependencyError                // External services (e.g., metrics backend)
	FatalBug                       // Panic or unexpected state
)

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hungpdn/mcp2grule/internal/api/handler"
	"github.com/hungpdn/mcp2grule/internal/config"
	"github.com/hungpdn/mcp2grule/internal/pkg/logger"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server represents the MCP server with its handler and underlying mcp.Server instance.
type Server struct {
	mcpHandler *handler.MCPHandler
	server     *mcp.Server
}

// NewServer creates a new MCP server instance with the given application name, version, and handler.
func NewServer(appName, verison string, mcpHandler *handler.MCPHandler) *Server {

	mcpServer := mcp.NewServer(&mcp.Implementation{Name: appName, Version: verison}, nil)

	srv := &Server{mcpHandler: mcpHandler, server: mcpServer}
	return srv
}

// Run starts the MCP server based on the configured transport method (stdio or HTTP).
func (s *Server) Run(ctx context.Context) error {
	// Register tools
	s.AddTools()

	// Handle shutdown signals
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	go func() {
		<-sigChan
		logger.Infof("Received shutdown signal")
		cancel()

		// For stdio, close stdin to unblock the Listen call
		if config.App.MCPTransport == config.MCPTransportStdio {
			_ = os.Stdin.Close()
		}
	}()

	// Start the server based on the configured transport
	switch config.App.MCPTransport {
	case config.MCPTransportStdio:
		if err := s.runStdio(ctx); err != nil {
			return err
		}
	case config.MCPTransportSSE:
		httpHandler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server { return s.server })
		// TODO: enable auth
		// httpHandler = middleware.Auth(httpHandler, config.App.StreamableHTTPTransport.AuthToken)
		srv := &http.Server{
			Addr:    config.App.HTTPTransport.HttpAddr(),
			Handler: httpHandler,
		}
		if err := s.runHTTPServer(ctx, srv, config.App.MCPTransport); err != nil {
			return err
		}
	case config.MCPTransportStreamableHTTP:
		httpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return s.server }, nil)
		srv := &http.Server{
			Addr:    config.App.HTTPTransport.HttpAddr(),
			Handler: httpHandler,
		}
		if err := s.runHTTPServer(ctx, srv, config.App.MCPTransport); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown MCP_TRANSPORT: %s", config.App.MCPTransport)
	}

	return nil
}

// runHTTPServer starts the MCP server using HTTP transport.
func (s *Server) runHTTPServer(ctx context.Context, srv *http.Server, transport config.MCPTransport) error {

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		logger.Infof("%s server shutting down...", transport)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("%s shutdown error: %v", transport, err)
		}

		select {
		case err := <-serverErr:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				return fmt.Errorf("%s erver error during shutdown: %v", transport, err)
			}
		case <-shutdownCtx.Done():
			logger.Warnf("%s server did not stop gracefully within timeout", transport)
		}
	}

	return nil
}

// runStdio starts the MCP server using standard input/output for communication.
func (s *Server) runStdio(ctx context.Context) error {
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

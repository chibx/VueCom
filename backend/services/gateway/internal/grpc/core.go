package grpc

import (
	"context"
	"net"
	"sync"

	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type manager struct {
	mu       sync.Mutex
	listener *bufconn.Listener
	server   *grpc.Server
	once     sync.Once // ensure started only once
}

var _conn *grpc.ClientConn
var global = &manager{}

// registerServices — call this at startup with all your RegisterXXXServer funcs
func registerServices(registerFns ...func(*grpc.Server)) {
	global.once.Do(func() {
		global.mu.Lock()
		global.listener = bufconn.Listen(bufSize)
		global.server = grpc.NewServer(
		// Add global interceptors here if needed (logging, tracing, auth)
		// grpc.ChainUnaryInterceptor(yourInterceptors...)
		)

		for _, fn := range registerFns {
			fn(global.server)
		}

		go func() {
			if err := global.server.Serve(global.listener); err != nil && err != grpc.ErrServerStopped {
				utils.Logger().Fatal("", zap.Error(err)) // or use a proper logger.Fatal
			}
		}()

		global.mu.Unlock()
	})
}

// clientConn — get a connection to the in-memory server (call from any module)
func clientConn() *grpc.ClientConn {
	var err error
	global.mu.Lock()
	if global.listener == nil {
		global.mu.Unlock()
		utils.Logger().Fatal("inproc not started — call RegisterServices first")
	}
	lis := global.listener // safe copy under lock
	global.mu.Unlock()

	if _conn == nil {
		_conn, err = grpc.NewClient(
			"passthrough://bufconn",
			grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
				return lis.DialContext(ctx)
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			utils.Logger().Fatal("", zap.Error(err)) // handle gracefully in real code
		}
	}

	return _conn
}

// shutdown — call on app exit
func shutdown() {
	global.mu.Lock()
	defer global.mu.Unlock()
	if global.server != nil {
		global.server.GracefulStop()
	}
	if global.listener != nil {
		global.listener.Close()
	}
}

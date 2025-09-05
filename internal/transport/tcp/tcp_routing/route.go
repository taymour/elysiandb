package tcprouting

import (
	"net"

	"github.com/taymour/elysiandb/internal/log"
	"github.com/taymour/elysiandb/internal/transport/tcp/handler"
	"github.com/taymour/elysiandb/internal/transport/tcp/parser"
)

func RouteLine(line []byte, c net.Conn) []byte {
	cmd, query := parser.FirstWordBytes(line)

	switch {
	case parser.EqASCII(cmd, []byte("EXIT")):
		_ = c.Close()
		return []byte("Goodbye!")

	case parser.EqASCII(cmd, []byte("GET")):
		return handler.HandleGet(query)

	case parser.EqASCII(cmd, []byte("SET")):
		return handler.HandleSet(query)
	}

	log.Error("Unknown command:", string(cmd))

	return []byte("ERR unknown command")
}

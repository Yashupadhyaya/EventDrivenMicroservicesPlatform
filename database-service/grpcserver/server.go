package grpcserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Yashupadhyaya/database-service/database"
	"github.com/Yashupadhyaya/database-service/proto"
)

type Server struct {
	proto.UnimplementedDatabaseServiceServer
	DB *sql.DB
}

func (s *Server) Query(ctx context.Context, req *proto.QueryRequest) (*proto.QueryResponse, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	// Convert []string to []interface{} for database query
	params := make([]interface{}, len(req.Params))
	for i, param := range req.Params {
		params[i] = param
	}

	rows, err := database.Query(s.DB, req.Sql, params...)
	if err != nil {
		return nil, err
	}

	response := &proto.QueryResponse{}
	for _, row := range rows {
		rowMap := make(map[string]string)
		for key, value := range row {
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("expected string value, got %T", value)
			}
			rowMap[key] = strValue
		}
		// Wrap rowMap in a proto.Row and append to response.Rows
		response.Rows = append(response.Rows, &proto.Row{Columns: rowMap})
	}
	return response, nil
}

func (s *Server) Command(ctx context.Context, req *proto.CommandRequest) (*proto.CommandResponse, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	// Convert []string to []interface{} for database command
	params := make([]interface{}, len(req.Params))
	for i, param := range req.Params {
		params[i] = param
	}

	rowsAffected, err := database.Execute(s.DB, req.Sql, params...)
	if err != nil {
		return nil, err
	}
	return &proto.CommandResponse{RowsAffected: rowsAffected}, nil
}

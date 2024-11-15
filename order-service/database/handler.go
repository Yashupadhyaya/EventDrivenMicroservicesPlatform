package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Yashupadhyaya/order-service/models"
	proto "github.com/Yashupadhyaya/order-service/proto"
	"google.golang.org/grpc"
)

var client proto.DatabaseServiceClient

func InitDatabase(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client = proto.NewDatabaseServiceClient(conn)
	return nil
}

// func convertParamsToStrings(params []interface{}) []string {
// 	var strParams []string
// 	for _, param := range params {
// 		switch v := param.(type) {
// 		case string:
// 			strParams = append(strParams, v)
// 		case int:
// 			strParams = append(strParams, strconv.Itoa(v))
// 		case int64:
// 			strParams = append(strParams, strconv.FormatInt(v, 10))
// 		case float64:
// 			strParams = append(strParams, strconv.FormatFloat(v, 'f', -1, 64))
// 		default:
// 			strParams = append(strParams, fmt.Sprint(v))
// 		}
// 	}
// 	return strParams
// }

func CreateOrder(order *models.Order) error {
	// Prepare the SQL command and parameters
	fmt.Println("CreateOrder")
	sql := "INSERT INTO orders (id, customer_id, status) VALUES ($1, $2, $3)"
	params := []interface{}{order.ID, order.CustomerID, order.Status}

	// Use the gRPC client to call the Command method on the Database Service
	_, err := client.Command(context.Background(), &proto.CommandRequest{
		Sql:    sql,
		Params: convertParamsToStrings(params),
	})
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	return nil
}

func convertParamsToStrings(params []interface{}) []string {
	strParams := make([]string, len(params))
	for i, param := range params {
		switch v := param.(type) {
		case int, int64, float32, float64:
			strParams[i] = fmt.Sprintf("%v", v)
		case string:
			strParams[i] = v
		default:
			strParams[i] = fmt.Sprintf("%v", param)
		}
	}
	return strParams
}

func GetOrderById(orderId string) (*models.Order, error) {
	// Prepare the SQL query and parameters
	sql := "SELECT id, customer_id, status FROM orders WHERE id = $1"
	params := []interface{}{orderId}

	// Use the gRPC client to call the Query method on the Database Service
	resp, err := client.Query(context.Background(), &proto.QueryRequest{
		Sql:    sql,
		Params: convertParamsToStrings(params),
	})
	if err != nil {
		return nil, err
	}

	// Convert the response to your application's Order model
	return convertToOrder(resp.Rows), nil
}

// Helper functions to convert between gRPC and application models

func convertToOrder(rows []*proto.Row) *models.Order {
	if len(rows) == 0 {
		return nil
	}

	firstRow := rows[0]
	order := &models.Order{
		ID:         firstRow.Columns["id"],
		CustomerID: firstRow.Columns["customer_id"],
		Status:     firstRow.Columns["status"],
		Items:      make([]models.OrderItem, 0),
	}

	for _, row := range rows {
		item := models.OrderItem{
			ProductID: row.Columns["product_id"],
			Quantity:  mustParseInt(row.Columns["quantity"]),
		}
		order.Items = append(order.Items, item)
	}

	return order
}

func mustParseInt(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		// Handle the error according to your application's logic
		// For now, panic to demonstrate an error scenario
		panic(fmt.Sprintf("cannot parse integer: %v", err))
	}
	return intValue
}

// func mustParseInt(s string) int {
// 	i, err := strconv.Atoi(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return i
// }

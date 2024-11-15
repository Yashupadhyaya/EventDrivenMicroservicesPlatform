
package models

type QueryRequest struct {
    SQL    string `json:"sql"`
    Params []string `json:"params"`
}

type QueryResponse struct {
    Rows []map[string]string `json:"rows"`
}

type CommandRequest struct {
    SQL    string `json:"sql"`
    Params []string `json:"params"`
}

type CommandResponse struct {
    RowsAffected int64 `json:"rows_affected"`
}

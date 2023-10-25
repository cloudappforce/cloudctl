package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/cloudappforce/cloudctl/pkg/api/models"
	"github.com/cloudappforce/cloudctl/pkg/templates"
	"github.com/go-resty/resty/v2"
)

func FormatDatabaseListResponse(resp *resty.Response) error {
	body := resp.Body()
	if resp.StatusCode() != 200 {
		return fmt.Errorf(resp.Status())
	}
	var data []models.ListDatabasesResponse
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	columns := []string{"NAME", "TYPE", "INSTANCES", "SIZE"}
	var rows [][]string
	for _, v := range data {
		rows = append(rows, []string{v.Name, "postgres", fmt.Sprint(v.Instances), v.Storage.Size})
	}
	return templates.NewTableFormatter().Format(columns, rows)
}

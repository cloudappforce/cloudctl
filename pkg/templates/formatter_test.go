package templates_test

import (
	"testing"

	"github.com/cloudappforce/cloudctl/pkg/templates"
)

func TestFormatter(t *testing.T) {
	testCases := [][]string{
		{"db-shopping-prod", "postgres", "3", "50Gi", "db-shopping-prod-r", "db-shopping-prod-ro", "db-shopping-prod-rw"},
		{"db-shopping-prod", "postgres", "3", "50Gi", "db-shopping-prod-r", "db-shopping-prod-ro", "db-shopping-prod-rw"},
		{"db-shopping-prod", "postgres", "3", "50Gi", "db-shopping-prod-r", "db-shopping-prod-ro", "db-shopping-prod-rw"},
		{"db-shopping-prod", "postgres", "3", "50Gi", "db-shopping-prod-r", "db-shopping-prod-ro", "db-shopping-prod-rw"},
		{"db-shopping-prod", "postgres", "3", "50Gi", "db-shopping-prod-r", "db-shopping-prod-ro", "db-shopping-prod-rw"},
		{"db-shopping-prod", "postgres", "3", "50Gi", "db-shopping-prod-r", "db-shopping-prod-ro", "db-shopping-prod-rw"},
		{"db-shopping-prod", "postgres", "3", "50Gi", "db-shopping-prod-r", "db-shopping-prod-ro", "db-shopping-prod-rw"},
		{"db-shopping-prod", "postgres", "3", "50Gi", "db-shopping-prod-r", "db-shopping-prod-ro", "db-shopping-prod-rw"},
	}
	tableFormatter := templates.NewTableFormatter()

	columns := []string{"NAME", "TYPE", "INSTANCES", "SIZE", "DATABASE-R", "DATABASE-RO", "DATABASE-RW"}

	err := tableFormatter.Format(columns, testCases)
	if err != nil {
		t.Fail()
	}

}

package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudappforce/cloudctl/pkg/api/models"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/pflag"
)

type Client struct {
	JWT    string
	Host   string
	Scheme string
}

type CreateDatabaseInput struct {
	Size      string            `json:"size"`
	Instances string            `json:"instances"`
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels,omitempty"`
}

type CreateIngressInput struct {
	Name string `json:"name"`
}

type DeleteIngressInput struct {
	Name string `json:"name"`
}

type ListDatabaseInput struct {
	Labels map[string]string `json:"labels,omitempty"`
}
type DeleteDatabaseInput struct {
	Name string
}

const (
	ArgSize      = "size"
	ArgInstances = "instances"
	ArgType      = "type"
	ArgName      = "name"
	ArgLabels    = "labels"
)

func NewCreateDatabaseInputFromFlags(flags *pflag.FlagSet) (*CreateDatabaseInput, error) {
	v := &CreateDatabaseInput{}
	err := v.parse(flags)
	if err != nil {
		return nil, err
	}
	return v, nil
}
func NewDeleteInputFromFlags(flags *pflag.FlagSet) (*DeleteDatabaseInput, error) {
	v := &DeleteDatabaseInput{}
	err := v.parse(flags)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func NewListDatabaseInputFromFlags(flags *pflag.FlagSet) (*ListDatabaseInput, error) {
	v := &ListDatabaseInput{}
	err := v.parse(flags)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (v *DeleteDatabaseInput) parse(flags *pflag.FlagSet) error {
	name, err := flags.GetString(ArgName)
	if err != nil {
		return err
	}
	v.Name = name
	return nil
}

func (v *CreateDatabaseInput) parse(flags *pflag.FlagSet) error {
	size, err := flags.GetString(ArgSize)
	if err != nil {
		return err
	}
	instances, err := flags.GetString(ArgInstances)
	if err != nil {
		return err
	}
	t, err := flags.GetString(ArgType)
	if err != nil {
		return err
	}
	name, err := flags.GetString(ArgName)
	if err != nil {
		return err
	}
	labels, err := flags.GetStringToString(ArgLabels)
	if err != nil {
		return err
	}
	v.Name = name
	v.Size = size
	v.Type = t
	v.Instances = instances
	v.Labels = labels
	return nil
}

func (v *ListDatabaseInput) parse(flags *pflag.FlagSet) error {
	// TODO: support namespaces and labels
	// currently there are no flags for listing databases
	return nil
}

func (c *Client) hostWithScheme() string {
	return fmt.Sprintf("%s://%s", c.Scheme, c.Host)
}

// Create a new database cluster
func (c *Client) CreateDatabase(ctx context.Context, input *CreateDatabaseInput) (*models.CreateDatabaseResponse, error) {
	restyClient := resty.New()
	restyClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.JWT))
	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(input).
		Post(fmt.Sprintf("%s/api/databases", c.hostWithScheme()))
	if err != nil {
		// panic(err)
		return nil, err
	}
	m := &models.CreateDatabaseResponse{}
	err = json.Unmarshal(resp.Body(), m)
	return m, err
}

// Create a new network ingress (exposes the service to the public)
func (c *Client) CreateIngress(ctx context.Context, input *CreateIngressInput) (*models.CreateIngressResponse, error) {
	restyClient := resty.New()
	restyClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.JWT))
	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(input).
		Post(fmt.Sprintf("%s/api/ingress", c.hostWithScheme()))
	if err != nil {
		return nil, err
	}
	m := &models.CreateIngressResponse{}
	err = json.Unmarshal(resp.Body(), m)
	return m, err
}

// Create a new network ingress (exposes the service to the public)
func (c *Client) DeleteIngress(ctx context.Context, input *DeleteIngressInput) (*models.DeleteResponse, error) {
	restyClient := resty.New()
	restyClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.JWT))
	resp, err := restyClient.R().Delete(fmt.Sprintf("%s/api/ingress/%s", c.hostWithScheme(), input.Name))
	if err != nil {
		return nil, err
	}
	m := &models.DeleteResponse{}
	err = json.Unmarshal(resp.Body(), m)
	return m, err
}

// Create a new deployment
func (c *Client) CreateDeployment(ctx context.Context, input *models.CreateDeploymentInput) (*models.CreateDeploymentResponse, error) {
	restyClient := resty.New()
	restyClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.JWT))
	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(input).
		Post(fmt.Sprintf("%s/api/deployments", c.hostWithScheme()))
	if err != nil {
		return nil, err
	}
	m := &models.CreateDeploymentResponse{}
	err = json.Unmarshal(resp.Body(), m)
	return m, err
}

// Create a new database cluster
func (c *Client) DeleteDatabase(ctx context.Context, input *DeleteDatabaseInput) error {

	restyClient := resty.New()
	restyClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.JWT))
	_, err := restyClient.R().Delete(fmt.Sprintf("%s/api/databases/%s", c.hostWithScheme(), input.Name))
	if err != nil {
		return err
	}

	return nil
}

// List database clusters available
func (c *Client) ListDatabases(ctx context.Context, input *ListDatabaseInput) (*resty.Response, error) {

	restyClient := resty.New()
	restyClient.SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.JWT))
	return restyClient.R().Get(fmt.Sprintf("%s/api/databases", c.hostWithScheme()))

}

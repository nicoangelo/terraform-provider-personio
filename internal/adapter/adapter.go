package adapter

import (
	"context"

	personio "github.com/giantswarm/personio-go/v1"
)

const (
	ApiBaseUrlDefault string = personio.DefaultBaseUrl
)

type PersonioAdapter struct {
	Client *personio.Client
}

func NewAdapter(apiBaseUrl string, clientId string, clientSecret string) (*PersonioAdapter, error) {
	credentials := personio.Credentials{ClientId: clientId, ClientSecret: clientSecret}
	client, err := personio.NewClient(context.TODO(), apiBaseUrl, credentials)

	if err == nil {
		return &PersonioAdapter{
			Client: client,
		}, nil
	}
	return nil, err
}

func (p *PersonioAdapter) GetEmployees() (employees []Employee, err error) {
	pe, err := p.Client.GetEmployees()
	if err != nil {
		return employees, err
	}
	for _, v := range pe {
		employees = append(employees, NewEmployee(v))
	}
	return employees, nil
}

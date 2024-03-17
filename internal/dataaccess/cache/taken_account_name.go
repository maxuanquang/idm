package cache

import (
	"context"
)

var (
	setKeyNameTakenAccountName string = "taken_account_name_set"
)

type TakenAccountName interface {
	Add(ctx context.Context, accountName string) error
	Has(ctx context.Context, accountName string) (bool, error)
}

func NewTakenAccountName(client Client) (TakenAccountName, error) {
	return &takenAccountName{
		client: client,
	}, nil
}

type takenAccountName struct {
	client Client
}

// Add implements TakenAccountName.
func (t *takenAccountName) Add(ctx context.Context, accountName string) error {
	err := t.client.AddToSet(ctx, setKeyNameTakenAccountName, accountName)
	if err != nil {
		return err
	}

	return nil
}

// Has implements TakenAccountName.
func (t *takenAccountName) Has(ctx context.Context, accountName string) (bool, error) {
	exists, err := t.client.IsValueInSet(ctx, setKeyNameTakenAccountName, accountName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

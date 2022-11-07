package kernel_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kasperbrandtpedersen/kernel"
)

type RegisterItem struct {
	User string
	SKU  string
}

func (cmd *RegisterItem) Stream() string {
	return fmt.Sprintf("inventory_item_%v", cmd.SKU)
}

type ItemRegistered struct {
	SKU          string
	EventVersion int
	EventAt      time.Time
	EventBy      string
}

type InventoryItem struct {
	SKU        string
	Created    time.Time
	CreatedBy  string
	Modified   time.Time
	ModifiedBy string
	Version    int
	Quantity   int
}

func (i *ItemRegistered) Version() int {
	return i.EventVersion
}

func (i *ItemRegistered) At() time.Time {
	return i.EventAt
}

func (i *ItemRegistered) By() string {
	return i.EventBy
}

func (i *InventoryItem) On(e kernel.Event) bool {
	switch v := e.(type) {
	case *ItemRegistered:
		i.SKU = v.SKU
		i.Created = v.EventAt
		i.CreatedBy = v.EventBy

	default:
		return false
	}

	i.Modified = e.At()
	i.ModifiedBy = e.By()
	i.Version = e.Version()

	return true
}

func ItemRegister(state *InventoryItem, cmd *RegisterItem) []kernel.Event {
	return []kernel.Event{
		&ItemRegistered{
			SKU:          cmd.SKU,
			EventVersion: 0,
			EventAt:      time.Now(),
			EventBy:      "Foo Bar",
		},
	}

}

func initialState() *InventoryItem {
	return &InventoryItem{}
}

func TestKernel(t *testing.T) {
	k := kernel.New(
		kernel.Emits(&ItemRegistered{}),
		kernel.Decide(ItemRegister, initialState),
	)

	err := k.Dispatch(&RegisterItem{
		SKU:  "123456",
		User: "Foo Bar",
	})

	if err != nil {
		t.Fatal(err)
	}
}

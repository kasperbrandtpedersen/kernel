package kernel_test

import (
	"fmt"
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

type RestockItem struct {
	User     string
	SKU      string
	Quantity int
}

func (cmd *RestockItem) Stream() string {
	return fmt.Sprintf("inventory_item_%v", cmd.SKU)
}

type ShipItem struct {
	User     string
	SKU      string
	Quantity int
}

func (cmd *ShipItem) Stream() string {
	return fmt.Sprintf("inventory_item_%v", cmd.SKU)
}

type ItemRegistered struct {
	kernel.EventModel

	SKU string
}

type ItemRestocked struct {
	kernel.EventModel

	SKU          string
	Incoming     int
	CurrentStock int
}

type ItemShipped struct {
	kernel.EventModel

	SKU          string
	OutGoing     int
	CurrentStock int
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

func (i *InventoryItem) On(e kernel.Event) bool {
	switch v := e.(type) {
	case *ItemRegistered:
		i.SKU = v.SKU
		i.Created = v.EventAt
		i.CreatedBy = v.EventBy

	case *ItemRestocked:
		i.Quantity += v.Incoming

	case *ItemShipped:
		i.Quantity -= v.OutGoing

	default:
		return false
	}

	i.Modified = e.At()
	i.ModifiedBy = e.By()
	i.Version = e.Version()

	return true
}

func ItemRegiser(state *InventoryItem, cmd *RegisterItem) []kernel.Event {
	return []kernel.Event{
		&ItemRegistered{
			SKU: cmd.SKU,
			EventModel: kernel.EventModel{
				EventVersion: state.Version + 1,
				EventAt:      time.Now(),
				EventBy:      cmd.User,
			},
		},
	}

}

func Example() {
	k := kernel.New(
		kernel.Decide(ItemRegiser),
	)

	k.Dispatch(&RegisterItem{
		SKU:  "123456",
		User: "Foo Bar",
	})
}

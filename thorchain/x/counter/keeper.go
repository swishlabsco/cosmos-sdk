package counter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	costGet      sdk.Gas = 10
	costAdd      sdk.Gas = 10
	costSubtract sdk.Gas = 10
)

// Keeper of the counter store
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a counter keeper
func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
	return keeper
}

// Get returns the counter.
func (keeper Keeper) Get(ctx sdk.Context) int {
	return get(ctx, keeper.storeKey)
}

// Add adds amount to the counter.
func (keeper Keeper) Add(ctx sdk.Context, amount int) sdk.Error {
	return add(ctx, keeper.storeKey, amount)
}

// Add adds amount to the counter.
func (keeper Keeper) Subtract(ctx sdk.Context, amount int) sdk.Error {
	return subtract(ctx, keeper.storeKey, amount)
}

func get(ctx sdk.Context, key sdk.StoreKey) int {
	ctx.GasMeter().ConsumeGas(costGet, "get")

	//get
	amount := 4

	return amount
}

func add(ctx sdk.Context, key sdk.StoreKey, amount int) sdk.Error {
	ctx.GasMeter().ConsumeGas(costAdd, "add")
	//get
	//currentCount := get(ctx, key)

	//add
	//newCount := currentCount + amount

	//checks?
	//error

	//set

	//checks?
	//error

	return nil
}

func subtract(ctx sdk.Context, key sdk.StoreKey, amount int) sdk.Error {
	ctx.GasMeter().ConsumeGas(costSubtract, "subtract")
	//get
	//currentCount := get(ctx, key)

	//subtract
	//newCount := currentCount - amount

	//checks?
	//error

	//set

	//checks?
	//error

	return nil
}

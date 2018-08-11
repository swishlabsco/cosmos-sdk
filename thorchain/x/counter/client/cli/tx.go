package cli

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/thorchain/x/counter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagAmount = "amount"
)

// GetCmdAdd will create an add tx and sign it with the given key
func GetCmdAdd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Create and sign an add tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			// get the from/to address
			from, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}

			fromAcc, err := ctx.QueryStore(auth.AddressStoreKey(from), ctx.AccountStore)
			if err != nil {
				return err
			}

			// Check if account was found
			if fromAcc == nil {
				return errors.Errorf("No account with address %s was found in the state.\nAre you sure there has been a transaction involving it?", from)
			}

			// parse amount trying to be added
			amountStr := viper.GetString(flagAmount)
			amount, err := strconv.ParseInt(amountStr, 10, 64)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := counter.NewMsgAdd(fromAcc, int(amount))
			err = ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, []sdk.Msg{msg}, cdc)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().String(flagAmount, "", "Amount to add")

	return cmd
}

// GetCmdSubtract will create an add tx and sign it with the given key
func GetCmdSubtract(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subtract",
		Short: "Create and sign a subtract tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))
			// get the from/to address
			from, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}

			fromAcc, err := ctx.QueryStore(auth.AddressStoreKey(from), ctx.AccountStore)
			if err != nil {
				return err
			}

			// Check if account was found
			if fromAcc == nil {
				return errors.Errorf("No account with address %s was found in the state.\nAre you sure there has been a transaction involving it?", from)
			}

			// parse amount trying to be added
			amountStr := viper.GetString(flagAmount)
			amount, err := strconv.ParseInt(amountStr, 10, 64)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := counter.NewMsgSubtract(fromAcc, int(amount))
			err = ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, []sdk.Msg{msg}, cdc)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().String(flagAmount, "", "Amount to subtract")
	cmd.MarkFlagRequired("amount")

	return cmd
}

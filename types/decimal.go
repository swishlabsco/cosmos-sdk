package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

// NOTE: never use new(Dec) or else we will panic unmarshalling into the
// nil embedded big.Int
type Dec struct {
	*big.Int `json:"int"`
}

// number of decimal places
const Precision = 10

func precisionInt() *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(Precision), nil)
}

// nolint - common values
func ZeroDec() Dec { return Dec{big.NewInt(0)} }
func OneDec() Dec  { return Dec{precisionInt()} }

// get the precision multiplier
func precisionMultiplier(prec int64) *big.Int {
	if prec > Precision {
		panic("too much precision")
	}
	zerosToAdd := Precision - prec
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(zerosToAdd), nil)
	return multiplier
}

// create a new Dec from integer assuming whole numbers
// CONTRACT: prec <= Precision
func NewDec(i, prec int64) Dec {
	return Dec{
		new(big.Int).Mul(big.NewInt(i), precisionMultiplier(prec)),
	}
}

// create a new Dec from big integer assuming whole numbers
// CONTRACT: prec !> Precision
func NewDecFromBigInt(i *big.Int, prec int64) Dec {
	return Dec{
		new(big.Int).Mul(i, precisionMultiplier(prec)),
	}
}

// create a new Dec from big integer assuming whole numbers
// CONTRACT: prec !> Precision
func NewDecFromInt(i Int, prec int64) Dec {
	return Dec{
		new(big.Int).Mul(i.BigInt(), precisionMultiplier(prec)),
	}
}

// create a decimal from a decimal string (ex. "1234.5678")
func NewDecFromStr(str string) (d Dec, err Error) {
	if len(str) == 0 {
		return d, ErrUnknownRequest("decimal string is empty")
	}

	// first extract any negative symbol
	neg := false
	if string(str[0]) == "-" {
		neg = true
		str = str[1:]
	}

	if len(str) == 0 {
		return d, ErrUnknownRequest("decimal string is empty")
	}

	strs := strings.Split(str, ".")
	lenDecs := 0
	combinedStr := strs[0]
	if len(strs) == 2 {
		lenDecs = len(strs[1])
		if lenDecs == 0 || len(combinedStr) == 0 {
			return d, ErrUnknownRequest("bad decimal length")
		}
		combinedStr = combinedStr + strs[1]
	} else if len(strs) > 2 {
		return d, ErrUnknownRequest("too many periods to be a decimal string")
	}

	if lenDecs > Precision {
		return d, ErrUnknownRequest("too much Precision in decimal")
	}

	// add some extra zero's to correct to the Precision factor
	zerosToAdd := Precision - lenDecs
	zeros := fmt.Sprintf(`%0`+strconv.Itoa(zerosToAdd)+`s`, "")
	combinedStr = combinedStr + zeros

	combined, ok := new(big.Int).SetString(combinedStr, 10)
	if !ok {
		return d, ErrUnknownRequest("bad string to integer conversion")
	}
	if neg {
		combined = new(big.Int).Neg(combined)
	}
	return Dec{combined}, nil
}

//nolint
func (d Dec) IsZero() bool      { return (d.Int).Sign() == 0 } // Is equal to zero
func (d Dec) Equal(d2 Dec) bool { return (d.Int).Cmp(d2.Int) == 0 }
func (d Dec) GT(d2 Dec) bool    { return (d.Int).Cmp(d2.Int) == 1 }             // greater than
func (d Dec) GTE(d2 Dec) bool   { return (d.Int).Cmp(d2.Int) >= 0 }             // greater than or equal
func (d Dec) LT(d2 Dec) bool    { return (d.Int).Cmp(d2.Int) == -1 }            // less than
func (d Dec) LTE(d2 Dec) bool   { return (d.Int).Cmp(d2.Int) <= 0 }             // less than or equal
func (d Dec) Neg() Dec          { return Dec{new(big.Int).Neg(d.Int)} }         // Is equal to zero
func (d Dec) Add(d2 Dec) Dec    { return Dec{new(big.Int).Add(d.Int, d2.Int)} } // addition
func (d Dec) Sub(d2 Dec) Dec    { return Dec{new(big.Int).Sub(d.Int, d2.Int)} } // subtraction

// multiplication
func (d Dec) Mul(d2 Dec) Dec {
	mul := new(big.Int).Mul(d.Int, d2.Int)
	chopped := BankerRoundChop(mul, Precision)
	return Dec{chopped}
}

// quotient
func (d Dec) Quo(d2 Dec) Dec {
	mul := new(big.Int).Mul(new(big.Int).Mul( // multiple Precision twice
		d.Int, precisionInt()), precisionInt())

	quo := new(big.Int).Quo(mul, d2.Int)
	chopped := BankerRoundChop(quo, Precision)
	return Dec{chopped}
}

func (d Dec) String() string {
	str := d.ToLeftPaddedWithDecimals(Precision)
	placement := len(str) - Precision
	if placement < 0 {
		panic("too few decimal digits")
	}
	return str[:placement] + "." + str[placement:]
}

// TODO panic if negative or if totalDigits < len(initStr)???
// evaluate as an integer and return left padded string
func (d Dec) ToLeftPaddedWithDecimals(totalDigits int8) string {
	intStr := d.Int.String()
	fcode := `%0` + strconv.Itoa(int(totalDigits)) + `s`
	return fmt.Sprintf(fcode, intStr)
}

// TODO panic if negative or if totalDigits < len(initStr)???
// evaluate as an integer and return left padded string
func (d Dec) ToLeftPadded(totalDigits int8) string {
	chopped := BankerRoundChop(d.Int, Precision)
	intStr := chopped.String()
	fcode := `%0` + strconv.Itoa(int(totalDigits)) + `s`
	return fmt.Sprintf(fcode, intStr)
}

//     ____
//  __|    |__   "chop 'em
//       ` \     round!"
// ___||  ~  _     -bankers
// |         |      __
// |       | |   __|__|__
// |_____:  /   | $$$    |
//              |________|

// nolint - go-cyclo
// chop of n digits, and banker round the digits being chopped off
// Examples:
//   BankerRoundChop(1005, 1) = 100
//   BankerRoundChop(1015, 1) = 102
//   BankerRoundChop(1500, 3) = 2
func BankerRoundChop(d *big.Int, n int64) (chopped *big.Int) {

	// remove the negative and add it back when returning
	if d.Sign() == -1 {
		d = new(big.Int).Neg(d)
		defer func() {
			chopped = new(big.Int).Neg(chopped)
		}()
	}

	// get the trucated quotient and remainder
	quo, rem, prec := big.NewInt(0), big.NewInt(0), precisionInt()
	quo, rem = quo.QuoRem(d, prec, rem)

	if rem.Sign() == 0 { // remainder is zero
		return quo
	}

	lenWhole := len(d.String())
	if quo.Sign() == 0 { // only the decimal places (ex. 0.1234)
		lenWhole++
	}
	lenQuo := len(quo.String())
	lenRem := len(rem.String())
	leadingZeros := lenWhole - (lenQuo + lenRem) // leading zeros removed from the remainder

	zerosToAdd := int64(lenRem - 1 + leadingZeros)
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(zerosToAdd), nil)
	fiveLine := new(big.Int).Mul(big.NewInt(5), multiplier)

	switch rem.Cmp(fiveLine) {
	case -1:
		chopped = quo
		return
	case 1:
		chopped = new(big.Int).Add(quo, big.NewInt(1))
		return
	default: // bankers rounding must take place
		str := quo.String()
		finalDig, err := strconv.Atoi(string(str[len(str)-1]))
		if err != nil {
			panic(err)
		}

		// always round to an even number
		if finalDig == 0 || finalDig == 2 || finalDig == 4 ||
			finalDig == 6 || finalDig == 8 {

			chopped = quo
			return
		}
		chopped = new(big.Int).Add(quo, big.NewInt(1))
		return
	}
}

// RoundInt64 rounds the decimal using bankers rounding
func (d Dec) RoundInt64() int64 {
	return BankerRoundChop(d.Int, Precision).Int64()
}

// RoundInt round the decimal using bankers rounding
func (d Dec) RoundInt() Int {
	return NewIntFromBigInt(BankerRoundChop(d.Int, Precision))
}

//___________________________________________________________________________________

// wraps d.MarshalText()
func (d Dec) MarshalAmino() (string, error) {
	if d.Int == nil {
		d.Int = new(big.Int)
	}
	bz, err := d.Int.MarshalText()
	return string(bz), err
}

// requires a valid JSON string - strings quotes and calls UnmarshalText
func (d *Dec) UnmarshalAmino(text string) (err error) {
	tempInt := new(big.Int)
	err = tempInt.UnmarshalText([]byte(text))
	if err != nil {
		return err
	}
	d.Int = tempInt
	return nil
}

// MarshalJSON defines custom encoding scheme
func (d Dec) MarshalJSON() ([]byte, error) {
	if d.Int == nil {
		d.Int = new(big.Int)
	}

	text, err := d.Int.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(text))
}

// UnmarshalJSON defines custom decoding scheme
func (d *Dec) UnmarshalJSON(bz []byte) error {
	if d.Int == nil {
		d.Int = new(big.Int)
	}

	var text string
	err := json.Unmarshal(bz, &text)
	if err != nil {
		return err
	}
	return d.Int.UnmarshalText([]byte(text))
}

//___________________________________________________________________________________
// helpers

// test if two decimal arrays are equal
func DecsEqual(d1s, d2s []Dec) bool {
	if len(d1s) != len(d2s) {
		return false
	}

	for i, d1 := range d1s {
		if !d1.Equal(d2s[i]) {
			return false
		}
	}
	return true
}

// intended to be used with require/assert:  require.True(DecEq(...))
func DecEq(t *testing.T, exp, got Dec) (*testing.T, bool, string, Dec, Dec) {
	return t, exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp, got
}

// minimum decimal between two
func MinDec(d1, d2 Dec) Dec {
	if d1.LT(d2) {
		return d1
	}
	return d2
}
package args

import (
	"errors"
	"fmt"
	"strconv"
)

type BpArgs struct {
	FromPack uint16
	ToPack   uint16
	PackPath string
}

func SetArgs(rawArgs []string) (a BpArgs, err error) {
	if len(rawArgs) != 4 {
		err = errors.New("incorrect number of arguments")
		return BpArgs{}, err
	}

	fromPack, err := strconv.ParseUint(rawArgs[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("%v argument not a number", rawArgs[1])
		return BpArgs{}, err
	}
	a.FromPack = uint16(fromPack)

	toPack, err := strconv.ParseUint(rawArgs[2], 10, 64)
	if err != nil {
		err := fmt.Errorf("%v argument not a number", rawArgs[2])
		return BpArgs{}, err
	}
	a.ToPack = uint16(toPack)

	if a.FromPack > a.ToPack {
		err := errors.New("initial pack number is greater than final pack number")
		return BpArgs{}, err
	}

	a.PackPath = rawArgs[3]
	if a.PackPath == "" {
		err := errors.New("no path provided")
		return BpArgs{}, err
	}

	return a, nil
}

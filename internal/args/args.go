package args

import (
	"errors"
	"fmt"
	"strconv"
)

type BpArgs struct {
	FromPack uint64
	ToPack   uint64
	PackPath string
}

func SetArgs(rawArgs []string) (a BpArgs, err error) {

	if len(rawArgs) != 4 {
		err = errors.New("incorrect number of arguments")
		return BpArgs{}, err
	}

	a.FromPack, err = strconv.ParseUint(rawArgs[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("%v argument not a number", rawArgs[1])
		return BpArgs{}, err
	}

	a.ToPack, err = strconv.ParseUint(rawArgs[2], 10, 64)
	if err != nil {
		err := fmt.Errorf("%v argument not a number", rawArgs[2])
		return BpArgs{}, err
	}

	if a.FromPack > a.ToPack {
		err := errors.New("initial pack number is greater than final pack number")
		return BpArgs{}, err
	}

	if rawArgs[3] == "" {
		err := errors.New("no path provided")
		return BpArgs{}, err
	}

	return a, nil
}

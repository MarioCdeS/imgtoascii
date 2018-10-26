package config

import "errors"

type Ramp uint8

const (
	Ramp10 Ramp = iota
	Ramp70
)

func (r *Ramp) String() string {
	switch *r {
	case Ramp10:
		return "ramp 10"

	case Ramp70:
		return "ramp 70"

	default:
		return "unknown ramp"
	}
}

func (r *Ramp) Set(value string) error {
	switch value {
	case "10":
		*r = Ramp10

	case "70":
		*r = Ramp70

	default:
		return errors.New("unknown ramp")
	}

	return nil
}

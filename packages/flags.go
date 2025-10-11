package flags

import (
	"fmt"
	"flag"
)

type Sun struct {
	sun  bool
}

func (s *Sun) String() string {
	if s.sun {
		return "true"
	}
	return "false"
}

func (s *Sun) Set(flag string) error {
	var sun bool
	_, err := fmt.Sscanf(flag, "%t", &sun)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	s.sun = sun
	return nil
}

func SunFlag(name string, value bool, usage string) (*bool) {
	f := Sun{sun : value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.sun
}
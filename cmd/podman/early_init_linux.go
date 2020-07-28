package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/containers/podman/v2/libpod/define"
	"github.com/pkg/errors"
)

func setRLimits() error {
	rlimits := new(syscall.Rlimit)
	rlimits.Cur = define.RLimitDefaultValue
	rlimits.Max = define.RLimitDefaultValue
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, rlimits); err != nil {
		if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, rlimits); err != nil {
			return errors.Wrapf(err, "error getting rlimits")
		}
		rlimits.Cur = rlimits.Max
		if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, rlimits); err != nil {
			return errors.Wrapf(err, "error setting new rlimits")
		}
	}
	return nil
}

func setUMask() {
	// Be sure we can create directories with 0755 mode.
	syscall.Umask(0022)
}

func earlyInitHook() {
	if err := setRLimits(); err != nil {
		fmt.Fprint(os.Stderr, "Failed to set rlimits: "+err.Error())
	}

	setUMask()
}

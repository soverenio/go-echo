// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package logoutput

import (
	"io"
	"os"
	"path/filepath"

	errors "github.com/insolar/vanilla/throw"

	"github.com/insolar/vanilla/throw"

	"github.com/soverenio/log/logcommon"
	"github.com/soverenio/log/outputsyslog"
)

type LogOutput uint8

const (
	StdOutOutput LogOutput = iota
	StdErrOutput
	SysLogOutput
	FileOutput
)

func (l LogOutput) String() string {
	switch l {
	case StdOutOutput:
		return "stdout"
	case StdErrOutput:
		return "stderr"
	case SysLogOutput:
		return "syslog"
	case FileOutput:
		return "file"
	}
	return string(l)
}

func (l LogOutput) IsConsole() bool {
	return l == StdOutOutput || l == StdErrOutput
}

var JSONConsoleWrapper func(io.Writer) io.Writer

func OpenLogBareOutput(output LogOutput, fmt logcommon.LogFormat, param string) (logcommon.BareOutput, error) {
	switch output {
	case StdOutOutput:
		o := logcommon.BareOutput{
			Writer:         os.Stdout,
			FlushFn:        os.Stdout.Sync,
			ProtectedClose: true,
		}

		if fmt == logcommon.JSONFormat && JSONConsoleWrapper != nil {
			o.Writer = JSONConsoleWrapper(o.Writer)
			if o.Writer == nil {
				return o, throw.FailHere("failed on JSONConsoleWrapper")
			}
		}

		return o, nil
	case StdErrOutput:
		o := logcommon.BareOutput{
			Writer:         os.Stderr,
			FlushFn:        os.Stderr.Sync,
			ProtectedClose: true,
		}

		if fmt == logcommon.JSONFormat && JSONConsoleWrapper != nil {
			o.Writer = JSONConsoleWrapper(o.Writer)
			if o.Writer == nil {
				return o, throw.FailHere("failed on JSONConsoleWrapper")
			}
		}

		return o, nil
	case FileOutput:
		w, err := os.OpenFile(param, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return logcommon.BareOutput{}, err
		}

		return logcommon.BareOutput{
			Writer:         w,
			FlushFn:        w.Sync,
			ProtectedClose: false,
		}, nil

	case SysLogOutput:
		executableName := filepath.Base(os.Args[0])
		w, err := outputsyslog.ConnectSyslogByParam(param, executableName)
		if err != nil {
			return logcommon.BareOutput{}, err
		}
		return logcommon.BareOutput{
			Writer:         w,
			FlushFn:        w.Flush,
			ProtectedClose: false,
		}, nil
	default:
		return logcommon.BareOutput{}, errors.New("unknown output " + output.String())
	}
}

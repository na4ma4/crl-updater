package crlmgr

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kballard/go-shellquote"
)

// Action encapsulates an action to be run at certain times.
type Action struct {
	cmd     string
	timeout time.Duration
	workdir string
}

// ActionFromString returns a Action from the configuration string.
func ActionFromString(in string) *Action {
	return &Action{
		cmd: in,
	}
}

// ErrNonZero returned when a action command returns non-zero.
var ErrNonZero = errors.New("action command returned non-zero response")

// Exec the syncronisation of a target from a source.
func (a *Action) Exec(ctx context.Context) (string, error) {
	args := a.splitCmd()

	if a.timeout > 0 {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, a.timeout)
		defer cancel()
	}

	exitCode, ob, oberr, err := a.wrapCmd(ctx, args)
	if err != nil {
		return "", fmt.Errorf("action command failed to execute: %w", err)
	}

	sb := &strings.Builder{}

	if _, err := io.Copy(sb, ob); err != nil {
		return "", fmt.Errorf("unable to read stdout from action command: %w", err)
	}

	if _, err := io.Copy(sb, oberr); err != nil {
		return "", fmt.Errorf("unable to read stderr from action command: %w", err)
	}

	output := sb.String()

	if exitCode != 0 {
		return output, fmt.Errorf("%w: %d", ErrNonZero, exitCode)
	}

	return output, nil
}

// splitCmd uses `shellquote` on non windows platforms.
func (a *Action) splitCmd() (o []string) {
	o, err := shellquote.Split(a.cmd)
	if err != nil {
		o = strings.Split(a.cmd, " ")
	}

	return
}

// runCmd runs a supplied command and returns the exitcode.
//nolint:nestif // it might be "deeply nested", but it's readable and confines this code to this function.
func (a *Action) runCmd(cmd *exec.Cmd) (exitCode int, cmdErr error) {
	cmdErr = cmd.Run()
	if cmdErr != nil {
		// try to get the exit code
		var exitError *exec.ExitError
		if errors.As(cmdErr, &exitError) {
			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = ws.ExitStatus()
			}
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH, in this situation,
			// exit code could not be get, and stderr will be empty string very likely, so we use
			// the default fail code, and format err to string and set to stderr
			exitCode = 3
			cmdErr = fmt.Errorf("check failed to run: %w", cmdErr)
		}
	} else if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		// success, exitCode should be 0 if go is ok
		exitCode = ws.ExitStatus()
	}

	return
}

// wrapCmd [!windows] uses syscall.Kill to kill process group for check.
func (a *Action) wrapCmd(
	ctx context.Context,
	args []string,
) (int, io.Reader, io.Reader, error) {
	wg := sync.WaitGroup{}
	ob := bytes.NewBuffer(nil)
	oberr := bytes.NewBuffer(nil)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint:gosec
	cmd.Dir = a.workdir

	if pb, err := cmd.StdoutPipe(); err == nil {
		a.ioCopyWaitGroup(&wg, ob, pb)
	}

	if pberr, err := cmd.StderrPipe(); err == nil {
		a.ioCopyWaitGroup(&wg, oberr, pberr)
	}

	exitCode, err := a.runCmd(cmd)

	wg.Wait()

	return exitCode, ob, oberr, err
}

// ioCopyWaitGroup adds one worker to a waitgroup and runs an io.Copy until completed,
// once completed it will call waitgroup.Done().
func (a *Action) ioCopyWaitGroup(wg *sync.WaitGroup, dst io.Writer, src io.Reader) {
	wg.Add(1)

	go func() {
		_, _ = io.Copy(dst, src)

		wg.Done()
	}()
}

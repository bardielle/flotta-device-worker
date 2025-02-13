package network

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

const (
	netfilter string = "nft"
	family    string = "inet" // for IPv4 and IPv6
)

type Error struct {
	exec.ExitError
	cmd        exec.Cmd
	msg        string
	exitStatus *int
}

func (e *Error) ExitStatus() int {
	if e.exitStatus != nil {
		return *e.exitStatus
	}
	return e.Sys().(syscall.WaitStatus).ExitStatus()
}

func (e *Error) Error() string {
	return fmt.Sprintf("running %v: exit status %v: %v", e.cmd.Args, e.ExitStatus(), e.msg)
}

func NewNetfilter() (*Netfilter, error) {
	path, err := exec.LookPath("nft")
	if err != nil {
		return nil, err
	}
	netfilter := &Netfilter{
		path: path,
	}
	return netfilter, nil
}

type Netfilter struct {
	path string
}

func (nft *Netfilter) AddTable(table string) error {
	args := []string{"add", "table", family, table}
	return nft.run(args)
}

func (nft *Netfilter) DeleteTable(table string) error {
	args := []string{"delete", "table", family, table}
	return nft.run(args)
}

func (nft *Netfilter) AddChain(table, chain string) error {
	args := []string{"add", "chain", family, table, chain, "{ type filter hook input priority 0 ; }"}
	return nft.run(args)
}

func (nft *Netfilter) DeleteChain(table, chain string) error {
	args := []string{"delete", "chain", family, table, chain}
	return nft.run(args)
}

func (nft *Netfilter) AddRule(table, chain, rule string) error {
	args := []string{"add", "rule", family, table, chain, rule}
	return nft.run(args)
}

func (nft *Netfilter) run(args []string) error {
	args = append([]string{nft.path}, args...)
	var stderr bytes.Buffer
	cmd := exec.Cmd{
		Path:   nft.path,
		Args:   args,
		Stderr: &stderr,
	}

	if err := cmd.Run(); err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			return &Error{*e, cmd, stderr.String(), nil}
		default:
			return err
		}
	}

	return nil
}

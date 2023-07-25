package gosh

import (
	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
	"sshx/gosh/ext/b64"
	"sshx/gosh/ext/pool"
	"testing"
)

var (
	mockShellScript = `
cpu_count=$(grep "^processor" /proc/cpuinfo|wc -l)
if [ "$cpu_count" -gt 1 ]; then
    echo "great than 1"
    exit 0
else
    echo "lower than or equal 1"
    exit 1
fi
`
	mockPyScript = `
import os

class A:

    def __init__(self, a, b):
        self.a = a
        self.b = b

    def print_abspath(self):
        print os.path.abspath(".")

    def print_a(self):
        print self.a

    def print_b(self):
        print self.b

a = A(1, 2)
a.print_abspath()
a.print_a()
a.print_b()
`
)

func TestParallelSSH(t *testing.T) {
	assert := assert.New(t)

	task := []*pool.Task{
		{
			ID:       "node-1-sh",
			Host:     "localhost",
			Port:     2210,
			UserName: "root",
			Password: "123",
			Handler:  "sh",
			Command:  mockShellScript,
		},
		{
			ID:       "node-2-sh",
			Host:     "localhost",
			Port:     2211,
			UserName: "root",
			Password: "123",
			Handler:  "sh",
			Command:  mockShellScript,
		},
		{
			ID:       "node-3-sh",
			Host:     "localhost",
			Port:     2212,
			UserName: "root",
			Password: "123",
			Handler:  "sh",
			Command:  mockShellScript,
		},
		{
			ID:       "node-1-py",
			Host:     "localhost",
			Port:     2210,
			UserName: "root",
			Password: "123",
			Handler:  "py",
			Command:  mockPyScript,
		},
		{
			ID:       "node-2-py",
			Host:     "localhost",
			Port:     2211,
			UserName: "root",
			Password: "123",
			Handler:  "py",
			Command:  mockPyScript,
		},
		{
			ID:       "node-3-py",
			Host:     "localhost",
			Port:     2212,
			UserName: "root",
			Password: "123",
			Handler:  "py",
			Command:  mockPyScript,
		},
	}

	decode, err := sonic.Marshal(task)
	assert.Nil(err)

	ciphers := b64.Encrypt(decode)

	trigger, err := pool.New(ciphers)
	assert.Nil(err)

	assert.Nil(trigger.Run())
}

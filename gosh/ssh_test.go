package gosh

import (
    "github.com/bytedance/sonic"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestParallelSSH(t *testing.T) {
    assert := assert.New(t)

    task := &Tasks{
        ID: "testcase-1",
        Tasks: []*Task{
            {
                ID:       "node-1",
                Host:     "localhost",
                Port:     2210,
                UserName: "root",
                Password: "123",
                Command: `snum=$(cat /etc/sysctl.conf|grep kernel.msgmni|wc -l)
if [ $snum -eq 1 ];then
   num=$(cat /etc/sysctl.conf|grep kernel.msgmni|grep "#"|wc -l)
   if [ $num -eq 1 ];then
      echo "kernel.msgmni already set,but not open,please open"
      exit 1
   else
     par1=$(cat /etc/sysctl.conf|grep kernel.msgmni|awk -F= '{print $2}'|awk '$1=$1')
     if [ $par1 -eq $msgmni ];then
        echo "kernel.msgmni already set,Compliance"
        exit 0
     else
        echo "kernel.msgmni alreday set,current $par1,but parameter value correct $msgmni,please modify"
        exit 1
     fi 
   fi
else
    echo "kernel.msgmni not set,please set,No Compliance"
    exit 1
fi
`,
            },
            {
                ID:       "node-2",
                Host:     "localhost",
                Port:     2211,
                UserName: "root",
                Password: "123",
                Command: `
msgmnb=16384
snum=$(cat /etc/sysctl.conf|grep kernel.msgmnb|wc -l)
if [ $snum -eq 1 ];then
   num=$(cat /etc/sysctl.conf|grep kernel.msgmnb|grep "#"|wc -l)
   if [ $num -eq 1 ];then
      echo "kernel.msgmnb already set,but not open,please open"
      exit 1
   else
     par1=$(cat /etc/sysctl.conf|grep kernel.msgmnb|awk -F= '{print $2}'|awk '$1=$1')
     if [ $par1 -eq $msgmnb ];then
        echo "kernel.msgmnb already set,Compliance"
        exit 0
     else
        echo "kernel.msgmnb alreday set,current $par1,but parameter value correct $msgmnb,please modify"
        exit 1
     fi 
   fi
else
   echo "kernel.msgmnb not set,No Compliance,please set"
   exit 1
fi
`,
            },
            {
                ID:       "node-3",
                Host:     "localhost",
                Port:     2211,
                UserName: "root",
                Password: "123",
                Command: `cpu_count=$(grep "^processor" /proc/cpuinfo|wc -l)
if [ "$cpu_count" -gt 1 ]; then
    echo "great than 1"
    exit 0
else
    echo "lower than or equal 1"
    exit 1
fi
`,
            },
        },
    }

    decode, err := sonic.Marshal(task)
    assert.Nil(err)

    ciphers := B64Encrypt(decode)

    trigger := NewTaskTrigger(ciphers)
    assert.Nil(trigger.Run())
}

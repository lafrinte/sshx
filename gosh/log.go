package gosh

import "sshx/gosh/zlog"

var log *zlog.Logger

func init() {
    log = zlog.New()
}

package zlog

import (
    "fmt"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/diode"
    "os"
    "time"
)

var GlobalLogger *Logger

func init() {
    GlobalLogger = New()
}

func New() *Logger {
    zerolog.SetGlobalLevel(zerolog.InfoLevel)

    if GlobalLogger == nil {
        wr := diode.NewWriter(os.Stderr, 10000, 10*time.Millisecond, func(missed int) {
            fmt.Printf("Logger Dropped %d messages", missed)
        })
        GlobalLogger = &Logger{
            zerolog.New(wr).With().Timestamp().Logger(),
        }
    }

    return GlobalLogger
}

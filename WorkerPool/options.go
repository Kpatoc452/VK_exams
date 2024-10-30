package main

import "WorkerPool/logger"


type OptionWP struct{
    Max int
    Logger logger.LoggerInterface
}


type Option func(o *OptionWP)


func NewOptionWP(opts ...Option) OptionWP{
    opt := OptionWP{
        Max: 10,
        Logger: logger.NewLogger(),
    }

    for _, o := range opts{
        o(&opt)
    }

    return opt
}
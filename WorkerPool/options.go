package main


type OptionWP struct{
    Max int
    Logger LoggerInterface
}


type Option func(o *OptionWP)


func NewOptionWP(opts ...Option) OptionWP{
    opt := OptionWP{
        Max: 10,
        Logger: NewLogger(),
    }

    for _, o := range opts{
        o(&opt)
    }

    return opt
}
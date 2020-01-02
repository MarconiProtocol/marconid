package tc

type NetemArgs struct {
  InterfaceName string
  Delay         uint32
  Loss          float32
  Duplicate     float32
  ReorderProb   float32
  CorruptProb   float32
}

type NetemReply struct {
}

type TbfArgs struct {
  InterfaceName   string
  Bandwidth       uint64
  LatencyInMillis float64
}

type TbfReply struct {
}

type ResetArgs struct {
  InterfaceName string
}

type ResetReply struct {
}

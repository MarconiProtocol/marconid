package tc

import (
  "fmt"
  "github.com/MarconiProtocol/netlink"
  "math"
)

// attach a Network Emulator qdisc to this network interface
// delay: millisecond
// loss: percentage
// duplicate: percentage
// reorderProb: percentage
// corruptProb: percentage
func SetNetem(interfaceName string, delay uint32, loss float32, duplicate float32, reorderProb float32, corruptProb float32) error {
  link, err := netlink.LinkByName(interfaceName)
  if err != nil {
    fmt.Printf("Failed to get the link info of interface %s, err = %s\n", interfaceName, err)
    return err
  }

  qdiscs, err := GetQdiscList(interfaceName)
  if err != nil {
    fmt.Printf("Failed to get the existing qdiscs of interface %s, err = %s\n", interfaceName, err)
    return err
  }

  netemQdisc := &netlink.Netem{
    QdiscAttrs: netlink.QdiscAttrs{
      LinkIndex: link.Attrs().Index,
      Handle:    netlink.MakeHandle(1, 0),
      Parent:    netlink.HANDLE_ROOT,
    },
    Latency:     delay * 15625, // TODO: Latency = 15625 is 1ms delay, need to find out what is behind the magic number
    Loss:        netlink.Percentage2u32(loss),
    Duplicate:   netlink.Percentage2u32(duplicate),
    ReorderProb: netlink.Percentage2u32(reorderProb),
    CorruptProb: netlink.Percentage2u32(corruptProb),
    Limit:       1000,
  }

  // if no qdisc is attached to this interface, add a netem qdisc
  if countQdisc(qdiscs) == 0 {
    if err := netlink.QdiscAdd(netemQdisc); err != nil {
      fmt.Printf("Failed to add tbf qdisc to %s, err = %s\n", interfaceName, err)
      return err
    }
    return nil
  }

  // else, delete all old qdiscs and add a tbf qdisc
  if err := DeleteAllQdisc(interfaceName); err != nil {
    fmt.Printf("Failed to remove old qdiscs from interface %s, err = %s\n", interfaceName, err)
    return err
  }
  if err := netlink.QdiscAdd(netemQdisc); err != nil {
    fmt.Printf("Failed to add netem qdisc to %s, err = %s\n", interfaceName, err)
    return err
  }
  return nil
}

// attach a token bucket filter qdisc to this network interface to limit its bandwidth
// bandwidth: rate in bytes/s
// latencyInMillis: latency in milliseconds
// referenced from https://github.com/AliyunContainerService/terway/blob/master/pkg/tc/tc.go
func SetTbf(interfaceName string, bandwidth uint64, latencyInMillis float64) error {
  link, err := netlink.LinkByName(interfaceName)
  if err != nil {
    fmt.Printf("Failed to get the link info of interface %s, err = %s\n", interfaceName, err)
    return err
  }

  burst := calculateBurst(bandwidth, link.Attrs().MTU)
  buffer := calculateBuffer(bandwidth, burst)
  latency := latencyInUsec(latencyInMillis)
  limit := calculateLimit(bandwidth, latency, burst)

  qdiscs, err := GetQdiscList(interfaceName)
  if err != nil {
    fmt.Printf("Failed to get the existing qdiscs of interface %s, err = %s\n", interfaceName, err)
    return err
  }

  tbfQdisc := &netlink.Tbf{
    QdiscAttrs: netlink.QdiscAttrs{
      LinkIndex: link.Attrs().Index,
      Handle:    netlink.MakeHandle(1, 0),
      Parent:    netlink.HANDLE_ROOT,
    },
    Rate:   bandwidth,
    Limit:  limit,
    Buffer: buffer,
  }

  // if no qdisc is attached to this interface, add a tbf qdisc
  if countQdisc(qdiscs) == 0 {
    if err := netlink.QdiscAdd(tbfQdisc); err != nil {
      fmt.Printf("Failed to add tbf qdisc to %s, err = %s\n", interfaceName, err)
      return err
    }
    return nil
  }

  // else, delete all old qdiscs and add a tbf qdisc
  if err := DeleteAllQdisc(interfaceName); err != nil {
    fmt.Printf("Failed to remove old qdiscs from interface %s, err = %s\n", interfaceName, err)
    return err
  }
  if err := netlink.QdiscAdd(tbfQdisc); err != nil {
    fmt.Printf("Failed to add tbf qdisc to %s, err = %s\n", interfaceName, err)
    return err
  }
  return nil
}

func calculateBurst(rate uint64, mtu int) uint32 {
  return uint32(math.Ceil(math.Max(float64(rate)/netlink.Hz(), float64(mtu))))
}

func time2Tick(time uint32) uint32 {
  return uint32(float64(time) * float64(netlink.TickInUsec()))
}

func calculateBuffer(rate uint64, burst uint32) uint32 {
  return time2Tick(uint32(float64(burst) * float64(netlink.TIME_UNITS_PER_SEC) / float64(rate)))
}

func calculateLimit(rate uint64, latency float64, buffer uint32) uint32 {
  return uint32(float64(rate)*latency/float64(netlink.TIME_UNITS_PER_SEC)) + buffer
}

func latencyInUsec(latencyInMillis float64) float64 {
  return float64(netlink.TIME_UNITS_PER_SEC) * (latencyInMillis / 1000.0)
}

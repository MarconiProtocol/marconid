package tc

import (
  "github.com/MarconiProtocol/netlink"
  "testing"
)

func TestSetTbf(t *testing.T) {
  interfaceName := "foo"
  var bandwidth uint64 = 1000       // 8000 bits/s
  var latencyInMillis float64 = 400 // 400 ms

  createDummyInterface(t, interfaceName)
  defer deleteDummyInterface(t, interfaceName)

  link, err := netlink.LinkByName(interfaceName)
  if err != nil {
    t.Fatal(err)
  }
  burst := calculateBurst(bandwidth, link.Attrs().MTU)
  buffer := calculateBuffer(bandwidth, burst)
  latency := latencyInUsec(latencyInMillis)
  limit := calculateLimit(bandwidth, latency, burst)

  if err := SetTbf(interfaceName, bandwidth, latencyInMillis); err != nil {
    t.Fatal(err)
  }

  qdiscList, err := GetQdiscList(interfaceName)
  if err != nil {
    t.Fatal(err)
  }

  qdisc, ok := qdiscList[0].(*netlink.Tbf)

  if !ok {
    t.Fatal("Qdisc type does not match")
  }
  if qdisc.Rate != bandwidth {
    t.Fatal("Tbf Rate doesn't match, expected", bandwidth, "got", qdisc.Rate)
  }
  if qdisc.Limit != limit {
    t.Fatal("Tbf Limit doesn't match, expected", limit, "got", qdisc.Limit)
  }
  if qdisc.Buffer != buffer {
    t.Fatal("Tbf Buffer doesn't match, expected", buffer, "got", qdisc.Buffer)
  }
}

func TestSetNetem(t *testing.T) {
  interfaceName := "foo"
  var delay uint32 = 100       // 100ms
  var loss float32 = 20        // 20% loss
  var duplicate float32 = 20   // 20% duplicate
  var reorderProb float32 = 30 // 30% reordering
  var corruptProb float32 = 40 // 40% corrupt

  createDummyInterface(t, interfaceName)
  defer deleteDummyInterface(t, interfaceName)

  if err := SetNetem(interfaceName, delay, loss, duplicate, reorderProb, corruptProb); err != nil {
    t.Fatal(err)
  }

  qdiscList, err := GetQdiscList(interfaceName)
  if err != nil {
    t.Fatal(err)
  }

  qdisc, ok := qdiscList[0].(*netlink.Netem)

  if !ok {
    t.Fatal("Qdisc type does not match")
  }
  if qdisc.Latency != delay*15625 {
    t.Fatal("Netem Latency doesn't match, expected", delay*15625, "got", qdisc.Latency)
  }
  if qdisc.Loss != netlink.Percentage2u32(loss) {
    t.Fatal("Netem Loss doesn't match, expected", netlink.Percentage2u32(loss), "got", qdisc.Loss)
  }
  if qdisc.Duplicate != netlink.Percentage2u32(duplicate) {
    t.Fatal("Netem Duplicate doesn't match, expected", netlink.Percentage2u32(duplicate), "got", qdisc.Duplicate)
  }
  if qdisc.ReorderProb != netlink.Percentage2u32(reorderProb) {
    t.Fatal("Netem ReorderProb doesn't match, expected", netlink.Percentage2u32(reorderProb), "got", qdisc.ReorderProb)
  }
  if qdisc.CorruptProb != netlink.Percentage2u32(corruptProb) {
    t.Fatal("Netem CorruptProb doesn't match, expected", netlink.Percentage2u32(corruptProb), "got", qdisc.CorruptProb)
  }
}

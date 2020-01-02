package tc

import (
  "github.com/MarconiProtocol/netlink"
  "testing"
)

func TestAddAndDeleteQdisc(t *testing.T) {
  interfaceName := "foo"
  var rate uint64 = 131072
  var latency float64 = 16793

  createDummyInterface(t, interfaceName)
  defer deleteDummyInterface(t, interfaceName)

  qdiscList, err := GetQdiscList(interfaceName)
  if err != nil {
    t.Fatal(err)
  }
  if len(qdiscList) != 0 {
    t.Fatal("The length of qdiscList should be 0")
  }

  if err := SetTbf(interfaceName, rate, latency); err != nil {
    t.Fatal(err)
  }
  qdiscList, err = GetQdiscList(interfaceName)
  if err != nil {
    t.Fatal(err)
  }
  if len(qdiscList) != 1 {
    t.Fatal("The length of qdiscList should be 1")
  }
  if err := DeleteAllQdisc(interfaceName); err != nil {
    t.Fatal(err)
  }

  qdiscList, err = GetQdiscList(interfaceName)
  if err != nil {
    t.Fatal(err)
  }
  if len(qdiscList) != 0 {
    t.Fatal("The length of qdiscList should be 0")
  }
}

// create a dummy network interface
func createDummyInterface(t *testing.T, interfaceName string) {
  if err := netlink.LinkAdd(&netlink.Ifb{LinkAttrs: netlink.LinkAttrs{Name: interfaceName}}); err != nil {
    t.Fatal(err)
  }
}

// delete the dummy network interface
func deleteDummyInterface(t *testing.T, interfaceName string) {
  link, err := netlink.LinkByName(interfaceName)
  if err != nil {
    t.Fatal(err)
  }

  err = netlink.LinkDel(link)
  if err != nil {
    t.Fatal(err)
  }
}

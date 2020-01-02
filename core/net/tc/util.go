package tc

import (
  "fmt"
  "github.com/MarconiProtocol/netlink"
)

// get all qdiscs of this interface
func GetQdiscList(interfaceName string) ([]netlink.Qdisc, error) {
  link, err := netlink.LinkByName(interfaceName)
  if err != nil {
    fmt.Printf("Failed to get the link info of interface %s, err = %s\n", interfaceName, err)
  }
  qdiscs, err := netlink.QdiscList(link)
  if err != nil {
    fmt.Printf("Failed to get the queue disciplines for interface %s, err = %s\n", interfaceName, err)
    return nil, err
  }
  return qdiscs, nil
}

// delete all qdiscs attached to this interface
func DeleteAllQdisc(interfaceName string) error {
  qdiscList, err := GetQdiscList(interfaceName)
  if err != nil {
    return err
  }
  for _, qdisc := range qdiscList {
    if qdisc.Attrs().Handle == netlink.HANDLE_NONE {
      continue
    }
    if err := netlink.QdiscDel(qdisc); err != nil {
      return err
    }
  }
  return nil
}

// return how many qdiscs attached to this interface, ignore the default qdisc
func countQdisc(qdiscList []netlink.Qdisc) int {
  count := 0
  for _, qdisc := range qdiscList {
    if qdisc.Attrs().Handle == netlink.HANDLE_NONE {
      continue
    }
    count++
  }
  return count
}

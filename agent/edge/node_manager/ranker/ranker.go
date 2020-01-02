package ranker

type RankingStrategy interface {
  Initialize()
  GetTopNodes(num int) []string
  AddNode(ipAddr string)
  RemoveNode(ipAddr string)
}

type Ranker struct {
  RankingStrategy RankingStrategy
}

/*
	Returns a new ranker object
*/
func NewRanker(strategy RankingStrategy) *Ranker {
  strategy.Initialize()
  ranker := Ranker{strategy}
  return &ranker
}

func (r *Ranker) GetTopNodes(num int) []string {
  return r.RankingStrategy.GetTopNodes(num)
}

func (r *Ranker) AddNode(ipAddr string) {
  r.RankingStrategy.AddNode(ipAddr)
}

func (r *Ranker) RemoveNode(ipAddr string) {
  r.RankingStrategy.RemoveNode(ipAddr)
}

func (r *Ranker) Initialize() {
  r.RankingStrategy.Initialize()
}

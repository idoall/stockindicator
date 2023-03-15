package utils

type IStrategy interface {
	AnalysisSide() SideData
}

// RunStrategies 运行多个策略
func RunStrategies(strategies ...IStrategy) []SideData {
	actions := make([]SideData, len(strategies))

	for i := 0; i < len(strategies); i++ {
		actions[i] = strategies[i].AnalysisSide()
	}

	return actions
}

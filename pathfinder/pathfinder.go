package pathfinder

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func RecalculatePath(fromX, fromY, toX, toY float64) []Position {
	// 간단한 직선 경로 리턴 (테스트용)
	path := []Position{
		{X: fromX, Y: fromY},
		{X: toX, Y: toY},
	}
	return path
}

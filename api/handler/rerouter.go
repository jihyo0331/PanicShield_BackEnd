package handler

import (
	"net/http"
	"ps_backend/pathfinder" // ✅ 여기랑

	"github.com/gin-gonic/gin"
)

type RerouteRequest struct {
	CurrentX float64 `json:"current_x"`
	CurrentY float64 `json:"current_y"`
	TargetX  float64 `json:"target_x"`
	TargetY  float64 `json:"target_y"`
}

func Reroute(c *gin.Context) {
	var req RerouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	newPath := pathfinder.RecalculatePath( // ✅ 여기도
		req.CurrentX, req.CurrentY,
		req.TargetX, req.TargetY,
	)

	c.JSON(http.StatusOK, gin.H{"path": newPath})
}

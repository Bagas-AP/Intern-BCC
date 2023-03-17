package Utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseIDs(c *gin.Context) (uint, int, uint, uint, error) {
	id, ok := c.Get("id")
	if !ok {
		return 0, 0, 0, 0, fmt.Errorf("missing ID value in context")
	}

	modelID, err := strconv.Atoi(c.Param("model_id"))
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid model ID: %s", err)
	}

	serviceID46, err := (strconv.ParseUint(c.Param("service_id"), 10, 64))
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid service ID: %s", err)
	}

	serviceID := uint(serviceID46)

	menuID64, err := strconv.ParseUint(c.Param("menu_id"), 10, 64)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid menu ID: %s", err)
	}

	menuID := uint(menuID64)

	return id.(uint), modelID, serviceID, menuID, nil
}

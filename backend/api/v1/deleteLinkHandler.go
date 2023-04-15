package v1

import (
	"github.com/gin-gonic/gin"
)

// DeleteLinkHandler                godoc
// @Summary      Удаляет ссылку.
// @Description  В строку запроса передаётся секретный ключ, по которому удаляется ссылка.
// @Produce      json
// @Success      204
// @Failure      500  {object}  errorResponse
// @Router       /api/v1/admin/{secret_key} [delete]
func (t *TaskServerV1) DeleteLinkHandler(c *gin.Context) {
	secretKey := c.Param("secretKey")

	// delete the link
	err := t.PgContext.DeleteUrl("secret_key", secretKey)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(204, gin.H{"message": "link deleted"})
}

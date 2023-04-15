package v1

import (
	"github.com/gin-gonic/gin"
	"shortener/configs"
	"shortener/db"
	"time"
)

// GetUrlsListHandler                godoc
// @Summary      Возвращает массив секретных ключей.
// @Description  В строку запроса передаётся токен, по которому возвращается массив.
// @Produce      json
// @Success      200  {object}  getUrlsListResponse
// @Failure 	 401  {object}  errorResponse "Если токен не найден"
// @Failure      500  {object}  errorResponse
// @Router       /api/v1/get_urls_list/{token} [get]
func (t *TaskServerV1) GetUrlsListHandler(c *gin.Context) {
	token := c.Param("token")
	authConfig := configs.NewAuthConfig()

	user, userErr := t.PgContext.GetUser("token", token)
	if _, ok := userErr.(*db.NotFoundUserError); ok {
		c.JSON(401, gin.H{"error": "user unauthorized"})
		return
	}
	if user.TokenCreatedAt.Time.Before(time.Now().UTC().Add(-authConfig.TokenLiveTime)) {
		c.JSON(401, gin.H{"error": "user unauthorized"})
		return
	}

	urls, err := t.PgContext.GetUrlsList(token)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = t.PgContext.UpdateUserTokenTime(token)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"urls": urls})
}

// response struct for the GetLinkInfoHandler
type getUrlsListResponse struct {
	Items []string `json:"items"`
}

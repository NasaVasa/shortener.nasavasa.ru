package v1

import (
	"github.com/gin-gonic/gin"
	"time"
)

// GetInfoByShortUrlHandler godoc
// @Summary      Возвращает информацию о публичной ссылке.
// @Description  В строку запроса передаётся краткая ссылка, по которой возвращается информация о ссылке.
// @Produce      json
// @Success      200  {object}  getLinkByShortResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/v1/url_info/{shortUrl} [get]
func (t *TaskServerV1) GetInfoByShortUrlHandler(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	// get the link info
	url, isExist, err := t.PgContext.GetUrl("short_url", shortUrl, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if isExist {
		c.JSON(200, gin.H{
			"short_url":          url.ShortUrl,
			"long_url":           url.LongUrl,
			"number_of_clicks":   url.UrlClicks,
			"dt_created_at":      url.UrlCreatedAt.Time.UTC(),
			"all_redirect_times": []time.Time{},
		})
		return
	}

	c.JSON(404, gin.H{"error": "url not found"})
}

type getLinkByShortResponse struct {
	ShortUrl         string      `json:"short_url"`
	LongUrl          string      `json:"long_url"`
	NumberOfClicks   int         `json:"number_of_clicks"`
	DtCreated        time.Time   `json:"dt_created_at"`
	AllRedirectTimes []time.Time `json:"all_redirect_times"`
}

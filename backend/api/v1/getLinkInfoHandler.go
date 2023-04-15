package v1

import (
	"github.com/gin-gonic/gin"
	"time"
)

// GetLinkInfoHandler                godoc
// @Summary      Возвращает информацию о ссылке.
// @Description  В строку запроса передаётся секретный ключ, по которому возвращается информация о ссылке.
// @Produce      json
// @Success      200  {object}  getLinkInfoResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/v1/admin/{secret_key} [get]
func (t *TaskServerV1) GetLinkInfoHandler(c *gin.Context) {
	secretKey := c.Param("secretKey")

	// get the link info
	url, isExist, err := t.PgContext.GetUrl("secret_key", secretKey, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if isExist {
		c.JSON(200, gin.H{
			"short_url":          url.ShortUrl,
			"long_url":           url.LongUrl,
			"number_of_clicks":   url.UrlClicks,
			"dt_created":         url.UrlCreatedAt.Time.UTC(),
			"dt_will_delete":     url.UrlWillDelete.Time.UTC(),
			"all_redirect_times": []time.Time{},
		})
		return
	}

	url, isExist, err = t.PgContext.GetUrl("secret_key", secretKey, true)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if isExist {
		if url.UrlWillDelete.Time.Before(time.Now().UTC()) {
			err = t.PgContext.DeleteUrl("short_url", url.ShortUrl)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(404, gin.H{"error": "url not found"})
			return
		}
		times, err := t.PgContext.GetAllRedirectTimes(url.Id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"short_url":          url.ShortUrl,
			"long_url":           url.LongUrl,
			"number_of_clicks":   url.UrlClicks,
			"dt_created":         url.UrlCreatedAt.Time.UTC(),
			"dt_will_delete":     url.UrlWillDelete.Time.UTC(),
			"all_redirect_times": times,
		})
		return
	}
	c.JSON(404, gin.H{"error": "url not found"})
}

// response struct for the GetLinkInfoHandler
type getLinkInfoResponse struct {
	ShortUrl         string      `json:"short_url"`
	LongUrl          string      `json:"long_url"`
	NumberOfClick    int         `json:"number_of_clicks"`
	DtCreated        string      `json:"dt_created"`
	DtWillDelete     string      `json:"dt_will_delete"`
	AllRedirectTimes []time.Time `json:"all_redirect_times"`
}

// response struct for the error
type errorResponse struct {
	Error string `json:"error"`
}

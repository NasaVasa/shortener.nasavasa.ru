package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// RedirectToLongHandler                godoc
// @Summary      Перенаправляет пользователя по ссылке.
// @Description  В строку запроса передаётся краткая ссылка, по ней происходит перенаправление на длинную ссылку.
// @Success      200 {object} responseRedirectToLong
// @Failure      404
// @Router       /redirect/{short_url} [get]
func (t *TaskServerV1) RedirectToLongHandler(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	url, isExist, err := t.PgContext.GetUrl("short_url", shortUrl, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if isExist {
		err = t.PgContext.InsertClickInfo(url.Id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		err = t.PgContext.IncrementUrlClicks(url.SecretKey)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"long_url": url.LongUrl})
		return
	}

	url, isExist, err = t.PgContext.GetUrl("short_url", shortUrl, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if isExist {
		if url.UrlWillDelete.Time.After(time.Now().UTC()) {
			err = t.PgContext.InsertClickInfo(url.Id)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			err = t.PgContext.IncrementUrlClicks(url.SecretKey)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"long_url": url.LongUrl})
			return
		}
		err = t.PgContext.DeleteUrl("short_url", url.ShortUrl)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "url not found"})
}

type responseRedirectToLong struct {
	LongUrl string `json:"long_url"`
}

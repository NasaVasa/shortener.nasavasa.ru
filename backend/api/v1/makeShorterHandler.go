package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"shortener/configs"
	"shortener/db"
	"shortener/utils"
	"time"
)

var liveUnits = []string{"SECONDS", "MINUTES", "HOURS", "DAYS"}

// MakeShorterHandler                godoc
// @Summary      Создаёт ссылку.
// @Description  В формате json передаётся url и опционально остальные параметры. В ответ возвращаются секретный ключ и краткая ссылка. Если не передан vip_key, то ссылка автоматически генерируется, а если передан, то vip_key и будет использоваться как краткий ключ.
// @Param        make_shorter_request body  requestMakeShorter  true  "Я не разобрался как тут указать, что всё кроме url nullable, поэтому напоминаю! И ещё, если всё-таки делаете vip-key, то обязательными становятся все поля. time_to_live_unit может принимать значения SECONDS, MINUTES, HOURS, DAYS"
// @Produce      json
// @Success      200  {object}  responseMakeShorter
// @Failure      400  {object}  errorResponse "Сообщения при различных ошибках валидации. Если передано пустое поле url, то возвращается ошибка "url is empty". Если передано некорректное значение time_to_live_unit или time_to_live <= 0, то возвращается ошибка "time to live unit or time to live is invalid". Если vip_key уже сущетсвует в базе данных, то возвращается ошибка "vip key is already in use"."
// @Failure      401  {object}  errorResponse "Если токен не найден"
// @Failure      500  {object}  errorResponse
// @Router       /api/v1/make_shorter [post]
func (t *TaskServerV1) MakeShorterHandler(c *gin.Context) {
	type MakeShorterRequest struct {
		Url            string `json:"url"`
		VipKey         string `json:"vip_key"`
		TimeToLive     int    `json:"time_to_live"`
		TimeToLiveUnit string `json:"time_to_live_unit"`
		Token          string `json:"token"`
	}
	urlConfig := configs.NewUrlConfig()
	prefix := urlConfig.UrlPrefix

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var request MakeShorterRequest
	err = json.Unmarshal(raw, &request)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if request.Url == "" {
		c.JSON(400, gin.H{"error": "url is empty"})
		return
	}

	// if request.Url don't have http:// or https:// prefix, then add it
	if !utils.IsUrl(request.Url) {
		request.Url = "http://" + request.Url
	}

	if request.VipKey == "" {
		url, isExist, err := t.PgContext.GetUrl("long_url", request.Url, false)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if isExist {
			c.JSON(200, gin.H{"short_url": prefix + url.ShortUrl, "secret_key": url.SecretKey})
			return
		}
		shortUrlGen := utils.GenerateShortUrl(t.PgContext)
		err = t.PgContext.InsertUrl(request.Url, shortUrlGen, false, pgtype.Timestamp{}, pgtype.UUID{})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		url, _, err = t.PgContext.GetUrl("short_url", shortUrlGen, false)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"short_url": prefix + shortUrlGen, "secret_key": url.SecretKey})
	} else {
		user, userErr := t.PgContext.GetUser("token", request.Token)
		if _, ok := userErr.(*db.NotFoundUserError); ok {
			c.JSON(401, gin.H{"error": "user unauthorized"})
			return
		}
		if !utils.IsInArray(request.TimeToLiveUnit, liveUnits) || request.TimeToLive <= 0 {
			c.JSON(400, gin.H{"error": "time to live unit or time to live is invalid"})
			return
		}
		url, isExist, err := t.PgContext.GetUrl("short_url", request.VipKey, true)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if isExist {
			if url.UrlWillDelete.Time.After(time.Now().UTC()) {
				c.JSON(400, gin.H{"error": "vip key is already in use"})
				return
			}
			err = t.PgContext.DeleteUrl("short_url", url.ShortUrl)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}

		newUrlWillDelete, err := utils.ParseLiveTime(request.TimeToLiveUnit, request.TimeToLive)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if userErr != nil {
			c.JSON(500, gin.H{"error": userErr.Error()})
			return
		}

		err = t.PgContext.InsertUrl(request.Url, request.VipKey, true, newUrlWillDelete, user.Id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		url, _, err = t.PgContext.GetUrl("short_url", request.VipKey, true)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"short_url": prefix + request.VipKey, "secret_key": url.SecretKey})
	}
}

// response struct for MakeShorterHandler
type responseMakeShorter struct {
	ShortUrl  string `json:"short_url"`
	SecretKey string `json:"secret_key"`
}

// request struct for MakeShorterHandler
type requestMakeShorter struct {
	Url            string `json:"url"`
	VipKey         string `json:"vip_key"`
	TimeToLive     int    `json:"time_to_live"`
	TimeToLiveUnit string `json:"time_to_live_unit"`
	Token          string `json:"token"`
}

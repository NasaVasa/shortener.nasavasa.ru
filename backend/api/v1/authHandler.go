package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"golang.org/x/crypto/bcrypt"
	"shortener/configs"
	"shortener/db"
	"shortener/utils"
	"time"
)

// struct for request body
type requestAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// AuthHandler                godoc
// @Summary      Авторизует пользователя.
// @Description  В формате json передаётся login и password. В ответ возвращается токен.
// @Param        auth_request body  requestAuth  true "make_shorter_request"
// @Produce      json
// @Success      200  {object}  responseAuth
// @Failure      400  {object}  errorResponse "Возвращает "bad password" если пароль не валидный, "user not found" если пользователь не найден, "wrong password" если пароль не совпадает"
// @Failure      500  {object}  errorResponse
// @Router       /api/v1/auth [post]
func (t *TaskServerV1) AuthHandler(c *gin.Context) {
	authConfig := configs.NewAuthConfig()

	var data requestAuth
	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// check if user exists
	user, userErr := t.PgContext.GetUser("login", data.Login)
	// if err is NotFoundUserError, then user not found
	if _, ok := userErr.(*db.NotFoundUserError); ok {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	if userErr != nil {
		c.JSON(500, gin.H{"error": userErr.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	// check password
	if err != nil {
		c.JSON(400, gin.H{"error": "wrong password"})
		return
	}

	// generate token if not exists or expired
	if user.Token == "" || time.Now().UTC().Sub(user.TokenCreatedAt.Time) > authConfig.TokenLiveTime {
		user.Token = utils.GenerateToken(t.PgContext).String()
		user.TokenCreatedAt = pgtype.Timestamp{Time: time.Now().UTC()}
		err = t.PgContext.UpdateUser(user)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(200, gin.H{"token": user.Token})
}

// struct for response body
type responseAuth struct {
	Token string `json:"token"`
}

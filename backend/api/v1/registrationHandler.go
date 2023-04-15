package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"shortener/utils"
)

// RegistrationHandler                godoc
// @Summary      Регистрирует пользователя.
// @Description  В формате json передаётся login, email и password. В ответ возвращается токен.
// @Param        registration_request body  requestUser  true "requestUser"
// @Produce      json
// @Success      200  {object}  responseAuth
// @Failure      400  {object}  errorResponse "Возвращается, если email не прошёл валидацию или пользователь с таким login или email уже существует."
// @Failure      500  {object}  errorResponse
// @Router       /api/v1/registration [post]
func (t *TaskServerV1) RegistrationHandler(c *gin.Context) {
	var data requestUser
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

	// check user email on regex
	isValid := utils.CheckEmail(data.Email)
	if !isValid {
		c.JSON(400, gin.H{"error": "email is not valid"})
		return
	}

	isExist := t.PgContext.CheckUser("login", data.Login)
	if isExist {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}
	isExist = t.PgContext.CheckUser("email", data.Email)
	if isExist {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}

	token := utils.GenerateToken(t.PgContext).String()

	err = t.PgContext.CreateUser(data.Login, utils.HashPassword(data.Password), data.Email, token)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token})
}

type requestUser struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

package users

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/logger"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/utils"
)

var log = logger.SetupLogger()
var tablename = models.User{}.TableName()

type UsersRouter struct {
}

func SetupAPIRoutes(r *gin.RouterGroup) {
	users := &UsersRouter{}
	r.GET("/me", users.findMe)
	r.PATCH("/me", users.updateMe)
}

type User struct {
	ID         uint      `json:"id"`
	Verified   bool      `json:"verified"`
	Name       string    `json:"name"      validate:"omitempty,min=2,max=100"`
	LastName   string    `json:"last_name" validate:"omitempty,min=2,max=100"`
	Email      string    `json:"email"     validate:"omitempty,email"`
	PhotoURL   string    `json:"photo_url"`
	Phone      string    `json:"phone"     validate:"omitempty,min=7,max=15"`
	CreationAt time.Time `json:"creation_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

func (h *UsersRouter) findMe(c *gin.Context) {
	token, err := utils.GetToken(c)
	if err != nil {
		log.Error("Error getting token", err)
		utils.Response(c, err)
		return
	}
	session, err := utils.ValidateSession(token)
	if err != nil {
		log.Error("Error validating session", err)
		utils.Response(c, err)
		return
	}

	conn := db.DefaultClient
	user := &User{}
	tx := conn.Table(tablename).Where("id = ?", session.ID).First(user)
	if tx.Error != nil {
		log.Error("Error getting users", tx.Error)
		utils.Response(c, tx.Error)
		return
	}

	c.JSON(200, user)
}

type UpdateValidationErrors struct {
	NameError     string `json:"name_error,omitempty"`
	LastNameError string `json:"last_name_error,omitempty"`
	PhoneError    string `json:"phone_error,omitempty"`
	EmailError    string `json:"email_error,omitempty"`
	PasswordError string `json:"password_error,omitempty"`
}

func (h *UsersRouter) updateMe(c *gin.Context) {
	token, err := utils.GetToken(c)
	if err != nil {
		log.Error("Error getting token", err)
		utils.Response(c, err)
		return
	}
	session, err := utils.ValidateSession(token)
	if err != nil {
		log.Error("Error validating session", err)
		utils.Response(c, err)
		return
	}

	var userInput User
	if err := c.BindJSON(&userInput); err != nil {
		log.Error("Error binding JSON", err)
		utils.Response(c, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(userInput); err != nil {
		log.Error("Error validating user input", err)
		errorsMap := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()

			switch tag {
			case "required":
				errorsMap[field] = "This field is required!"
			case "email":
				errorsMap[field] = "Invalid email format!"
			case "phone":
				errorsMap[field] = "Invalid phone number!"
			case "pattern":
				errorsMap[field] = "Password does not meet requirements!"
			default:
				errorsMap[field] = "Invalid field!"
			}
		}
		customErrors := UpdateValidationErrors{
			NameError:     errorsMap["Name"],
			LastNameError: errorsMap["LastName"],
			PhoneError:    errorsMap["Phone"],
			EmailError:    errorsMap["Email"],
			PasswordError: errorsMap["Password"],
		}
		c.JSON(http.StatusBadRequest, customErrors)
		return
	}

	conn := db.DefaultClient
	tx := conn.Table(tablename).Where("id = ?", session.ID).Updates(&userInput)
	if tx.Error != nil {
		log.Error("Error updating user", tx.Error)
		utils.Response(c, tx.Error)
		return
	}

	c.JSON(200, gin.H{"message": "Profile updated successfully"})
}
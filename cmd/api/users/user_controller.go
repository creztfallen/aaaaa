package users

import (
	"context"
	"fmt"
	"github.com/creztfallen/compiler_api/cmd/api/auth"
	"github.com/creztfallen/compiler_api/cmd/api/config"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(c *fiber.Ctx) error {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validationErr.Error()}})
	}

	hash := HashPassword(user.Password)

	newUser := User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Username: user.Username,
		Phone:    user.Phone,
		Password: hash,
		Email:    user.Email,
		Cnpj:     user.Cnpj}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": result}})
}

func Login(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	filter := bson.M{"email": user.Email}
	var foundUser User

	err := userCollection.FindOne(context.TODO(), filter).Decode(&foundUser)
	fmt.Println(err)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data:    &fiber.Map{"data": err}})
	}
	passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
	if !passwordIsValid {
		return c.Status(http.StatusUnauthorized).JSON(UserResponse{
			Status:  http.StatusUnauthorized,
			Message: msg,
			Data:    nil})
	}

	token, err := auth.GenerateJWT(foundUser.Email, foundUser.Username)

	return c.Status(http.StatusOK).JSON(UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &fiber.Map{
			"user":  foundUser,
			"token": token,
		},
	})
}

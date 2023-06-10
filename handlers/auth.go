package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"log"
	"pahomov_frolovsky_cson/handlers/response"
	"pahomov_frolovsky_cson/postgres"
)

var secret = "your-secret-key"

func GenerateToken(userID int, roleId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role_id": roleId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret)) // Replace with your own secret key
}

// HandlerLogin handles the POST request for user login
func HandlerLogin(ctx *fiber.Ctx) error {
	var request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	err := json.Unmarshal(ctx.Body(), &request)
	if err != nil {
		log.Println("Failed to decode JSON:", err)
		return response.ErrBadRequest.Send(ctx)
	}

	var userId int
	var password string
	var roleId int
	err = postgres.Conn.QueryRow(ctx.Context(), "SELECT id, password, role_id FROM public.users WHERE login = $1", request.Login).
		Scan(&userId, &password, &roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.ErrUnauthorized.SetData("Invalid credentials").AddMessage("Логин или пароль указан не верно").Send(ctx)
		}
		log.Println("Failed to execute query:", err)
		return response.ErrInternal.Send(ctx)
	}

	// Check if the entered password matches the one in the database
	if request.Password != password {
		return response.ErrUnauthorized.SetData("Invalid credentials").AddMessage("Логин или пароль указан не верно").Send(ctx)
	}

	// Generate a JWT token
	token, err := GenerateToken(userId, roleId)
	if err != nil {
		return response.ErrInternal.SetData("Failed to generate token " + err.Error()).Send(ctx)
	}

	// Return the token in the response
	return response.OK.SetData(map[string]string{"token": token}).Send(ctx)
}

func AuthMiddlewareAll(ctx *fiber.Ctx) error {
	return AuthMiddleware(ctx, []int{1, 2})
}
func AuthMiddlewareNewsMaker(ctx *fiber.Ctx) error {
	return AuthMiddleware(ctx, []int{2})
}
func AuthMiddlewareAdmin(ctx *fiber.Ctx) error {
	return AuthMiddleware(ctx, []int{1})
}

// AuthMiddleware is a middleware function to check if the user is authorized
func AuthMiddleware(ctx *fiber.Ctx, roles []int) error {
	// Extract the token from the Authorization header
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return response.ErrUnauthorized.SetData("Missing authorization header").Send(ctx)
	}

	// Verify the token
	tokenString := authHeader[len("Bearer "):]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil // Replace with your own secret key
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return response.ErrInternal.SetData("Invalid token").Send(ctx)
		}
		log.Println("Failed to parse token:", err)
		return response.ErrInternal.Send(ctx)
	}

	// Check if the token is valid
	if !token.Valid {
		return response.ErrInternal.SetData("Invalid token").Send(ctx)
	}
	for _, role := range roles {
		//fmt.Println("claims[\"role_id\"]", fmt.Sprintf("%v", claims["role_id"]), "==", strconv.Itoa(role))

		if claims["role_id"] == float64(role) {
			return ctx.Next()
		}
	}
	return response.ErrUnauthorized.SetData("У Вас нет необходимой роли!").Send(ctx)

	// Proceed to the next middleware or handler

}

//func IsRoleNewsMaker(t *jwt.Token) error {
//	if validateRole(t, 2) {
//		return nil
//	}
//	return errors.New("У Вас нет необходимой роли!")
//}
//
//func validateRole(t *jwt.Token, roleId int) bool {
//	claims := t.Claims.(jwt.MapClaims)
//	id := int(claims["role_id"].(float64))
//	return id == roleId
//}

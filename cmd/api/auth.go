package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/livin21/billmate/internal/store"
	"github.com/livin21/billmate/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func getUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	parsedUserId, parseErr:= uuid.Parse(ctx.Value("user_id").(string))
	if parseErr != nil {
		return uuid.UUID{}, parseErr
	}
	return parsedUserId, nil
}

func getUserNameFromContext(ctx context.Context) string {
	return ctx.Value("name").(string)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.serverError(w, err)
		return
	}

	u, err := app.store.Users.GetByEmail(r.Context(), user.Email)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		log.Println(err)
		app.unAuthorized(w)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   u.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"user_id": u.ID,
		"name":    u.Name,
		"iat":     time.Now().Unix(),
		"role":    u.Role,
	})

	signedToken, err := token.SignedString([]byte(app.config.jwtSecret))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJson(w, http.StatusOK, util.DataResponse{
		Data: map[string]interface{}{
			"token": signedToken,
		},
		Status: http.StatusOK,
	})
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if user.Email == "" || user.Password == "" {
		app.badRequest(w, "email and password are required")
		return
	}

	err = app.store.Users.Create(r.Context(), &user)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJson(w, http.StatusCreated, &util.MessageResponse{
		Message: "User created successfully",
		Status:  http.StatusCreated,
	})
}

func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unAuthorized(w)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.ParseWithClaims(
			tokenString,
			jwt.MapClaims{
				"email":   "",
				"exp":     time.Time{},
				"user_id": "",
				"name":    "",
				"iat":     time.Time{},
				"role":    "",
			},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, &jwt.ValidationError{}
				}
				return []byte(app.config.jwtSecret), nil
			})
		if err != nil {
			app.unAuthorized(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			app.unAuthorized(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims["user_id"])
		ctx = context.WithValue(ctx, "role", claims["role"])
		ctx = context.WithValue(ctx, "name", claims["name"])
		ctx = context.WithValue(ctx, "email", claims["email"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-auth/db"
	"go-auth/model"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

func errorInResponse(w http.ResponseWriter, status int, error model.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
	return
}

func responseByJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
	return
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var error model.Error

	fmt.Println(r.Body)

	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Emailは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "パスワードは必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	// パスワードのハッシュを生成
	// https://godoc.org/golang.org/x/crypto/bcrypt#GenerateFromPassword
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}

	user.Password = string(hash)

	db.Get().Create(&model.User{Password: user.Password, Email: user.Email})

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")

	// JSON 形式で結果を返却
	responseByJSON(w, user)

}

func createToken(user *model.User) (string, error) {

	var err error

	secret := "secret"

	// Token作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "_init_",
	})

	tokenSting, err := token.SignedString([]byte(secret))

	fmt.Println("tokenString", tokenSting)

	if err != nil {
		log.Fatal(err)
	}

	return tokenSting, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user model.User

	var error model.Error
	var jwt model.JWT
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email は必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "パスワードは、必須です。"
		errorInResponse(w, http.StatusBadRequest, error)
	}

	password := user.Password
	fmt.Println(password)

	raw := db.Get().Where("email = ?", user.Email).First(&user).Row()
	err := raw.Scan(&user)

	if err != nil {
		// https://golang.org/pkg/database/sql/#pkg-variables
		if err == sql.ErrNoRows {
			error.Message = "ユーザが存在しません。"
			errorInResponse(w, http.StatusBadRequest, error)
		}

	}

	hasedPassword := user.Password
	fmt.Println(hasedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(password))

	if err != nil {
		error.Message = "無効なパスワードです。"
		errorInResponse(w, http.StatusUnauthorized, error)
		return
	}

	token, err := createToken(&user)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	jwt.Token = token

	responseByJSON(w, jwt)

}

func VerifyEndpoint(w http.ResponseWriter, r *http.Request) {
	responseByJSON(w, "認証OK")
}
func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var errorObject model.Error

		// HTTP リクエストヘッダーを読み取る
		authHeader := r.Header.Get("Authorization")
		// Restlet Client から以下のような文字列を渡す
		// bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3Q5OUBleGFtcGxlLmNvLmpwIiwiaXNzIjoiY291cnNlIn0.7lJKe5SlUbdo2uKO_iLzzeGoxghG7SXsC3w-4qBRLvs
		bearerToken := strings.Split(authHeader, " ")
		fmt.Println("bearerToken: ", bearerToken)

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("エラーが発生しました。")
				}
				return []byte("secret"), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				errorInResponse(w, http.StatusUnauthorized, errorObject)
				return
			}
	
			if token.Valid {
				// レスポンスを返す
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				errorInResponse(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Token が無効です。"
			errorInResponse(w, http.StatusUnauthorized, errorObject)
			return
		}
	}
}

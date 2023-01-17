package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

var sampleSecretKey = []byte("SecretYouShouldHide")

func handlePage(writer http.ResponseWriter, request *http.Request) {
	_, err := generateJWT()
	if err != nil {
		log.Fatalln("Error generating JWT", err)
	}

	writer.Header().Set("Token", "%v")
	type_ := "application/json"
	writer.Header().Set("Content-Type", type_)
	var message Message
	err = json.NewDecoder(request.Body).Decode(&message)
	if err != nil {
		return
	}
	err = json.NewEncoder(writer).Encode(message)
	if err != nil {
		return
	}
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Unix()
	claims["authorized"] = true
	claims["user"] = "username"
	tokenString, err := token.SignedString(sampleSecretKey)
	fmt.Println(tokenString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// comment these
func verifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if len(authHeader) > 0 {
			jwtToken := authHeader[7:]
			fmt.Println(jwtToken)
			token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					writer.WriteHeader(http.StatusUnauthorized)
					_, err := writer.Write([]byte("Bad algorithm"))
					if err != nil {
						return nil, err
					}
				}
				return sampleSecretKey, nil
			})
			// parsing errors result
			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err2 := writer.Write([]byte(err.Error()))
				if err2 != nil {
					return
				}
			}
			// if there's a token
			if token.Valid {
				endpointHandler(writer, request)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
		// response for if there's no token header
	})
}

func extractClaims(_ http.ResponseWriter, request *http.Request) (string, error) {
	if request.Header["Token"] != nil {
		tokenString := request.Header["Token"][0]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return sampleSecretKey, nil
		})
		if err != nil {
			return "Error Parsing Token: ", err
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			username := claims["username"].(string)
			return username, nil
		}
	}

	return "unable to extract claims", nil
}

func authPage() {
	token, _ := generateJWT()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	req.Header.Set("Token", token)
	_, _ = client.Do(req)
}

func main() {
	str, err := generateJWT()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(str)
	http.HandleFunc("/home", verifyJWT(handlePage))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}

}

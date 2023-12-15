// handlers/user_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"example.com/m/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateUserEndpoint handles the creation of a new user
func CreateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser models.User
	json.NewDecoder(r.Body).Decode(&newUser)

	// Insert the new user into MongoDB
	collection := client.Database("GoBackEnd").Collection("Users")
	existingUser := collection.FindOne(context.TODO(), bson.M{"username": newUser.Username})

	// Kiểm tra lỗi khi tìm kiếm
	if existingUser.Err() == nil {
		// Người dùng đã tồn tại
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}
	newUser.TimeCreate = time.Now()
	token, e := createToken(newUser.Username, newUser.Password)
	if e != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}
	newUser.TokenKey = token
	// Kiểm tra xem biến `client` có nil hay không
	if client == nil {
		log.Println("Error: MongoDB client is nil.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err := collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newUser)
}

// GetUserEndpoint retrieves user data from MongoDB based on the username
func GetUserEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, e := ioutil.ReadAll(r.Body)
	var requestData map[string]interface{}

	e = json.Unmarshal(body, &requestData)

	if e != nil {
		http.Error(w, "Error decoding JSON body", http.StatusBadRequest)
		return
	}

	if e != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	username, ok := requestData["username"].(string)
	if !ok {
		http.Error(w, "Invalid or missing username", http.StatusBadRequest)
		return
	}

	// Query MongoDB for the user with the given username
	collection := client.Database("GoBackEnd").Collection("Users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, e := ioutil.ReadAll(r.Body)
	var requestData map[string]interface{}

	e = json.Unmarshal(body, &requestData)

	if e != nil {
		http.Error(w, "Error decoding JSON body", http.StatusBadRequest)
		return
	}

	if e != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	username, ok := requestData["username"].(string)
	if !ok {
		http.Error(w, "Invalid or missing username", http.StatusBadRequest)
		return
	}

	password, ok := requestData["password"].(string)
	if !ok {
		http.Error(w, "Invalid or missing password", http.StatusBadRequest)
		return
	}

	// Xác nhận người dùng
	success, result := authenticate(username, password)
	if success {
		// Trả về thông tin người dùng nếu đăng nhập thành công
		json.NewEncoder(w).Encode(result)
	} else {
		// Xác nhận không thành công
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, result.(string)), http.StatusUnauthorized)
	}
}

var tokenKey = []byte("your-secret-token-key")

func createToken(username string, password string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["password"] = password
	// Tạo token với khóa bí mật
	tokenString, err := token.SignedString(tokenKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func authenticate(username, password string) (bool, interface{}) {
	// Kiểm tra xem username có tồn tại trong MongoDB hay không
	collection := client.Database("GoBackEnd").Collection("Users")

	var response struct {
		Username    string    `json:"username"`
		Password    string    `json:"password"`
		TimeCreate  time.Time `json:"timecreate"`
		TimeLogin   time.Time `json:"timelogin"`
		Permissions string    `json:"permissions"`
		TokenKey    string    `json:"tokenkey"`
	}

	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&response)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, "Username not found"
		}
		return false, "Error checking username"
	}
	// update := bson.M{
	// 	"$set": bson.M{"timelogin": time.Now()},
	// }
	// Kiểm tra mật khẩu
	if response.Password != password {
		return false, "Incorrect password"
	}
	//  else {
	// 	_, err := collection.UpdateOne(context.TODO(), username, update)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// }

	// Trả về thông tin người dùng nếu kiểm tra thành công
	// Bạn có thể thêm các trường thông tin người dùng khác vào đây nếu cần
	currentTime := time.Now()

	return true, map[string]interface{}{
		"username":    username,
		"timecreate":  response.TimeCreate,
		"timelogin":   currentTime,
		"permissions": response.Permissions,
		"tokenkey":    response.TokenKey,
		// Thêm các trường thông tin người dùng khác nếu cần
	}
}

package models

import (
	"context"
	"errors"
	"fmt"
	"golang-user-authentication/database"
	"golang-user-authentication/dto"
	"golang-user-authentication/helpers"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type (
	User struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	jwtCustomClaims struct {
		Name   string `json:"name"`
		UserId int    `json:"userId"`
		jwt.RegisteredClaims
	}
)

func CreateUser(user *dto.UserCreateUserRequest) (dto.WebResopnse, error) {
	fmt.Printf("CreateUser(%+v) user.model.go \n", user)
	db := database.GetMysqlConnection()
	var res dto.WebResopnse

	// Mencari user dengan email sama
	rows, err := db.Query("select count(*) from users where email = ?", user.Email)
	if err != nil {
		log.Fatal(err)
		return res, err
	}

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
			return res, err
		}
	}

	if count > 0 {
		return res, errors.New("email sudah pernah digunakan")

	}

	// Prepare query sql
	stmt, err := db.Prepare("Insert into users(name, email, password) value(?,?,?)")
	if err != nil {
		log.Fatal(err)
		return res, err
	}

	// Enkripsi password
	user.Password, err = helpers.CreateHashedPassword(user.Password)
	if err != nil {
		return res, err
	}

	// Eksekusi query sql
	_, err = stmt.Exec(user.Name, user.Email, user.Password)
	if err != nil {
		log.Fatal(err)
		return res, err
	}
	defer rows.Close()

	res.Message = "Berhasil membuat user"
	res.Data = map[string]string{
		"name":  user.Name,
		"email": user.Email,
	}

	return res, nil
}

func GetUserInformation(userId int) (dto.WebResopnse, error) {
	fmt.Printf("GetUserInformation(%o) user.model.go \n", userId)
	db := database.GetMysqlConnection()

	var user dto.UserResponse
	var res dto.WebResopnse

	stmt, err := db.Prepare("Select id, name, email from users where id = ?")
	if err != nil {
		return res, err
	}

	rows, err := stmt.Query(userId)
	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return res, err
		}
	}

	res.Message = "Berhasil mengambil data user"

	res.Data = user

	return res, nil

}

func LoginUser(userRequest *dto.UserLoginRequest) (dto.WebResopnse, error) {
	fmt.Printf("LoginUser(%+v) user.model.go \n", userRequest)

	timeToken := time.Hour * time.Duration(1)
	db := database.GetMysqlConnection()
	redis := database.GetRedisConnection()

	var user User
	var res dto.WebResopnse

	stmt, err := db.Prepare("select id, name, password from users where email = ?")
	if err != nil {
		return res, err
	}

	result, err := stmt.Query(userRequest.Email)
	if err != nil {
		return res, errors.New("emails atau Password salah")
	}
	defer result.Close()

	for result.Next() {
		err = result.Scan(&user.Id, &user.Name, &user.Password)
		if err != nil {
			return res, err
		}
	}

	if !helpers.CheckHashedPassword(userRequest.Password, user.Password) {
		return res, errors.New("email atau password salah")
	}

	tokenOnRedisExists := redis.Get(context.Background(), strconv.Itoa(user.Id))
	if err := tokenOnRedisExists.Err(); err != nil {
		fmt.Println("token ga ada di redis")
	}

	tokenRedis, err := tokenOnRedisExists.Result()
	if err == nil {

		res.Data = map[string]string{"token": tokenRedis}
		res.Message = "Berhasil login"

		return res, nil
	}

	claims := &jwtCustomClaims{
		user.Name,
		user.Id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeToken)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return res, err
	}

	op1 := redis.Set(context.Background(), strconv.Itoa(user.Id), t, timeToken)
	if err := op1.Err(); err != nil {
		return res, err
	}

	res.Message = "Berhasil login"
	res.Data = map[string]string{"token": t}

	return res, nil
}

func GetAllUser() (dto.WebResopnse, error) {
	fmt.Println("GetAllUser() user.model.go")

	var res dto.WebResopnse
	db := database.GetMysqlConnection()

	rows, err := db.Query("select * from users")
	if err != nil {
		return res, err
	}
	defer rows.Close()

	var result []dto.UserResponse

	for rows.Next() {

		var each = User{}
		err = rows.Scan(&each.Id, &each.Name, &each.Email, &each.Password)
		if err != nil {
			return res, err
		}

		data := &dto.UserResponse{
			Id:    each.Id,
			Name:  each.Name,
			Email: each.Email,
		}

		result = append(result, *data)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	res.Message = "Berhasil mengambil user"
	res.Data = result

	return res, err

}

func UpdateUser(userRequest *dto.UserUpdateRequest, userId int) (dto.WebResopnse, error) {
	fmt.Printf("UpdateUser(%+v, %o) user.model.go \n", userRequest, userId)

	db := database.GetMysqlConnection()
	var res dto.WebResopnse
	var err error

	user, err := GetUserById(userId)
	if err != nil {
		return res, err
	}
	var resMessage = "Berhasil mengubah"

	if !(userRequest.Name == "") {
		user.Name = userRequest.Name
		resMessage = resMessage + " Name,"
	}
	if !(userRequest.Email == "") {
		user.Email = userRequest.Email
		resMessage = resMessage + " Email,"
	}
	if !(userRequest.Password == "") {
		user.Password, err = helpers.CreateHashedPassword(userRequest.Password)
		if err != nil {
			return res, err
		}
		resMessage = resMessage + " Password,"
	}

	stmt, err := db.Prepare("update users set name = ?, email = ?, password = ? where id = ?")
	if err != nil {
		return res, err
	}

	_, err = stmt.Exec(&user.Name, &user.Email, &user.Password, userId)
	if err != nil {
		return res, err
	}

	res.Message = resMessage
	res.Data, err = GetUserById(user.Id)
	if err != nil {
		return res, err
	}

	return res, nil

}

func GetUserById(userId int) (User, error) {
	fmt.Printf("GetUserByEmail(%o) user.model.go \n", userId)

	db := database.GetMysqlConnection()
	var user User

	err := db.QueryRow("select * from users where id = ?", userId).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil

}

func LogoutUser(userId int) (dto.WebResopnse, error) {
	redis := database.GetRedisConnection()
	var res dto.WebResopnse

	op1 := redis.Del(context.Background(), strconv.Itoa(userId))
	if err := op1.Err(); err != nil {
		fmt.Println(err)
	}

	res.Message = "Berhasil logout"
	res.Data = map[string]string{"message": "Berhasil logout"}

	return res, nil

}

package services

import (
	"context"
	"fmt"
	"strings"
	"time"
	"user-service/config"
	"user-service/constants"
	errConstant "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRepositoryRegistry
}

type IUserService interface {
	//ada 5 function
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.IRepositoryRegistry) IUserService {
	return &UserService{repository: repository}
}

func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {

	fmt.Println("[DEBUG] Login service dimulai")
	fmt.Println("[INFO] Mencari user dengan username:", req.Username)

	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		fmt.Println("[ERROR] Gagal menemukan user:", err)
		return nil, err
	}

	fmt.Println("[INFO] User ditemukan, cek password")

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Println("[ERROR] Password tidak cocok:")
		return nil, err
	}

	fmt.Println("[INFO] Password cocok, buat token JWT")

	expirationTime := time.Now().Add(time.Duration(config.Config.JwtExpirationTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		UserName:    user.UserName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        strings.ToLower(user.Role.Code),
	}

	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		fmt.Println("[ERROR] Gagal generate token:", err)
		return nil, err
	}
	fmt.Println("[INFO] Token JWT berhasil dibuat")
	fmt.Println("[INFO] Token JWT:", tokenString)

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil
}

func (u *UserService) isUsernameExist(ctx context.Context, username string) bool {

	user, err := u.repository.GetUser().FindByUsername(ctx, username)
	if err != nil {
		return false
	}
	if user != nil {
		return true
	}
	return false
}

func (u *UserService) isEmailExist(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, username)
	if err != nil {
		return false
	}
	if user != nil {
		return true
	}
	return false
}

func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if u.isUsernameExist(ctx, req.UserName) {
		return nil, errConstant.ErrUsernameExist
	}

	if u.isUsernameExist(ctx, req.Email) {
		return nil, errConstant.ErrEmailExist
	}

	if req.Password != req.ConfirmPassword {
		return nil, errConstant.ErrPasswordDoesNotMatch
	}

	createdUser, err := u.repository.GetUser().Register(ctx, &dto.RegisterRequest{
		Name:        req.Name,
		UserName:    req.UserName,
		Email:       req.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		RoleID:      constants.Customer,
	})

	if err != nil {
		return nil, err
	}

	user, err := u.repository.GetUser().FindByIDWithRole(ctx, createdUser.ID)
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        user.UUID,
			Name:        user.Name,
			UserName:    user.UserName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role.Code,
		},
	}

	return response, nil
}

func (u *UserService) Update(ctx context.Context, request *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	//Log awal fungsi
	fmt.Println("✅ [DEBUG-SERVICE] Memulai proses Update user dengan UUID:", uuid)
	fmt.Printf("✅ [DEBUG-SERVICE] Data yang diterima: %+v\n", request)

	var (
		password                  string
		checkUsername, checkEmail *models.User
		hashedPassword            []byte
		user, userResult          *models.User
		err                       error
		data                      dto.UserResponse
	)

	//log pencarian user
	fmt.Println("✅ [DEBUG-SERVICE] Mencari user dengan UUID:", uuid)
	//cari user berdasarkan uuid
	//jika tidak ada user, return error
	//jika ada user, lanjutkan

	user, err = u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-SERVICE] Gagal menemukan user:", err)
		return nil, err
	}
	fmt.Println("✅ [DEBUG-SERVICE] User ditemukan, melanjutkan proses update")
	fmt.Printf("✅ [DEBUG-SERVICE] Data user yang ditemukan: %+v\n", user)

	//cek apakah username sudah ada
	fmt.Println("✅ [DEBUG-SERVICE] Memeriksa apakah username sudah ada", request.Username)
	isUsernameExist := u.isUsernameExist(ctx, request.Username)
	fmt.Println("✅ [DEBUG-SERVICE] Username sudah ada:", isUsernameExist)
	//jika username sudah ada, cek apakah username yang diinputkan sama dengan username yang ada
	if isUsernameExist && user.UserName != request.Username {
		fmt.Println("[DEBUG]Mencari user dengan username:", request.Username)
		checkUsername, err = u.repository.GetUser().FindByUsername(ctx, request.Username)
		if err != nil {
			fmt.Println("❌ [ERROR-SERVICE] Error saat cek username", err)
			return nil, err
		}

		if checkUsername != nil {
			fmt.Println(("❌ [ERROR-SERVICE] Username sudah digunakan"))
			return nil, errConstant.ErrUsernameExist
		}
	}

	//cek apakah email sudah ada
	//jika ada, cek apakah email yang diinputkan sama dengan email yang ada
	//jika tidak sama, return error
	isEmailExist := u.isEmailExist(ctx, request.Email)
	if isEmailExist && user.Email != request.Email {
		checkEmail, err = u.repository.GetUser().FindByEmail(ctx, request.Email)
		if err != nil {
			return nil, err
		}

		if checkEmail != nil {
			return nil, errConstant.ErrEmailExist
		}
	}

	if request.Password != nil {
		if *request.Password != *request.ConfirmPassword {
			return nil, errConstant.ErrPasswordDoesNotMatch
		}

		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		password = string(hashedPassword)
	}

	userResult, err = u.repository.GetUser().Update(ctx, &dto.UpdateRequest{
		Name:        request.Name,
		Username:    request.Username,
		Password:    &password,
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
	}, uuid)
	if err != nil {
		return nil, err
	}

	data = dto.UserResponse{
		UUID:        userResult.UUID,
		Name:        userResult.Name,
		UserName:    userResult.UserName,
		Email:       userResult.Email,
		PhoneNumber: userResult.PhoneNumber,
	}

	return &data, nil
}

func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	fmt.Println("🛡️ [DEBUG-SERVICE] Memulai proses GetUserLogin")
	fmt.Println("🛡️ [DEBUG-SERVICE] Mengambil data user dari context")
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	fmt.Println("🛡️ [DEBUG-SERVICE] Data user login:", userLogin)
	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		UserName:    userLogin.UserName,
		Email:       userLogin.Email,
		PhoneNumber: userLogin.PhoneNumber,
		Role:        userLogin.Role,
	}
	fmt.Println("🛡️ [DEBUG-SERVICE] Mengembalikan data user login ke controller")

	return &data, nil
}

func (u *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-SERVICE] Gagal menemukan user:", err)
		return nil, err
	}

	fmt.Println("✅ [INFO-SERVICE] User ditemukan, mengembalikan data user")

	data := dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		UserName:    user.UserName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        strings.ToLower(user.Role.Code),
	}
	fmt.Println("✅ [INFO-SERVICE] Data user:", data)
	fmt.Println("✅ [INFO-SERVICE] Mengembalikan data user ke controller")
	return &data, nil
}

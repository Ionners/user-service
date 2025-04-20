package controllers

import (
	"fmt"
	"net/http"
	errWrap "user-service/common/error"
	"user-service/common/response"
	"user-service/domain/dto"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service services.IServiceRegistry
}

type IUserController interface {
	Login(*gin.Context)
	Register(*gin.Context)
	Update(*gin.Context)
	GetUserLogin(*gin.Context)
	GetUserByUUID(*gin.Context)
}

func NewUserController(service services.IServiceRegistry) IUserController {
	return &UserController{service: service}
}

func (u *UserController) Login(ctx *gin.Context) {
	fmt.Println("[DEBUG] Login handler dipanggil")
	request := &dto.LoginRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		fmt.Println("[ERROR] Gagal binding JSON:", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}
	fmt.Printf("[INFO] Payload masuk: Username=%s, Password=%s\n", request.Username, request.Password)

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		fmt.Println("[ERROR] Validasi request gagal:", err)
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	user, err := u.service.GetUser().Login(ctx, request)
	if err != nil {
		fmt.Println("[ERROR] Login service gagal:", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}
	fmt.Println("[INFO] Login berhasil, token akan dikirim")

	response.HttpResponse(response.ParamHttpResp{
		Code:  http.StatusOK,
		Data:  user.User,
		Token: &user.Token,
		Gin:   ctx,
	})
}

func (u *UserController) Register(ctx *gin.Context) {
	request := &dto.RegisterRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	user, err := u.service.GetUser().Register(ctx, request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: user.User,
		Gin:  ctx,
	})

}

func (u *UserController) Update(ctx *gin.Context) {
	fmt.Println("🔍 [DEBUG-CONTROLLER] Update handler dipanggil")

	request := &dto.UpdateRequest{}
	uuid := ctx.Param("uuid")
	fmt.Println("🔍 [DEBUG-CONTROLLER] UUID yang diterima dari parameter:", uuid)

	fmt.Println("🔍 [DEBUG-CONTROLLER] Mencoba binding JSON request")
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		fmt.Println("❌ [ERROR-CONTROLLER] Gagal binding JSON:", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}
	fmt.Printf("🔍 [DEBUG-CONTROLLER] Request data setelah binding: %v+\n", request)

	fmt.Println("🔍 [DEBUG-CONTROLLER] Memulai validasi request")
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		fmt.Println("❌ [ERROR-CONTROLLER] Validasi request gagal:", err)
		// Menggunakan http.StatusText untuk mendapatkan pesan status
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	fmt.Println("🔍 [DEBUG-CONTROLLER] Validasi berhasil, memanggil service untuk update user")
	user, err := u.service.GetUser().Update(ctx, request, uuid)
	if err != nil {
		fmt.Println("❌ [ERROR-CONTROLLER] Gagal memanggil service untuk update user:", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}
	fmt.Println("✅ [INFO-CONTROLLER] Update berhasil")
	fmt.Printf("✅ [INFO-CONTROLLER] Data user setelah update: %v+\n", user)

	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
	fmt.Println("✅ [INFO-CONTROLLER] Mengembalikan response ke client")
	fmt.Println("✅ [INFO-CONTROLLER] Response berhasil dikirim")
	fmt.Println("✅ [INFO-CONTROLLER] Proses update selesai")
	fmt.Println("=======================================")
}

func (u *UserController) GetUserLogin(ctx *gin.Context) {
	fmt.Println("🔍 [DEBUG-CONTROLLER] GetUserLogin handler dipanggil")
	fmt.Println("🔍 [DEBUG-CONTROLLER] Mengambil data user dari context")
	user, err := u.service.GetUser().GetUserLogin(ctx.Request.Context())
	if err != nil {
		fmt.Println("❌ [ERROR-CONTROLLER] Gagal mendapatkan user login:", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		fmt.Println("✅ [INFO-CONTROLLER] Mengembalikan response error")
		fmt.Println("=======================================")
		return
	}

	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (u *UserController) GetUserByUUID(ctx *gin.Context) {
	user, err := u.service.GetUser().GetUserByUUID(ctx.Request.Context(), ctx.Param("uuid"))
	if err != nil {
		fmt.Println("❌ [ERROR-CONTROLLERS] Gagal mendapatkan user:", err)
		response.HttpResponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		fmt.Println("✅ [INFO-CONTROLLERS] Mengembalikan response error")
		return
	}
	fmt.Println("✅ [INFO-CONTROLLERS] User ditemukan, mengembalikan data user")
	fmt.Println("✅ [INFO-CONTROLLERS] Data user:", user)

	response.HttpResponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

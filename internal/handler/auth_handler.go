package handler

import (
	"api-practice/internal/auth"
	. "api-practice/internal/dto"
	"api-practice/internal/model"
	"api-practice/internal/service"
	"api-practice/internal/validator"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"net/smtp"
	"os"
	"sync"
	"time"
)

type UserHandler struct {
	service      service.UserService
	tokenService auth.TokenService
	emailCodes   map[string]string
	tempUsers    map[string]model.User
	emailCodesMu sync.Mutex
	tempUsersMu  sync.Mutex
}

func NewUserHandler(service service.UserService, tokenService auth.TokenService) *UserHandler {
	return &UserHandler{
		service:      service,
		tokenService: tokenService,
		emailCodes:   make(map[string]string),
		tempUsers:    make(map[string]model.User),
	}
}

// Register godoc
// @Summary Register new user
// @Description Registers new user and sends verification code
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration data"
// @Success 200 {object} dto.SuccessBaseResponse
// @Failure 400 {object} dto.ErrorBaseResponse
// @Failure 500 {object} dto.ErrorBaseResponse
// @Router /api/auth/reg [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if err := validator.Struct(req); err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	h.tempUsersMu.Lock()
	h.tempUsers[req.Email] = model.User{
		Email:    req.Email,
		Password: req.Password,
	}
	h.tempUsersMu.Unlock()

	h.emailCodesMu.Lock()
	h.emailCodes[req.Email] = code
	h.emailCodesMu.Unlock()

	go func(email string) {
		time.Sleep(5 * time.Minute)
		h.emailCodesMu.Lock()
		delete(h.emailCodes, email)
		h.emailCodesMu.Unlock()

		h.tempUsersMu.Lock()
		delete(h.tempUsers, email)
		h.tempUsersMu.Unlock()
	}(req.Email)

	if err := sendEmailSMTP(req.Email, "Email Verification", fmt.Sprintf("Your verification code is: %s", code)); err != nil {
		return Error(c, "Failed to send verification code: "+err.Error(), fiber.StatusInternalServerError)
	}

	return SuccessNoData(c, "Verification code sent. Please verify your email.", fiber.StatusOK)
}

// VerifyRegistration godoc
// @Summary Verify user registration
// @Description Verifies user's email with received code
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.VerifyRequest true "User registration data"
// @Success 201 {object} dto.SuccessBaseResponse
// @Failure 400 {object} dto.ErrorBaseResponse
// @Failure 401 {object} dto.ErrorBaseResponse
// @Router /api/auth/verify [post]
func (h *UserHandler) VerifyRegistration(c *fiber.Ctx) error {
	var req VerifyRequest

	if err := c.BodyParser(&req); err != nil {
		return Error(c, "invalid request body", fiber.StatusBadRequest)
	}

	h.emailCodesMu.Lock()
	expectedCode, ok := h.emailCodes[req.Email]
	h.emailCodesMu.Unlock()

	if !ok || req.Code != expectedCode {
		return Error(c, "invalid or expired verification code", fiber.StatusUnauthorized)
	}

	h.tempUsersMu.Lock()
	user, ok := h.tempUsers[req.Email]
	h.tempUsersMu.Unlock()

	if !ok {
		return Error(c, "user data not found, please register again", fiber.StatusBadRequest)
	}

	if err := h.service.Register(&user); err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	h.emailCodesMu.Lock()
	delete(h.emailCodes, req.Email)
	h.emailCodesMu.Unlock()

	h.tempUsersMu.Lock()
	delete(h.tempUsers, req.Email)
	h.tempUsersMu.Unlock()

	return SuccessNoData(c, "user successfully registered and verified", fiber.StatusCreated)
}

// Authenticate godoc
// @Summary Authenticate user
// @Description auth user and return jwt
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequest true "User login credentials"
// @Success 200 {object} dto.TokenBaseResponse
// @Failure 400 {object} dto.ErrorBaseResponse
// @Failure 500 {object} dto.ErrorBaseResponse
// @Router /api/auth/login [post]
func (h *UserHandler) Authenticate(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	authUser, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	token, err := h.tokenService.GenerateToken(authUser.ID)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	return Success(c, "authenticated", &TokenResponse{Token: token}, fiber.StatusOK)
}

func sendEmailSMTP(to, subject, body string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	plainAuth := smtp.PlainAuth("", from, pass, host)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(host+":"+port, plainAuth, from, []string{to}, msg)
}

func (h *UserHandler) GetService() service.UserService {
	return h.service
}

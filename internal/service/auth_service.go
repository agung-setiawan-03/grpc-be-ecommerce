package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity"
	jwtentity "github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity/jwt"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/repository"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
	cacheService   *gocache.Cache
}

func (as *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if req.Password != req.PasswordConfirmation {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Password dan konfirmasi password tidak sama"),
		}, nil

	}

	user, err := as.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Email sudah terdaftar"),
		}, nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	newUser := entity.User{
		Id:        uuid.NewString(),
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		RoleCode:  entity.UserRoleCustomer,
		CreatedAt: time.Now(),
		CreatedBy: &req.FullName,
	}

	err = as.authRepository.InsertUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("User berhasil didaftarkan"),
	}, nil
}

func (as *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := as.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("Email tidak terdaftar"),
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "Password salah")
		}
		return nil, err
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtentity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.RoleCode,
	})
	secretKey := os.Getenv("JWT_SECRET_KEY")
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Base:        utils.SuccessResponse("Login berhasil"),
		AccessToken: accessToken,
	}, nil
}

func (as *authService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {

	jwtToken, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tokenClaims, err := jwtentity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, err
	}

	as.cacheService.Set(jwtToken, "", time.Duration(tokenClaims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)

	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout berhasil"),
	}, nil
}

func NewAuthService(authRepository repository.IAuthRepository, cahceService *gocache.Cache) IAuthService {
	return &authService{
		authRepository: authRepository,
		cacheService:   cahceService,
	}
}

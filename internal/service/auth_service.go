package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/repository"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	"github.com/AgungSetiawan/grpc-be-ecommerce/pb/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	// Cek password dan konfirmasi password
	if req.Password != req.PasswordConfirmation {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Password dan konfirmasi password tidak sama"),
		}, nil

	}

	// Cek email ke dalam database
	user, err := as.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// Jika email sudah ada, return error
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

	// Jika email belum ada, simpan kredensial user ke dalam database
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
	// Cek apakah email ada di dalam database
	user, err := as.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("Email tidak terdaftar"),
		}, nil
	}

	// Cek apakah password yang dikirim sama dengan password yang ada di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "Password salah")
		}
		return nil, err
	}

	// Generate token JWT
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.JwtClaims{
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

	// Return response
	return &auth.LoginResponse{
		Base:        utils.SuccessResponse("Login berhasil"),
		AccessToken: accessToken,
	}, nil
}

func (as *authService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// Dapatkan token dari metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Tidak ada metadata pada context")
	}

	bearerToken, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Tidak ada token pada metadata")
	}

	if len(bearerToken) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Token tidak ditemukan")
	}

	tokenSplit := strings.Split(bearerToken[0], " ")
	if len(tokenSplit) != 2 {
		return nil, status.Errorf(codes.Unauthenticated, "Format token salah")
	}

	if tokenSplit[0] != "Bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "Token harus diawali dengan Bearer")
	}

	jwtToken := tokenSplit[1]

	// Kembalikan token tadi hingga menjadi entity jwt
	tokenClaims, err := jwt.ParseWithClaims(jwtToken, &entity.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Metode signing token tidak valid %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Token tidak valid")
	}

	if !tokenClaims.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "Token sudah tidak berlaku")
	}

	var claims *entity.JwtClaims
	if claims, ok = tokenClaims.Claims.(*entity.JwtClaims); !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Token tidak valid")
	}

	// Masukkan token dari metadata ke dalam memory db / cache
	as.cacheService.Set(jwtToken, "", time.Duration(claims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)

	// Kembalikan response
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

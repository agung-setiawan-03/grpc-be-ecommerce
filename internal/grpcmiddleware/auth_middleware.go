package grpcmiddleware

import (
	"context"

	jwtentity "github.com/AgungSetiawan/grpc-be-ecommerce/internal/entity/jwt"
	"github.com/AgungSetiawan/grpc-be-ecommerce/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

type authMiddleware struct {
	cacheService *gocache.Cache
}

var publicAPIs = map[string]bool{
	"/auth.AuthService/Login":                   true,
	"/auth.AuthService/Register":                true,
	"/product.ProductService/DetailProduct":     true,
	"/product.ProductService/ListProduct":       true,
	"/product.ProductService/HighlightProducts": true,
}

func (am *authMiddleware) AuthMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if publicAPIs[info.FullMethod] {
		return handler(ctx, req)
	}

	tokenStr, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, ok := am.cacheService.Get(tokenStr)
	if ok {
		return nil, utils.UnauthenticatedResponse()
	}

	claims, err := jwtentity.GetClaimsFromToken(tokenStr)
	if err != nil {
		return nil, err
	}

	ctx = claims.SetTokenContext(ctx)

	resp, err = handler(ctx, req)

	return resp, err
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}

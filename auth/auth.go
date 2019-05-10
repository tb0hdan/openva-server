package auth

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"

	"google.golang.org/grpc/metadata"
)

type Authenticator struct {
	authFileName string
	authData     map[string]string
}

func (a *Authenticator) ReadAuthData(fileName string) (authMap map[string]string, err error) {
	authMap = make(map[string]string)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if len(strings.Split(line, ":")) != 2 {
			err = errors.New(fmt.Sprintf("wrong auth format on line %d", lines))
			break
		}
		authMap[strings.Split(line, ":")[0]] = strings.Split(line, ":")[1]
	}

	if err != nil {
		return nil, err
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return
}

func (a *Authenticator) VerifyToken(token string) (status bool, errMsg string) {
	for systemUUID, systemToken := range a.authData {
		if token == systemToken {
			status = true
			log.Printf("Validated token for %s", systemUUID)
			break
		}
	}
	if !status {
		errMsg = "no valid token found"
	}
	return
}

func (a *Authenticator) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		tokenValues := values["token"]

		if len(tokenValues) == 0 {
			http.Error(w, "403 Forbidden. No token.", http.StatusForbidden)
			return

		}

		if ok, errMsg := a.VerifyToken(tokenValues[0]); !ok {
			http.Error(w, fmt.Sprintf("403 Forbidden. %s", errMsg), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authenticator) GetTokenFromContext(ctx context.Context) (token string, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("authorization failed, no metadata")
	}
	if len(meta.Get("authorization")) == 0 {
		return "", errors.New("authorization failed, invalid token")
	}

	token = meta.Get("authorization")[0]
	return
}

func (a *Authenticator) MyGRPCAuthFunction(ctx context.Context) (newContext context.Context, err error) {
	token, err := a.GetTokenFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrap(err, "grpc auth failed")
	}

	if ok, errMsg := a.VerifyToken(token); !ok {
		return ctx, errors.New(fmt.Sprintf("authorization failed, %s", errMsg))
	}

	return ctx, nil
}

func NewAuthenticator(authFileName string) (*Authenticator, error) {
	authenticator := &Authenticator{
		authFileName: authFileName,
	}
	authData, err := authenticator.ReadAuthData(authFileName)
	if err != nil {
		return nil, err
	}
	authenticator.authData = authData

	return authenticator, nil
}

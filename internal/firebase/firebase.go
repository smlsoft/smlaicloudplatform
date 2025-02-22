package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
)

type IFirebaseAdapter interface {
	ValidateToken(token string) (*UserInfo, error)
}

type FirebaseAdapter struct {
	app *firebase.App
}

func NewFirebaseAdapter() IFirebaseAdapter {

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		// log.Fatalf("error initializing app: %v\n", err)
	}

	return &FirebaseAdapter{
		app: app,
	}
}

func (f *FirebaseAdapter) ValidateToken(token string) (*UserInfo, error) {

	ctx := context.Background()
	client, err := f.app.Auth(ctx)
	if err != nil {
		//log.Fatalf("error getting Auth client: %v\n", err)
		return nil, err
	}
	authToken, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		//log.Fatalf("error verifying ID token: %v\n", err)
		return nil, err
	}

	userName := ""
	if authToken.Claims["name"] != nil {
		userName = authToken.Claims["name"].(string)
	} else {
		userName = authToken.Claims["email"].(string)
	}

	userInfo := &UserInfo{
		SignInProvider: authToken.Firebase.SignInProvider,
		Email:          authToken.Claims["email"].(string),
		UserId:         authToken.UID,
		Name:           userName,
	}

	return userInfo, nil
}

package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KrzysztofSieczkiewicz/go-authentication-service-test/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// properly should be stored in .env
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(
	email string,
	firstName string,
	lastName string,
	userType string,
	userId string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_type:  userType,
		Uid:        userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", updatedAt})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
}

func ValidateToken(singedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		singedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("The token is invalid")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("The token is expired")
		return
	}

	return claims, msg
}

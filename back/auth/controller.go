package auth

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// CognitoClient is an interface that defines the methods that we need to implement
type CognitoClient interface {
	SignUp(username, password string) error
	SignIn(username, password string) (string, error)
	VerifyToken(token string) error
}

type awsCognitoClient struct {
	cognitoClient *cognito.CognitoIdentityProvider
	appClientID   string
}

// NewCognitoClient returns a new CognitoClient
func NewCognitoClient(cognitoRegion string, cognitoAppClientID string) CognitoClient {
	conf := &aws.Config{Region: aws.String(cognitoRegion)}

	sess, err := session.NewSession(conf)
	if err != nil {
		panic(err)
	}
	client := cognito.New(sess)

	return &awsCognitoClient{
		cognitoClient: client,
		appClientID:   cognitoAppClientID,
	}

}

func (c *awsCognitoClient) SignUp(email, password string) error {
	input := &cognito.SignUpInput{
		ClientId: aws.String(c.appClientID),
		Username: aws.String(email),
		Password: aws.String(password),
	}

	_, err := c.cognitoClient.SignUp(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *awsCognitoClient) SignIn(email, password string) (string, error) {
	input := &cognito.InitiateAuthInput{
		AuthFlow: aws.String(cognito.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(email),
			"PASSWORD": aws.String(password),
		},
		ClientId: aws.String(c.appClientID),
	}

	result, err := c.cognitoClient.InitiateAuth(input)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return *result.AuthenticationResult.IdToken, nil
}

func (c *awsCognitoClient) VerifyToken(token string) error {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(token),
	}

	_, err := c.cognitoClient.GetUser(input)
	if err != nil {
		return err
	}

	return nil
}

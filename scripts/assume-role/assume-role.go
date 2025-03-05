package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var roles = map[string]string{
	"staging:admin":     "arn:aws:iam::533267310428:role/OrganizationAccountAccessRole",
	"staging:terraform": "arn:aws:iam::533267310428:role/TerraformRole",
}

func main() {
	// Check if ARN is provided
	if len(os.Args) < 3 {
		printUsage(nil)
		os.Exit(1)
	}

	roleName := os.Args[1] // Get ARN from arguments
	command := os.Args[2]
	args := os.Args[3:]

	roleARN, ok := roles[roleName]
	if !ok {
		printUsage(fmt.Errorf("unknown role: %s", roleName))
		os.Exit(1)
	}

	err := assumeRole(roleARN)
	if err != nil {
		printUsage(err)
		os.Exit(1)
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}

func assumeRole(roleARN string) error {
	// Get the current user to generate a session name
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}
	// Use the username and timestamp for the session name
	sessionName := fmt.Sprintf("%s-%d", currentUser.Username, time.Now().Unix())

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}

	svc := sts.New(sess)
	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleARN),
		RoleSessionName: aws.String(sessionName),
	}

	resp, err := svc.AssumeRole(params)
	if err != nil {
		return fmt.Errorf("error assuming role: %v", err)
	}

	err = os.Setenv("AWS_DEFAULT_REGION", "us-west-2")
	if err != nil {
		return fmt.Errorf("error setting environment variables: %v", err)
	}

	// Set the environment variables for the current session
	err = os.Setenv("AWS_ACCESS_KEY_ID", *resp.Credentials.AccessKeyId)
	if err != nil {
		return fmt.Errorf("error setting environment variables: %v", err)
	}
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", *resp.Credentials.SecretAccessKey)
	if err != nil {
		return fmt.Errorf("error setting environment variables: %v", err)
	}
	err = os.Setenv("AWS_SESSION_TOKEN", *resp.Credentials.SessionToken)
	if err != nil {
		return fmt.Errorf("error setting environment variables: %v", err)
	}

	return nil
}

func printUsage(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("Usage: assume-role <role-name> [command]\nAvailable roles:\n")
	for k := range roles {
		fmt.Printf(" %s", k)
	}
	fmt.Println("\n\nExample:\n\tassume-role staging aws s3 ls")
}

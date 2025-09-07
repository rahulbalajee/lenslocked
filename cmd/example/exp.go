package main

import (
	"context"
	"fmt"

	userctx "github.com/rahulbalajee/lenslocked/context/context"
	"github.com/rahulbalajee/lenslocked/models"
)

func main() {
	ctx := context.Background()

	user := &models.User{
		ID:           1,
		Email:        "test@test.co",
		PasswordHash: "",
	}

	ctx = userctx.WithUser(ctx, user)

	user = userctx.User(ctx)

	fmt.Println(user.Email)
}

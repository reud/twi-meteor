package v2_test

import (
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/reud/twi-meteor/env"
	"github.com/reud/twi-meteor/infra/http"
	"testing"
)

func TestLikingUsers(t *testing.T) {
	err := godotenv.Overload()
	if err != nil {
		t.Error(err)
	}

	con := env.GetTwitterConfig()
	likeData, err := http.LikingUsers("1448269297948581889", con.BearerToken)
	if err != nil {
		t.Error(err)
	}
	pp.Print(likeData)
}

package tests

import (
	modelAPI "crypto_scam/internal/api/handler/user/model"
	"crypto_scam/internal/converter"
	modelDB "crypto_scam/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"reflect"
	"testing"
)

func TestUserToGetUserResponse(t *testing.T) {
	type args struct {
		from *modelDB.User
	}

	var (
		id             = gofakeit.Int64()
		firstname      = gofakeit.Name()
		blocks         = gofakeit.Int64()
		record         = gofakeit.Int64()
		daysStreak     = gofakeit.Int16()
		invitedFriends = gofakeit.Int16()
		isPremium      = gofakeit.Bool()
		league         = gofakeit.Int16()
		award          = gofakeit.Bool()
	)

	input1 := &modelDB.User{
		Id:             id,
		Firstname:      firstname,
		Blocks:         blocks,
		Record:         record,
		DaysStreak:     daysStreak,
		InvitedFriends: invitedFriends,
		IsPremium:      isPremium,
		League:         league,
		Award:          award,
	}
	output1 := &modelAPI.GetUserResponse{
		Id:             id,
		Firstname:      firstname,
		Blocks:         blocks,
		Record:         record,
		DaysStreak:     daysStreak,
		InvitedFriends: invitedFriends,
		IsPremium:      isPremium,
		League:         league,
		Award:          award,
	}

	input2 := &modelDB.User{
		Id:             id,
		Firstname:      "",
		Blocks:         blocks,
		Record:         record,
		DaysStreak:     daysStreak,
		InvitedFriends: invitedFriends,
		IsPremium:      isPremium,
		League:         league,
		Award:          award,
	}
	output2 := &modelAPI.GetUserResponse{
		Id:             id,
		Firstname:      "",
		Blocks:         blocks,
		Record:         record,
		DaysStreak:     daysStreak,
		InvitedFriends: invitedFriends,
		IsPremium:      isPremium,
		League:         league,
		Award:          award,
	}

	tests := []struct {
		name string
		args args
		want *modelAPI.GetUserResponse
	}{
		{
			name: "Тест №1",
			args: args{
				from: input1,
			},
			want: output1,
		},
		{
			name: "Тест №2",
			args: args{
				from: input2,
			},
			want: output2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := converter.UserToGetUserResponse(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserToGetUserResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

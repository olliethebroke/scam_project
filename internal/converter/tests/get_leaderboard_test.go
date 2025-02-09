package tests

import (
	modelAPI "crypto_scam/internal/api/handler/user/model"
	"crypto_scam/internal/converter"
	modelDB "crypto_scam/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"reflect"
	"testing"
)

func TestLeadersToGetLeaderboardResponse(t *testing.T) {
	type args struct {
		from map[int16][]*modelDB.Leader
	}

	var (
		id        = gofakeit.Int64()
		position  = gofakeit.Int16()
		firstname = gofakeit.BeerName()
		blocks    = gofakeit.Int64()
	)
	input1 := make(map[int16][]*modelDB.Leader)
	input1[1] = []*modelDB.Leader{
		{
			Id:        id,
			Position:  position,
			Firstname: firstname,
			Blocks:    blocks,
		},
	}

	output1 := make(map[int16][]*modelAPI.GetLeaderboardResponse)
	output1[1] = []*modelAPI.GetLeaderboardResponse{
		{
			Id:        id,
			Position:  position,
			Firstname: firstname,
			Blocks:    blocks,
		},
	}

	input2 := make(map[int16][]*modelDB.Leader)
	input2[1] = []*modelDB.Leader{
		{
			Id:        id,
			Position:  position,
			Firstname: "",
			Blocks:    blocks,
		},
	}

	output2 := make(map[int16][]*modelAPI.GetLeaderboardResponse)
	output2[1] = []*modelAPI.GetLeaderboardResponse{
		{
			Id:        id,
			Position:  position,
			Firstname: "",
			Blocks:    blocks,
		},
	}

	tests := []struct {
		name string
		args args
		want map[int16][]*modelAPI.GetLeaderboardResponse
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
			if got := converter.LeadersToGetLeaderboardResponse(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LeadersToGetLeaderboardResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

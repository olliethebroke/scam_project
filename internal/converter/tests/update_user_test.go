package tests

import (
	modelAPI "crypto_scam/internal/api/handler/user/model"
	"crypto_scam/internal/converter"
	modelDB "crypto_scam/internal/repository/model"
	"reflect"
	"testing"
)

func TestUpdateUserRequestToUpdateUser(t *testing.T) {
	type args struct {
		from *modelAPI.UpdateUserRequest
	}
	tests := []struct {
		name string
		args args
		want *modelDB.Update
	}{
		{
			name: "Тест №1",
			args: args{
				from: &modelAPI.UpdateUserRequest{
					Blocks: 279127898728,
					Record: 329891291321,
				},
			},
			want: &modelDB.Update{
				Blocks: 279127898728,
				Record: 329891291321,
			},
		},
		{
			name: "Тест №2",
			args: args{
				from: &modelAPI.UpdateUserRequest{
					Blocks: -1,
					Record: -2,
				},
			},
			want: &modelDB.Update{
				Blocks: -1,
				Record: -2,
			},
		},
		{
			name: "Тест №3",
			args: args{
				from: &modelAPI.UpdateUserRequest{
					Blocks: 0,
					Record: 0,
				},
			},
			want: &modelDB.Update{
				Blocks: 0,
				Record: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := converter.UpdateUserRequestToUpdateUser(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUserRequestToUpdateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

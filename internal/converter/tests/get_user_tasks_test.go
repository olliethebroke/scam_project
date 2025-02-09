package tests

import (
	modelAPI "crypto_scam/internal/api/handler/user/model"
	"crypto_scam/internal/converter"
	modelDB "crypto_scam/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"reflect"
	"testing"
)

func TestTasksToGetUserTasksResponse(t *testing.T) {
	type args struct {
		from []*modelDB.Task
	}

	var (
		id          = gofakeit.Int16()
		description = gofakeit.BeerName()
		reward      = gofakeit.Int32()
		actionType  = " \"\\n|\\t\""
		actionData  = gofakeit.MovieName()
		isCompleted = gofakeit.Bool()
	)
	input := make([]*modelDB.Task, 0, 1)

	input = append(input, &modelDB.Task{
		Id:          id,
		Description: description,
		Reward:      reward,
		ActionType:  actionType,
		ActionData:  actionData,
		IsCompleted: isCompleted,
	})

	output := make([]*modelAPI.GetUserTasksResponse, 0, 1)
	output = append(output, &modelAPI.GetUserTasksResponse{
		Id:          id,
		Description: description,
		Reward:      reward,
		ActionType:  actionType,
		ActionData:  actionData,
		IsCompleted: isCompleted,
	})

	tests := []struct {
		name string
		args args
		want []*modelAPI.GetUserTasksResponse
	}{
		{
			name: "Тест №1",
			args: args{
				from: input,
			},
			want: output,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := converter.TasksToGetUserTasksResponse(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TasksToGetUserTasksResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

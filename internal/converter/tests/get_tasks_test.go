package tests

import (
	modelAPI "crypto_scam/internal/api/handler/admin/model"
	"crypto_scam/internal/converter"
	modelDB "crypto_scam/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"reflect"
	"testing"
)

func TestTasksToGetTasksResponse(t *testing.T) {
	type args struct {
		from []*modelDB.Task
	}

	var (
		id          = gofakeit.Int16()
		description = gofakeit.BeerName()
		reward      = gofakeit.Int32()
		actionType  = "watch"
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

	output := make([]*modelAPI.GetTasksResponse, 0, 1)
	output = append(output, &modelAPI.GetTasksResponse{
		Id:          id,
		Description: description,
		Reward:      reward,
		ActionType:  actionType,
		ActionData:  actionData,
	})

	tests := []struct {
		name string
		args args
		want []*modelAPI.GetTasksResponse
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
			if got := converter.TasksToGetTasksResponse(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TasksToGetTasksResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

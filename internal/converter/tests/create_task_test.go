package tests

import (
	modelAPI "crypto_scam/internal/api/handler/admin/model"
	"crypto_scam/internal/converter"
	modelDB "crypto_scam/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"reflect"
	"testing"
)

func TestCreateTaskRequestToTask(t *testing.T) {
	type args struct {
		from *modelAPI.CreateTaskRequest
	}
	tests := []struct {
		name string
		args args
		want *modelDB.Task
	}{
		{
			name: "Тест №1",
			args: args{
				from: &modelAPI.CreateTaskRequest{
					Description: "Subscribe to the truth!",
					Reward:      500,
					ActionType:  "subscribe",
					ActionData:  "https://gay-not-ok.com",
				},
			},
			want: &modelDB.Task{
				Id:          0,
				Description: "Subscribe to the truth!",
				Reward:      500,
				ActionType:  "subscribe",
				ActionData:  "https://gay-not-ok.com",
				IsCompleted: false,
			},
		},
		{
			name: "Тест №2",
			args: args{
				from: &modelAPI.CreateTaskRequest{
					Description: "",
					Reward:      500,
					ActionType:  "\n",
					ActionData:  "💜",
				},
			},
			want: &modelDB.Task{
				Id:          0,
				Description: "",
				Reward:      500,
				ActionType:  "\n",
				ActionData:  "💜",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := converter.CreateTaskRequestToTask(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTaskRequestToTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskToCreateTaskResponse(t *testing.T) {
	type args struct {
		from *modelDB.Task
	}
	var (
		id          = gofakeit.Int16()
		description = gofakeit.Animal()
		reward      = gofakeit.Int32()
		actionType  = gofakeit.Breakfast()
		actionData  = gofakeit.Adjective()
	)
	input := &modelDB.Task{
		Id:          id,
		Description: description,
		Reward:      reward,
		ActionType:  actionType,
		ActionData:  actionData,
	}
	output := &modelAPI.CreateTaskResponse{
		Id: id,
	}
	tests := []struct {
		name string
		args args
		want *modelAPI.CreateTaskResponse
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
			if got := converter.TaskToCreateTaskResponse(tt.args.from); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskToCreateTaskResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

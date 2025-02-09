package tests

import (
	modelDB "crypto_scam/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
)

var (
	adminId          int64 = 1320748911
	adminInitDataRaw       = "query_id=AAFvB7lOAAAAAG8HuU4zoAQY&" +
		"user=%7B%22id%22%3A1320748911%2C%22" +
		"first_name%22%3A%22OLLie%2C%20The%20Broke%22%2C%22" +
		"last_name%22%3A%22%22%2C%22username%22%3A%22oLLieTheBroke%22%2C%22" +
		"language_code%22%3A%22ru%22%2C%22" +
		"allows_write_to_pm%22%3Atrue%2C%22" +
		"photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F38sH16L9KW-Gn-Ct6XS_CFbc3qNa48N5ccQZvY80fFc.svg%22%7D&" +
		"auth_date=1738963365&signature=az82DqUJQdHKXfC3h3zJHyQO0-2SAKq73d2n3B-1H8C8fUhdXNlATlHdg8rMO-cCI2dOWvM8gYSC957snfRlBg&" +
		"hash=e2aa79b98603af01a3e45c1020bf376ee59d69c2b79fdbdbe1cfa760dc8f620a"

	users = []*modelDB.User{
		{
			Id:        902594319,
			Firstname: "S. Morpheus Ulixes amicus est",
		},
		{
			Id:        adminId,
			Firstname: "OLLie, The Broke",
		},
	}
	tasks = []*modelDB.Task{
		{
			Description: "Drink some oslinaya mocha",
			Reward:      5000,
			ActionType:  "drink",
			ActionData:  gofakeit.BeerName(),
		},
		{
			Description: "Kill an animal in Minecraft",
			Reward:      5000,
			ActionType:  "kill",
			ActionData:  gofakeit.MinecraftAnimal(),
		},
		{
			Description: "Scam morpheus",
			Reward:      10000,
			ActionType:  "scam",
			ActionData:  gofakeit.MovieName(),
		},
	}
	completedTasks = map[int64][]int16{
		adminId: {
			2,
			1,
		},
	}
	leaders = map[int16][]*modelDB.Leader{
		3: {
			{
				Id:        gofakeit.Int64(),
				Firstname: gofakeit.BeerName(),
				Position:  gofakeit.Int16(),
				Blocks:    gofakeit.Int64(),
			},
			{
				Id:        gofakeit.Int64(),
				Firstname: gofakeit.BeerName(),
				Position:  gofakeit.Int16(),
				Blocks:    gofakeit.Int64(),
			},
			{
				Id:        gofakeit.Int64(),
				Firstname: gofakeit.BeerName(),
				Position:  gofakeit.Int16(),
				Blocks:    gofakeit.Int64(),
			},
			{
				Id:        gofakeit.Int64(),
				Firstname: gofakeit.BeerName(),
				Position:  gofakeit.Int16(),
				Blocks:    gofakeit.Int64(),
			},
		},
	}
)

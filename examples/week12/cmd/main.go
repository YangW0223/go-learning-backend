package main

import (
	"fmt"

	"github.com/yang/go-learning-backend/examples/week12"
)

func main() {
	portfolio := week12.NewPortfolio()
	portfolio.MarkCompleted("http", "week02 delete api")
	portfolio.MarkCompleted("concurrency", "week03 timeout + chatroom")
	portfolio.MarkCompleted("db", "week05 tx + pagination")
	portfolio.MarkCompleted("auth", "week07 jwt middleware")

	caps := portfolio.Capabilities()
	score := portfolio.Score(6)
	narrative := week12.BuildNarrative("todo-service", caps)

	fmt.Println("score:", score)
	fmt.Println("narrative:", narrative)
	fmt.Println("next stage (mq):", week12.NextStagePlan("mq"))
}

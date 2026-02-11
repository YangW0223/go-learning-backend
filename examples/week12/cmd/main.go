// 详细注释: package main
package main

// 详细注释: import (
import (
	// 详细注释: "fmt"
	"fmt"

	// 详细注释: "github.com/yang/go-learning-backend/examples/week12"
	"github.com/yang/go-learning-backend/examples/week12"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: portfolio := week12.NewPortfolio()
	portfolio := week12.NewPortfolio()
	// 详细注释: portfolio.MarkCompleted("http", "week02 delete api")
	portfolio.MarkCompleted("http", "week02 delete api")
	// 详细注释: portfolio.MarkCompleted("concurrency", "week03 timeout + chatroom")
	portfolio.MarkCompleted("concurrency", "week03 timeout + chatroom")
	// 详细注释: portfolio.MarkCompleted("db", "week05 tx + pagination")
	portfolio.MarkCompleted("db", "week05 tx + pagination")
	// 详细注释: portfolio.MarkCompleted("auth", "week07 jwt middleware")
	portfolio.MarkCompleted("auth", "week07 jwt middleware")

	// 详细注释: caps := portfolio.Capabilities()
	caps := portfolio.Capabilities()
	// 详细注释: score := portfolio.Score(6)
	score := portfolio.Score(6)
	// 详细注释: narrative := week12.BuildNarrative("todo-service", caps)
	narrative := week12.BuildNarrative("todo-service", caps)

	// 详细注释: fmt.Println("score:", score)
	fmt.Println("score:", score)
	// 详细注释: fmt.Println("narrative:", narrative)
	fmt.Println("narrative:", narrative)
	// 详细注释: fmt.Println("next stage (mq):", week12.NextStagePlan("mq"))
	fmt.Println("next stage (mq):", week12.NextStagePlan("mq"))
	// 详细注释: }
}

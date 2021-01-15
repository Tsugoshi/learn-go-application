package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type CLI struct {
	game Game
	in   *bufio.Scanner
	out  io.Writer
}

func NewCLI(game Game, in io.Reader, out io.Writer) *CLI {
	return &CLI{
		game: game,
		in:   bufio.NewScanner(in),
		out:  out,
	}
}

const PlayerGreeting = "Please enter the number of players: "
const BadStartInput = "Expected number of players"

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerGreeting)

	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, err := strconv.Atoi(numberOfPlayersInput)

	if err != nil {
		fmt.Fprint(cli.out, BadStartInput)
		return
	}

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)
	cli.game.Finish(winner)
}

func extractWinner(line string) string {
	return strings.Replace(line, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

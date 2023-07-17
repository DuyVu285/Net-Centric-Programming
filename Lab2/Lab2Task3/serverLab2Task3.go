package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

type HangmanGame struct {
	word          string
	description   string
	guessed       []bool
	players       []*Player
	currentPlayer int
}

type Player struct {
	conn    net.Conn
	score   int
	hasTurn bool
	guess   byte
	name    string
}

func readWordsFromFile(filename string) ([]HangmanGame, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	games := []HangmanGame{}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		game := HangmanGame{
			word:          parts[0],
			description:   parts[1],
			guessed:       make([]bool, len(parts[0])),
			players:       []*Player{},
			currentPlayer: 0,
		}
		games = append(games, game)
	}

	return games, nil
}

func getRandomGame(games []HangmanGame) HangmanGame {
	rand.Seed(time.Now().Unix())
	idx := rand.Intn(len(games))
	return games[idx]
}

func startGame(players []*Player) {
	games, err := readWordsFromFile("words.txt")
	if err != nil {
		log.Fatal(err)
	}

	game := getRandomGame(games)

	for _, player := range players {
		player.conn.Write([]byte("Starting new game...\n"))
		player.conn.Write([]byte(fmt.Sprintf("Description: %s\n", game.description)))
		player.conn.Write([]byte(fmt.Sprintf("Word: %s\n", getGuessedWord(game.word, game.guessed))))
	}

	for {
		player := players[game.currentPlayer]

		if !player.hasTurn {
			continue
		}

		player.conn.Write([]byte("Your turn to guess a letter: "))

		player.conn.SetDeadline(time.Now().Add(30 * time.Second))

		buf := make([]byte, 1024)
		n, err := player.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Println(err)
			return
		}

		guess := buf[:n][0]

		correct := false
		for i := 0; i < len(game.word); i++ {
			if game.word[i] == guess {
				game.guessed[i] = true
				correct = true
			}
		}

		score := 0
		if correct {
			score = countOccurrences(game.word, guess) * 10
			player.score += score
		}

		for _, p := range game.players {
			if p == player {
				continue
			}
			if correct {
				p.conn.Write([]byte(fmt.Sprintf("%s guessed %c and earned %d points\n", player.name, guess, score)))
			} else {
				p.conn.Write([]byte(fmt.Sprintf("%s guessed %c and failed\n", player.name, guess)))
			}
			p.conn.Write([]byte(fmt.Sprintf("Word: %s\n", getGuessedWord(game.word, game.guessed))))
		}

		if isGameOver(game.guessed) {
			winner := getWinner(game.players)
			for _, p := range game.players {
				p.conn.Write([]byte(fmt.Sprintf("Game over! %s has won with a score of %d\n", winner.name, winner.score)))
			}
			return
		}

		game.currentPlayer = (game.currentPlayer + 1) % len(game.players)
	}
}

func getGuessedWord(word string, guessed []bool) string {
	var guessedWord string

	for i := 0; i < len(word); i++ {
		if guessed[i] {
			guessedWord += string(word[i])
		} else {
			guessedWord += "_"
		}
	}

	return guessedWord
}

func countOccurrences(word string, letter byte) int {
	count := 0
	for i := 0; i < len(word); i++ {
		if word[i] == letter {
			count++
		}
	}
	return count
}

func isGameOver(guessed []bool) bool {
	for _, g := range guessed {
		if !g {
			return false
		}
	}
	return true
}

func getWinner(players []*Player) *Player {
	maxScore := -1
	var winner *Player
	for _, p := range players {
		if p.score > maxScore {
			maxScore = p.score
			winner = p
		}
	}
	return winner
}
func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Server started, listening on port 8080")

	var players []*Player
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		player := &Player{
			conn:    conn,
			score:   0,
			hasTurn: false,
			guess:   0,
			name:    "",
		}
		players = append(players, player)

		if len(players) >= 2 {
			for i := range players {
				players[i].hasTurn = true
				if i == 0 {
					players[i].name = "Player 1"
				} else if i == 1 {
					players[i].name = "Player 2"
				} else {
					players[i].name = fmt.Sprintf("Player %d", i+1)
				}
			}

			startGame(players)

			for _, player := range players {
				player.hasTurn = false
				player.score = 0
				player.name = ""
			}
			players = nil
		}
	}
}

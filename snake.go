// snake.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type Point struct{ r, c int }

type SnakeGame struct {
	rows, cols int
	maze       [][]byte
	snake      []Point
	dir        Point
	food       Point
	score      int
	speed      time.Duration
	gameOver   bool
	running    bool
}

func NewSnakeGame(rows, cols int) *SnakeGame {
	g := &SnakeGame{
		rows:    rows,
		cols:    cols,
		speed:   150 * time.Millisecond,
		running: true,
	}
	g.generateMaze()
	g.spawnSnake()
	g.spawnFood()
	return g
}

func (g *SnakeGame) generateMaze() {
	g.maze = make([][]byte, g.rows)
	for i := range g.maze {
		g.maze[i] = make([]byte, g.cols)
		for j := range g.maze[i] {
			g.maze[i][j] = '#'
		}
	}
	// DFS carve
	stack := []Point{{1, 1}}
	g.maze[1][1] = ' '
	dirs := [][2]int{{-2, 0}, {2, 0}, {0, -2}, {0, 2}}
	for len(stack) > 0 {
		r, c := stack[len(stack)-1].r, stack[len(stack)-1].c
		var neighbours []int
		for i, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr > 0 && nr < g.rows-1 && nc > 0 && nc < g.cols-1 && g.maze[nr][nc] == '#' {
				neighbours = append(neighbours, i)
			}
		}
		if len(neighbours) > 0 {
			idx := neighbours[rand.Intn(len(neighbours))]
			dr, dc := dirs[idx][0], dirs[idx][1]
			nr, nc := r+dr, c+dc
			g.maze[r+dr/2][c+dc/2] = ' '
			g.maze[nr][nc] = ' '
			stack = append(stack, Point{nr, nc})
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	// borders
	for i := 0; i < g.rows; i++ {
		g.maze[i][0] = '#'
		g.maze[i][g.cols-1] = '#'
	}
	for j := 0; j < g.cols; j++ {
		g.maze[0][j] = '#'
		g.maze[g.rows-1][j] = '#'
	}
}

func (g *SnakeGame) spawnSnake() {
	center := Point{g.rows / 2, g.cols / 2}
	if g.maze[center.r][center.c] != ' ' {
		for r := 1; r < g.rows-1; r++ {
			for c := 1; c < g.cols-1; c++ {
				if g.maze[r][c] == ' ' {
					center = Point{r, c}
					break
				}
			}
			if center.r == g.rows/2 && center.c == g.cols/2 {
				// if not found, break outer
				break
			}
		}
	}
	g.snake = []Point{center}
	g.dir = Point{0, 1}
}

func (g *SnakeGame) spawnFood() {
	var empty []Point
	for r := 1; r < g.rows-1; r++ {
		for c := 1; c < g.cols-1; c++ {
			if g.maze[r][c] == ' ' {
				occupied := false
				for _, p := range g.snake {
					if p.r == r && p.c == c {
						occupied = true
						break
					}
				}
				if !occupied {
					empty = append(empty, Point{r, c})
				}
			}
		}
	}
	if len(empty) > 0 {
		g.food = empty[rand.Intn(len(empty))]
	} else {
		g.gameOver = true
	}
}

func (g *SnakeGame) update() {
	if g.gameOver {
		return
	}
	head := g.snake[0]
	newHead := Point{head.r + g.dir.r, head.c + g.dir.c}
	// wall collision
	if newHead.r <= 0 || newHead.r >= g.rows-1 || newHead.c <= 0 || newHead.c >= g.cols-1 || g.maze[newHead.r][newHead.c] == '#' {
		g.gameOver = true
		return
	}
	// self collision
	for _, p := range g.snake {
		if p == newHead {
			g.gameOver = true
			return
		}
	}
	// move
	g.snake = append([]Point{newHead}, g.snake...)
	if newHead == g.food {
		g.score++
		g.speed = time.Duration(float64(g.speed) * 0.95)
		if g.speed < 50*time.Millisecond {
			g.speed = 50 * time.Millisecond
		}
		g.spawnFood()
		if g.gameOver {
			return
		}
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}
}

func (g *SnakeGame) draw() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Printf("Score: %d  Speed: %d\n", g.score, int(1000/int(g.speed.Milliseconds())))
	// Build display
	display := make([][]byte, g.rows)
	for i := range display {
		display[i] = make([]byte, g.cols)
		copy(display[i], g.maze[i])
	}
	for _, p := range g.snake {
		display[p.r][p.c] = 'O'
	}
	if !g.gameOver && g.food != (Point{}) {
		display[g.food.r][g.food.c] = '*'
	}
	display[g.snake[0].r][g.snake[0].c] = 'X'
	for _, row := range display {
		fmt.Println(string(row))
	}
	if g.gameOver {
		fmt.Println("GAME OVER! Press R to restart, Q to quit.")
	}
}

func (g *SnakeGame) handleInput() {
	// We'll use a simple non-blocking read using a goroutine and channels.
	// For simplicity, we'll use `fmt.Scanf`? In Go, we can use `bufio` and a goroutine.
	// We'll implement a basic version using `os.Stdin` with `bufio` and a separate goroutine.
	// This is a simplified version; for production, consider `tcell` or `termbox`.
	// We'll use a channel to send key events.
	keyChan := make(chan string)
	go func() {
		for g.running {
			var b [1]byte
			_, err := os.Stdin.Read(b[:])
			if err != nil {
				continue
			}
			ch := b[0]
			if ch == 27 { // escape
				keyChan <- "esc"
				continue
			}
			switch ch {
			case 'q', 'Q':
				keyChan <- "quit"
			case 'r', 'R':
				keyChan <- "restart"
			case 'w', 'W':
				keyChan <- "up"
			case 's', 'S':
				keyChan <- "down"
			case 'a', 'A':
				keyChan <- "left"
			case 'd', 'D':
				keyChan <- "right"
			}
		}
	}()
	for g.running {
		select {
		case key := <-keyChan:
			if key == "quit" {
				g.running = false
				return
			}
			if key == "restart" && g.gameOver {
				g.__init__(g.rows, g.cols)
				return
			}
			if !g.gameOver {
				switch key {
				case "up":
					if g.dir != (Point{1, 0}) {
						g.dir = Point{-1, 0}
					}
				case "down":
					if g.dir != (Point{-1, 0}) {
						g.dir = Point{1, 0}
					}
				case "left":
					if g.dir != (Point{0, 1}) {
						g.dir = Point{0, -1}
					}
				case "right":
					if g.dir != (Point{0, -1}) {
						g.dir = Point{0, 1}
					}
				}
			}
		default:
			// no key
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (g *SnakeGame) __init__(rows, cols int) {
	*g = *NewSnakeGame(rows, cols)
}

func (g *SnakeGame) run() {
	go g.handleInput()
	for g.running {
		g.update()
		g.draw()
		time.Sleep(g.speed)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game := NewSnakeGame(21, 21)
	game.run()
}

# snake.py
import random
import os
import time
import sys
from collections import deque

class MazeSnake:
    def __init__(self, rows=21, cols=21):
        self.rows = rows
        self.cols = cols
        self.maze = []
        self.snake = deque()
        self.direction = (0, 1)  # (dr, dc)
        self.food = None
        self.score = 0
        self.speed = 0.15
        self.game_over = False
        self.running = True
        self.generate_maze()
        self.spawn_snake()
        self.spawn_food()

    def generate_maze(self):
        # Initialize grid with walls
        self.maze = [['#'] * self.cols for _ in range(self.rows)]
        # DFS to carve paths (start at (1,1))
        stack = [(1, 1)]
        self.maze[1][1] = ' '
        while stack:
            r, c = stack[-1]
            # Get unvisited neighbours (2 cells away)
            neighbours = []
            for dr, dc in [(-2,0), (2,0), (0,-2), (0,2)]:
                nr, nc = r + dr, c + dc
                if 0 < nr < self.rows-1 and 0 < nc < self.cols-1 and self.maze[nr][nc] == '#':
                    neighbours.append((nr, nc, dr, dc))
            if neighbours:
                nr, nc, dr, dc = random.choice(neighbours)
                # Remove wall between
                self.maze[r + dr//2][c + dc//2] = ' '
                self.maze[nr][nc] = ' '
                stack.append((nr, nc))
            else:
                stack.pop()
        # Ensure border walls remain
        for r in range(self.rows):
            self.maze[r][0] = '#'
            self.maze[r][self.cols-1] = '#'
        for c in range(self.cols):
            self.maze[0][c] = '#'
            self.maze[self.rows-1][c] = '#'

    def spawn_snake(self):
        # Place snake in the centre (or near centre)
        center_r = self.rows // 2
        center_c = self.cols // 2
        if self.maze[center_r][center_c] != ' ':
            # find first empty cell
            for r in range(1, self.rows-1):
                for c in range(1, self.cols-1):
                    if self.maze[r][c] == ' ':
                        center_r, center_c = r, c
                        break
                else:
                    continue
                break
        self.snake = deque([(center_r, center_c)])
        self.direction = (0, 1)

    def spawn_food(self):
        # find empty cell not occupied by snake
        empty = []
        for r in range(1, self.rows-1):
            for c in range(1, self.cols-1):
                if self.maze[r][c] == ' ' and (r, c) not in self.snake:
                    empty.append((r, c))
        if empty:
            self.food = random.choice(empty)
        else:
            self.game_over = True

    def update(self):
        if self.game_over:
            return
        # Move snake
        dr, dc = self.direction
        head_r, head_c = self.snake[0]
        new_r, new_c = head_r + dr, head_c + dc
        # Check collision with wall
        if not (0 < new_r < self.rows-1 and 0 < new_c < self.cols-1) or self.maze[new_r][new_c] == '#':
            self.game_over = True
            return
        # Check collision with self
        if (new_r, new_c) in self.snake:
            self.game_over = True
            return
        # Move
        self.snake.appendleft((new_r, new_c))
        # Check food
        if (new_r, new_c) == self.food:
            self.score += 1
            self.speed = max(0.05, self.speed - 0.005)
            self.spawn_food()
            # If no food, game over
            if self.game_over:
                return
        else:
            self.snake.pop()

    def draw(self):
        os.system('cls' if os.name == 'nt' else 'clear')
        print(f"Score: {self.score}  Speed: {int(1/self.speed)}")
        # Build display grid
        display = [list(row) for row in self.maze]
        for r, c in self.snake:
            display[r][c] = 'O'
        if self.food:
            display[self.food[0]][self.food[1]] = '*'
        display[self.snake[0][0]][self.snake[0][1]] = 'X'
        # Print
        for row in display:
            print(''.join(row))
        if self.game_over:
            print("GAME OVER! Press R to restart, Q to quit.")

    def handle_input(self):
        import threading
        def get_key():
            try:
                import termios, tty
                fd = sys.stdin.fileno()
                old = termios.tcgetattr(fd)
                try:
                    tty.setraw(fd)
                    ch = sys.stdin.read(1)
                    if ch == '\x1b':
                        ch2 = sys.stdin.read(1)
                        if ch2 == '[':
                            ch3 = sys.stdin.read(1)
                            if ch3 == 'A': return 'up'
                            elif ch3 == 'B': return 'down'
                            elif ch3 == 'C': return 'right'
                            elif ch3 == 'D': return 'left'
                    return ch.lower()
                finally:
                    termios.tcsetattr(fd, termios.TCSADRAIN, old)
            except:
                return None
        while self.running:
            key = get_key()
            if key == 'q':
                self.running = False
                break
            if key == 'r' and self.game_over:
                self.__init__(self.rows, self.cols)
                break
            if not self.game_over:
                if key == 'up' and self.direction != (1, 0):
                    self.direction = (-1, 0)
                elif key == 'down' and self.direction != (-1, 0):
                    self.direction = (1, 0)
                elif key == 'left' and self.direction != (0, 1):
                    self.direction = (0, -1)
                elif key == 'right' and self.direction != (0, -1):
                    self.direction = (0, 1)
            time.sleep(0.05)

    def run(self):
        # Start input thread
        import threading
        input_thread = threading.Thread(target=self.handle_input, daemon=True)
        input_thread.start()
        while self.running:
            self.update()
            self.draw()
            time.sleep(self.speed)
        print("Goodbye!")

if __name__ == "__main__":
    game = MazeSnake(21, 21)
    game.run()

// Snake.java
import java.io.*;
import java.util.*;
import java.util.concurrent.*;

public class Snake {
    static final int ROWS = 21, COLS = 21;
    char[][] maze;
    List<int[]> snake = new ArrayList<>();
    int[] dir = {0, 1};
    int[] food;
    int score = 0;
    int speed = 150; // ms
    boolean gameOver = false;
    boolean running = true;
    Random rand = new Random();

    public Snake() {
        generateMaze();
        spawnSnake();
        spawnFood();
        new Thread(this::handleInput).start();
        gameLoop();
    }

    void generateMaze() {
        maze = new char[ROWS][COLS];
        for (int i = 0; i < ROWS; i++) Arrays.fill(maze[i], '#');
        Stack<int[]> stack = new Stack<>();
        stack.push(new int[]{1, 1});
        maze[1][1] = ' ';
        int[][] dirs = {{-2,0},{2,0},{0,-2},{0,2}};
        while (!stack.isEmpty()) {
            int[] cur = stack.peek();
            List<Integer> neigh = new ArrayList<>();
            for (int i=0; i<4; i++) {
                int nr = cur[0]+dirs[i][0], nc = cur[1]+dirs[i][1];
                if (nr>0 && nr<ROWS-1 && nc>0 && nc<COLS-1 && maze[nr][nc]=='#')
                    neigh.add(i);
            }
            if (!neigh.isEmpty()) {
                int idx = neigh.get(rand.nextInt(neigh.size()));
                int dr = dirs[idx][0], dc = dirs[idx][1];
                int nr = cur[0]+dr, nc = cur[1]+dc;
                maze[cur[0]+dr/2][cur[1]+dc/2] = ' ';
                maze[nr][nc] = ' ';
                stack.push(new int[]{nr, nc});
            } else {
                stack.pop();
            }
        }
        for (int i=0; i<ROWS; i++) { maze[i][0]='#'; maze[i][COLS-1]='#'; }
        for (int j=0; j<COLS; j++) { maze[0][j]='#'; maze[ROWS-1][j]='#'; }
    }

    void spawnSnake() {
        int cr=ROWS/2, cc=COLS/2;
        if (maze[cr][cc] != ' ') {
            outer: for (int r=1; r<ROWS-1; r++)
                for (int c=1; c<COLS-1; c++)
                    if (maze[r][c]==' ') { cr=r; cc=c; break outer; }
        }
        snake.clear();
        snake.add(new int[]{cr, cc});
        dir = new int[]{0, 1};
    }

    void spawnFood() {
        List<int[]> empty = new ArrayList<>();
        for (int r=1; r<ROWS-1; r++)
            for (int c=1; c<COLS-1; c++) {
                if (maze[r][c]==' ') {
                    boolean occ = false;
                    for (int[] p : snake) if (p[0]==r && p[1]==c) { occ=true; break; }
                    if (!occ) empty.add(new int[]{r, c});
                }
            }
        if (!empty.isEmpty()) food = empty.get(rand.nextInt(empty.size()));
        else gameOver = true;
    }

    void update() {
        if (gameOver) return;
        int[] head = snake.get(0);
        int nr = head[0] + dir[0], nc = head[1] + dir[1];
        if (nr<=0 || nr>=ROWS-1 || nc<=0 || nc>=COLS-1 || maze[nr][nc]=='#') {
            gameOver = true; return;
        }
        for (int[] p : snake) if (p[0]==nr && p[1]==nc) { gameOver=true; return; }
        snake.add(0, new int[]{nr, nc});
        if (nr==food[0] && nc==food[1]) {
            score++;
            speed = (int)Math.max(50, speed*0.95);
            spawnFood();
            if (gameOver) return;
        } else {
            snake.remove(snake.size()-1);
        }
    }

    void draw() {
        System.out.print("\033[H\033[2J");
        System.out.printf("Score: %d  Speed: %d\n", score, 1000/speed);
        char[][] disp = new char[ROWS][COLS];
        for (int i=0; i<ROWS; i++) System.arraycopy(maze[i], 0, disp[i], 0, COLS);
        for (int[] p : snake) disp[p[0]][p[1]] = 'O';
        if (!gameOver) disp[food[0]][food[1]] = '*';
        disp[snake.get(0)[0]][snake.get(0)[1]] = 'X';
        for (char[] row : disp) System.out.println(new String(row));
        if (gameOver) System.out.println("GAME OVER! Press R to restart, Q to quit.");
    }

    void handleInput() {
        try {
            while (running) {
                int ch = System.in.read();
                if (ch == -1) continue;
                char c = (char) ch;
                if (c == 'q' || c == 'Q') { running = false; System.exit(0); }
                if (c == 'r' || c == 'R') {
                    if (gameOver) { restart(); return; }
                }
                if (!gameOver) {
                    if (c == 'w' && dir[0] != 1) dir = new int[]{-1, 0};
                    else if (c == 's' && dir[0] != -1) dir = new int[]{1, 0};
                    else if (c == 'a' && dir[1] != 1) dir = new int[]{0, -1};
                    else if (c == 'd' && dir[1] != -1) dir = new int[]{0, 1};
                }
            }
        } catch (IOException e) {}
    }

    void restart() {
        // Reinitialize the game state
        generateMaze();
        spawnSnake();
        spawnFood();
        score = 0;
        speed = 150;
        gameOver = false;
    }

    void gameLoop() {
        while (running) {
            update();
            draw();
            try { Thread.sleep(speed); } catch (InterruptedException e) {}
        }
    }

    public static void main(String[] args) {
        new Snake();
    }
}

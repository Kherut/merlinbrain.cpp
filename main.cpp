#include <ncurses.h>
#include <string>
#include <cstring>
#include <vector>

using std::vector;

WINDOW *create_newwin(int height, int width, int starty, int startx);
void destroy_win(WINDOW *local_win);
void execute(const char *command, std::vector<std::string> queue);
void addQueue(int index, std::vector<std::string> queue, WINDOW *queue_win, const char* command);

int main(int argc, char *argv[])
{	
	WINDOW *devices_win, *command_win, *queue_win;
	int x, y;
	char command[256];
	std::vector<std::string> queue;
	int index = 1;

	initscr();
	cbreak();

	keypad(stdscr, TRUE);

	attron(A_BOLD);
	mvprintw(0, COLS / 2 - 6, "Merlin Brain");
	attroff(A_BOLD);

    refresh();

    y = 1;
	x = 1;
	devices_win = create_newwin(LINES - 1, COLS / 3 - 1, y, x);
    mvwprintw(devices_win, 0, 1, "Devices");

    y = 1;
    x = COLS / 3 + 1;
    command_win = create_newwin(3, COLS * 2 / 3 - 1, y, x);
	mvwprintw(command_win, 0, 1, "Command");
	mvwprintw(command_win, 1, 1, " > ");

	y = 4;
    x = COLS / 3 + 1;
    queue_win = create_newwin(LINES - 4, COLS * 2 / 3 - 1, y, x);
	mvwprintw(queue_win, 0, 1, "Queue");

	wrefresh(devices_win);
	wrefresh(command_win);
	wrefresh(queue_win);

	while (true) {
        mvwgetstr(command_win, 1, 4, command);
		
		if(strcmp(command, "force_exit") == 0) {
			endwin();
			return 0;
		}

		addQueue(index, queue, queue_win, command);
		index++;

		for(int i = 4; i < COLS * 2 / 3 - 2; i++)
			mvwprintw(command_win, 1, i, " ");
	}
		
	endwin();
	return 0;
}

WINDOW *create_newwin(int height, int width, int starty, int startx) {
	WINDOW *local_win;

	local_win = newwin(height, width, starty, startx);
	box(local_win, 0 , 0);
	wrefresh(local_win);

	return local_win;
}

void destroy_win(WINDOW *local_win) {	
	wborder(local_win, ' ', ' ', ' ',' ',' ',' ',' ',' ');
	
	wrefresh(local_win);
	delwin(local_win);
}

void addQueue(int index, std::vector<std::string> queue, WINDOW *queue_win, const char* command) {
	mvwprintw(queue_win, index, 1, command);

	wrefresh(queue_win);
}
#include <ncurses.h>
#include <string>
#include <sstream>
#include <cstring>
#include <vector>
#include <iterator>

using std::vector;

template<typename Out>
void split(const std::string &s, char delim, Out result) {
    std::stringstream ss(s);
    std::string item;
    while (std::getline(ss, item, delim)) {
        *(result++) = item;
    }
}

std::vector<std::string> split(const std::string &s, char delim) {
    std::vector<std::string> elems;
    split(s, delim, std::back_inserter(elems));
    return elems;
}

WINDOW *create_newwin(int height, int width, int starty, int startx);
void destroy_win(WINDOW *local_win);
void execute(const char *command);
void printWindow(WINDOW *win, std::vector<std::string> list);
void clear(WINDOW *win, int width, int height);

int main(int argc, char *argv[])
{	
	WINDOW *devices_win, *command_win, *queue_win;
	int x, y;
	int dW, dH;
	int cW, cH;
	int qW, qH;
	char command[256];
	std::vector<std::string> queue, devices;

	initscr();
	cbreak();

	keypad(stdscr, TRUE);

	attron(A_BOLD);
	mvprintw(0, COLS / 2 - 6, "Merlin Brain");
	attroff(A_BOLD);

    refresh();

	dW = COLS / 3 - 1;
	dH = LINES - 1;
    y = 1;
	x = 1;
	devices_win = create_newwin(dH, dW, y, x);
    mvwprintw(devices_win, 0, 1, "Devices");

	cW = COLS * 2 / 3 - 1;
	cH = 3;
    y = 1;
    x = COLS / 3 + 1;
    command_win = create_newwin(cH, cW, y, x);
	mvwprintw(command_win, 0, 1, "Command");
	mvwprintw(command_win, 1, 1, " > ");

	qW = COLS * 2 / 3 - 1;
	qH = LINES - 4;
	y = 4;
    x = COLS / 3 + 1;
    queue_win = create_newwin(qH, qW, y, x);
	mvwprintw(queue_win, 0, 1, "Queue");

	wrefresh(devices_win);
	wrefresh(command_win);
	wrefresh(queue_win);

	std::vector<std::string> args;

	while (true) {
        mvwgetstr(command_win, 1, 4, command);
		
		if(strcmp(command, "forceexit") == 0) {
			endwin();
			return 0;
		}

		queue.push_back(std::string(command));

		args = split(std::string(command), ' ');

		//EXECUTE
		if(strcmp(args.at(0).c_str(), "clear") == 0) {
			if(args.size() > 1) {
				if(strcmp(args.at(1).c_str(), "queue") == 0) {
					clear(queue_win, qW, qH);
					queue.clear();
				}
			}
		}

		else if(strcmp(args.at(0).c_str(), "pop") == 0) {
			if(args.size() > 1) {
				if(strcmp(args.at(1).c_str(), "front") == 0) {
					queue.erase(queue.begin());
					queue.erase(queue.end());
				}

				else if(strcmp(args.at(1).c_str(), "back") == 0) {
					queue.erase(queue.end());
					queue.erase(queue.end());
				}
			}
		}

		clear(queue_win, qW, qH);
		printWindow(queue_win, queue);

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

void printWindow(WINDOW *win, std::vector<std::string> list) {
	int position = 1;

	for (auto element = list.begin(); element != list.end(); ++element) {
		mvwprintw(win, position, 1, (*element).c_str());

		position++;
	}

	wrefresh(win);
}

void clear(WINDOW *win, int width, int height) {
	for(int i = 1; i < height - 1; i++)
		mvwprintw(win, i, 1, std::string(width - 2, ' ').c_str());
	
	wrefresh(win);
}
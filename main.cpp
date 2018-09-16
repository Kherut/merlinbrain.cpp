#include <ncurses.h>
#include <string>

WINDOW *create_newwin(int height, int width, int starty, int startx);
void destroy_win(WINDOW *local_win);

int main(int argc, char *argv[])
{	WINDOW *devices_win, *command_win, *queue_win;
	int x, y;
	int ch;

	initscr();
	cbreak();

	keypad(stdscr, TRUE);

    /*printw(const_cast<char*>(std::to_string(LINES).c_str()));
    printw("\n");
    printw(const_cast<char*>(std::to_string(COLS).c_str()));*/

	refresh();

    //y = 0;
	//x = 0;
	//devices_win = create_newwin(LINES, COLS / 3, y, x);
    
    y = COLS / 3;
    x = 0;
    command_win = create_newwin(4, COLS * 2 / 3, y, x);

	while((ch = getch()) != KEY_F(1))
	{	
        
	}
		
	endwin();
	return 0;
}

WINDOW *create_newwin(int height, int width, int starty, int startx)
{	WINDOW *local_win;

	local_win = newwin(height, width, starty, startx);
	box(local_win, 0 , 0);		/* 0, 0 gives default characters 
					 * for the vertical and horizontal
					 * lines			*/
	wrefresh(local_win);		/* Show that box 		*/

	return local_win;
}

void destroy_win(WINDOW *local_win)
{	
	wborder(local_win, ' ', ' ', ' ',' ',' ',' ',' ',' ');
	
	wrefresh(local_win);
	delwin(local_win);
}

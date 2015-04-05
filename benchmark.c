#include <stdlib.h>
#include <stdio.h>

#define colon(i, l, s) character(':', i, l, s)
#define point(i, l, s) character('.', i, l, s)

inline bool digit(int& i, char* line, const size_t size) 
{
	char c = line[i];
	if (c <= '0' || c >= '9')
		return false;
	if (i > size)
		return false;
	i++;
	return true;
}

inline bool character(char which, int& i, char* line, const size_t size)
{
	char c = line[i];
	if (c != which)
		return false;
	if (i > size)
		return false;
	i++;
	return true;	
}

int main()
{
	char* line = NULL;
	size_t size;
	while (getline(&line, &size, stdin) != -1) {
		bool hit = false;
		for (int i = 0; i < size; i++) {
			// look for hh:mm:ss.mmm	
			if (digit(i, line, size) &&
				digit(i, line, size) &&
				colon(i, line, size) &&
				digit(i, line, size) &&
				digit(i, line, size) &&
				colon(i, line, size) &&
				digit(i, line, size) &&
				digit(i, line, size) &&
				point(i, line, size) &&
				digit(i, line, size) &&
				digit(i, line, size) &&
				digit(i, line, size)) {
				hit = true;
				break;	
			} else {
				continue;
			}
		}
		if ( hit == true )
			printf("*");
		else 
			printf(" ");
		printf("%s", line);
	}
	return 0;
}


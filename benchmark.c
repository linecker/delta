#include <stdlib.h>
#include <stdio.h>

inline bool digit(int& i, char* line, const size_t size) 
{
	char c = line[i];
	if (c <= '0' || c >= '9')
		return false;
	if (i >= size)
		return false;
	i++;
	return true;
}

inline bool colon(int& i, char* line, const size_t size) 
{
	char c = line[i];
	if (c != ':')
		return false;
	if (i >= size)
		return false;
	i++;
	return true;
}

inline bool point(int& i, char* line, const size_t size) 
{
	char c = line[i];
	if (c != '.')
		return false;
	if (i >= size)
		return false;
	i++;
	return true;
}

int main()
{
	char* line = NULL;
	size_t size;
	while (getline(&line, &size, stdin) != -1) {
		printf("%s", line);
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
				printf("timestamp found!");	
			} else {
				continue;
			}
		}
	}
	return 0;
}


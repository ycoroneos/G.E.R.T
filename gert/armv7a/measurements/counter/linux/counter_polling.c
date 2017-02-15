#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <err.h>
#include <unistd.h>

int main()
{
  int exporter = open("/sys/class/gpio/export", O_WRONLY);
  if (exporter<0) {
    printf("fail to open exporter\n");
    return -1;
  }
  write(exporter, "91", 2);
  close(exporter);

  int direction = open("/sys/class/gpio/gpio91/direction", O_WRONLY);
  if (direction<0) {
    printf("fail to open direction\n");
    return -1;
  }
  write(direction, "in", 2);
  close(direction);

  int value = open("/sys/class/gpio/gpio91/value", O_RDONLY);
  if (value<0) {
    printf("fail to open value\n");
    return -1;
  }
  
  int i=0;
  int count=0;
  char oldc='0';
  char c='0';
  for (i=0;i<10000000/2;++i)
  {
    read(value, &c, 1);
    lseek(value, 0, SEEK_SET);
    if (c==oldc)
     {
       continue;
     }
    else
     {
       oldc=c;
       ++count;
       //printf("read: %d, count is %d\n", oldc, count);
     }	
  }
printf("count is %d\n", count);
  close(value);
  return 0;
}

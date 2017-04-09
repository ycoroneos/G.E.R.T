#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <err.h>
#include <unistd.h>

int main()
{
  //make 91 an input
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


  //make 191 the output
  exporter = open("/sys/class/gpio/export", O_WRONLY);
  if (exporter<0) {
    printf("fail to open exporter\n");
    return -1;
  }
  write(exporter, "191", 3);
  close(exporter);

  direction = open("/sys/class/gpio/gpio191/direction", O_WRONLY);
  if (direction<0) {
    printf("fail to open direction\n");
    return -1;
  }
  write(direction, "out", 3);
  close(direction);

  int writeval = open("/sys/class/gpio/gpio191/value", O_WRONLY);
  if (writeval<0) {
    printf("fail to open value\n");
    return -1;
  }
  char zero='0';
  char one='1';
  write(writeval, &zero, 1);

  ///wait for an input
  int i=0;
  int count=0;
  char oldc='1';
  char c='1';
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
       write(writeval, &one, 1);
       oldc=c;
       //printf("read: %d, count is %d\n", oldc, count);
     }
  }
  close(value);
  return 0;
}

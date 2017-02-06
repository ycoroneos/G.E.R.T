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
  write(direction, "out", 3);
  close(direction);

  int value = open("/sys/class/gpio/gpio91/value", O_WRONLY);
  if (value<0) {
    printf("fail to open value\n");
    return -1;
  }

  for (int i=0; i<1000000; i++)
  {
    write(value, "0", 1);
    write(value, "1", 1);
  }
  close(value);
  return 0;
}

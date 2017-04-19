#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <fcntl.h>
#include <err.h>
#include <unistd.h>


int pin_read(const char *pin)
{
  int exporter = open("/sys/class/gpio/export", O_WRONLY);
  if (exporter<0) {
    printf("fail to open exporter\n");
    return -1;
  }
  write(exporter, pin, 2);
  close(exporter);

  char dirstr[80];
  strcpy(dirstr, "/sys/class/gpio");
  strcat(dirstr, pin);
  strcat(dirstr, "/direction");

  int direction = open(dirstr, O_WRONLY);
  if (direction<0) {
    printf("fail to open direction\n");
    return -1;
  }
  write(direction, "in", 2);
  close(direction);

  char valstr[80];
  strcpy(valstr, "/sys/class/gpio");
  strcat(valstr, pin);
  strcat(valstr, "/value");

  int value = open(valstr, O_RDONLY);
  if (value<0) {
    printf("fail to open value\n");
    return -1;
  }
}

int main()
{
  
  int gp3_27 = pin_read("91");
  int gp6_31 = pin_read("191");
  int gp1_24 = pin_read("24");
  int gp7_8 = pin_read("200");

  int count=0;
  char c1,c2,c3,c4;
  for (int i=0;i<10000000/2;++i)
  {
    read(gp3_27, &c1, 1);
    read(gp6_31, &c2, 1);
    read(gp1_24, &c3, 1);
    read(gp7_8, &c4, 1);
    if (c1=='1')
    {
      ++count;
    }
    if (c2=='1')
    {
      ++count;
    }
    if (c3=='1')
    {
      ++count;
    }
    if (c4=='1')
    {
      ++count;
    }
    lseek(gp3_27, 0, SEEK_SET);
    lseek(gp6_31, 0, SEEK_SET);
    lseek(gp1_24, 0, SEEK_SET);
    lseek(gp7_8, 0, SEEK_SET);
  }
  printf("count is %d\n", count);
  close(gp3_27);
  close(gp6_31);
  close(gp1_24);
  close(gp7_8);
  return 0;
}

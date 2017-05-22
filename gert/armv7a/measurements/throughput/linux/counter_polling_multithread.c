#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <fcntl.h>
#include <err.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <time.h>
#include <stdlib.h>
#include <string.h>


int pin_read(const char *pin)
{
  int exporter = open("/sys/class/gpio/export", O_WRONLY);
  if (exporter<0) {
    printf("fail to open exporter\n");
    return -1;
  }
  write(exporter, pin, strlen(pin));
  close(exporter);

  char dirstr[80];
  strcpy(dirstr, "/sys/class/gpio/gpio");
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
  strcpy(valstr, "/sys/class/gpio/gpio");
  strcat(valstr, pin);
  strcat(valstr, "/value");

  int value = open(valstr, O_RDONLY);
  if (value<0) {
    printf("fail to open value\n");
    return -1;
  }
}

void poll(int pinfd)
{
  int i=0;
  int count=0;
  char c, oldc;
  c='0';
  oldc='0';
  time_t oldtime=time(NULL);
  time_t start_time = oldtime;
  while(difftime(time(NULL), start_time) <= 10.0)
  {
    time_t curtime=time(NULL);
    if (difftime(curtime, oldtime)>=1.0) {
      oldtime=curtime;
      printf("%d count -> %d\n", pinfd, count);
      count=0;
    }

    read(pinfd, &c, 1);
    if (c=='1' && oldc=='0')
    {
      ++count;
    }
    oldc=c;
    lseek(pinfd, 0, SEEK_SET);
  }
  printf("count is %d\n", count);
}

int main(int argc, char **argv)
{
  int nprocs = atoi(argv[1]);

  int input_pins[4];
  input_pins[0] = pin_read("91");
  input_pins[1] = pin_read("191");
  input_pins[2] = pin_read("24");
  input_pins[3] = pin_read("200");

  int i=0;
  for (i=0; i<nprocs; ++i)
  {
    if (fork()==0)
    {
      poll(input_pins[i]);
      exit(0);
    }
  }

  int child=0;
  for (child=0; child<nprocs; ++child)
  {
    wait(NULL);
  }

  for (i=0; i<nprocs; ++i)
  {
    close(input_pins[i]);
  }
  return 0;
}

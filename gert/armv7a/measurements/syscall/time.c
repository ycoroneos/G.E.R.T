#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>
#include <time.h>


int comp(const void *e1, const void *e2) {
  unsigned a = *((unsigned *)e1);
  unsigned b = *((unsigned *)e2);
  if (a>b) {
    return 1;
  }
  if (b>a) {
    return -1;
  }
  return 0;
}

#define COUNT 100000


int main() {
  struct timespec oldtime, newtime;
  //struct timespec results[COUNT];
  unsigned results[COUNT];

  int i=0;
  volatile pid_t pid;
  for (i=0; i<COUNT; ++i) {
    clock_gettime(CLOCK_REALTIME, &oldtime);
    pid=getppid();
    clock_gettime(CLOCK_REALTIME, &newtime);
    results[i] = newtime.tv_nsec - oldtime.tv_nsec;
  }
  printf("qsort\n");
  qsort(results, COUNT, sizeof(unsigned), comp);
  printf("median time for getppid: %dns\n", results[COUNT/2]);
  return 0;
}

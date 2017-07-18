// Copyright 2017 Yanni Coroneos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef CONSOLE_H
#define CONSOLE_H
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdarg.h>
#include "pl011_uart.h"
#define KNRM  "\x1B[0m"
#define KRED  "\x1B[31m"
#define KGRN  "\x1B[32m"
#define KYEL  "\x1B[33m"
#define KBLU  "\x1B[34m"
#define KMAG  "\x1B[35m"
#define KCYN  "\x1B[36m"
#define KWHT  "\x1B[37m"
void consoleinit(void);
void csprintf(char*,char*, ...);
void cputs(char *);
void switchuart();
void panic(char*) __attribute__((noreturn));
void printfmt(void *putdat, const char *fmt, ...);
void vprintfmt(void *putdat, const char *fmt, va_list ap);
int vcprintf(const char *fmt, va_list ap);
int cprintf(const char *fmt, ...);
size_t strnlen(const char *s, size_t maxlen);
#endif

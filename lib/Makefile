CC = gcc

ifndef os
   os = linux
endif

ifndef arch
   arch = 386
endif

CCFLAGS = -fPIC -Wall -Wextra -D$(os) -finline-functions -O3 -fno-strict-aliasing -fvisibility=hidden

ifeq ($(os), darwin)
   CCFLAGS += -DDARWIN -x objective-c -fobjc-arc
endif

ifeq ($(arch), 386)
   CCFLAGS += -DIA32
endif

ifeq ($(arch), amd64)
   CCFLAGS += -DAMD64
endif

OBJS = systray_$(os)_$(arch).o
DIR = $(shell pwd)

all: ../libsystray_$(os)_$(arch).a

%_$(os)_$(arch).o: %_$(os).c
	$(CC) $(CCFLAGS) -o $@ $< -c

libsystray_$(os)_$(arch).a: $(OBJS)
	ar -rcs $@ $^

../libsystray_$(os)_$(arch).a: libsystray_$(os)_$(arch).a
	cp libsystray_$(os)_$(arch).a ../

clean:
	rm -f *.o *.a ../libsystray*.a

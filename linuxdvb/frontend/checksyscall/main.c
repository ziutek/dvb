#include <sys/ioctl.h>
#include <stdio.h>
#include <linux/dvb/frontend.h>

int main() {
	printf("%x %x\n",  FE_SET_PROPERTY, FE_GET_PROPERTY);
	return 0;
}
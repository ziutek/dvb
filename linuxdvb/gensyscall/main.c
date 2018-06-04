#include <sys/ioctl.h>
#include <stdio.h>
#include <linux/dvb/frontend.h>
#include <linux/dvb/dmx.h>

typedef struct tuple tuple;
struct tuple {
	const char *name;
	int val;
};

tuple table[] = {
	{"_FE_GET_INFO", FE_GET_INFO},
	{"_FE_GET_FRONTEND", FE_GET_FRONTEND},
	{"_FE_SET_FRONTEND", FE_SET_FRONTEND},
	{"_FE_GET_EVENT", FE_GET_EVENT},
	{"_FE_READ_STATUS", FE_READ_STATUS},
	{"_FE_READ_BER", FE_READ_BER},
	{"_FE_READ_SIGNAL_STRENGTH", FE_READ_SIGNAL_STRENGTH},
	{"_FE_READ_SNR", FE_READ_SNR},
	{"_FE_READ_UNCORRECTED_BLOCKS", FE_READ_UNCORRECTED_BLOCKS},
	{"_FE_SET_TONE", FE_SET_TONE},
	{"_FE_SET_VOLTAGE", FE_SET_VOLTAGE},
	{"_FE_ENABLE_HIGH_LNB_VOLTAGE", FE_ENABLE_HIGH_LNB_VOLTAGE},

	//{"_FE_SET_PROPERTY", FE_SET_PROPERTY},
	//{"_FE_GET_PROPERTY", FE_GET_PROPERTY},

	{"_DMX_START", DMX_START},
	{"_DMX_STOP", DMX_STOP},
	{"_DMX_SET_BUFFER_SIZE", DMX_SET_BUFFER_SIZE},
	{"_DMX_SET_FILTER", DMX_SET_FILTER},
	{"_DMX_SET_PES_FILTER", DMX_SET_PES_FILTER},
	{"_DMX_GET_PES_PIDS", DMX_GET_PES_PIDS},
	{"_DMX_GET_STC", DMX_GET_STC},
	{"_DMX_ADD_PID", DMX_ADD_PID},
	{"_DMX_REMOVE_PID", DMX_REMOVE_PID},
};

int
main() {
	unsigned int i;
	for (i = 0; i < sizeof(table) / sizeof(tuple); i++) {
		printf("%s = 0x%08x\n", table[i].name, table[i].val);
	}
	struct dtv_property p;
	printf(
		"sizeof(dtv_property) = %d, sizeof(dtv_property.u) = %d\n",
		sizeof(p), sizeof(p.u)
	);
	return 0;
}

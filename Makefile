all:
	make -C ./driver-location
	make -C ./gateway
	make -C ./zombie-driver

test:
	make -C ./driver-location test
	make -C ./gateway test
	make -C ./zombie-driver test

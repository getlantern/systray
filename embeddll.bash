#!/bin/bash

###############################################################################
#
# This script regenerates the source file that embeds systray.dll
#
###############################################################################

go-bindata -nomemcopy -nocompress -pkg systray -prefix dll -o systraydll_windows.go dll
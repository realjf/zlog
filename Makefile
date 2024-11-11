# ##############################################################################
# # File: Makefile                                                             #
# # Project: zlog                                                              #
# # Created Date: 2024/11/11 11:17:41                                          #
# # Author: realjf                                                             #
# # -----                                                                      #
# # Last Modified: 2024/11/11 11:18:27                                         #
# # Modified By: realjf                                                        #
# # -----                                                                      #
# #                                                                            #
# ##############################################################################




B ?= master
M ?= update

.PHONY: push
push:
	@git add -A && git commit -m "${M}" && git push origin ${B}

-include gomk/main.mk
-include local/Makefile

superclean: clean
ifeq ($(unameS),Windows)
ifneq ($(wildcard wwwotp),)
	@powershell -c Remove-Item -Force -Recurse ./wwwotp
endif
else
	@rm -fr wwwotp
endif

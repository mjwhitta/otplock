-include gomk/main.mk
-include local/Makefile

superclean: clean
ifeq ($(unameS),windows)
ifneq ($(wildcard wwwotp),)
	@remove-item -force -recurse ./wwwotp
endif
else
	@rm -f -r wwwotp
endif

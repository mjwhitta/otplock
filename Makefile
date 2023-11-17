-include gomk/main.mk
-include local/Makefile

ifneq ($(unameS),windows)
spellcheck:
	@codespell -f -L hilighter -S ".git,*.pem"
endif

superclean: clean
ifeq ($(unameS),windows)
ifneq ($(wildcard wwwotp),)
	@powershell -c Remove-Item -Force -Recurse ./wwwotp
endif
else
	@rm -f -r wwwotp
endif

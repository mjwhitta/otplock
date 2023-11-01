-include gomk/main.mk
-include local/Makefile

ifneq ($(unameS),Windows)
spellcheck:
	@codespell -f -L hilighter -S ".git,*.pem"
endif

superclean: clean
ifeq ($(unameS),Windows)
ifneq ($(wildcard wwwotp),)
	@powershell -c Remove-Item -Force -Recurse ./wwwotp
endif
else
	@rm -fr wwwotp
endif

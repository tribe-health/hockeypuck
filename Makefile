
PANDOC=pandoc

HTML_FILES=community.html \
configuration.html \
contributing.html \
install-source.html \
install-tarball.html \
install-ubuntu.html \
juju.html \
pre-populating.html \
README.html \
running.html

all: $(HTML_FILES)

%.html: %.md
	$(PANDOC) -f markdown_github+backtick_code_blocks -t html -o $@ $<
	sed -i 's/\.md\"/.html"/g' $@

clean:
	$(RM) $(HTML_FILES)

.PHONY: all clean

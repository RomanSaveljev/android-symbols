all : landing.html incremental-storing.html symbols-listing.html retrieving.html usage-printing.html figuring-build.html sourcing-aliases.html
	@echo Done converting markdown to HTML

%.html : %.mkd
	markdown $< >$@

ALL_FILES := $(shell find ./* -type f)

all : $(addsuffix .aaaa,$(ALL_FILES))
	@echo $(ALL_FILES)

define PROCESSING_RULE
${1}.aaaa : ${1}
	split -a 4 -b 4M ${1} ${1}.
	rm -f ${1}
endef

$(foreach FILE,$(ALL_FILES),\
  $(eval $(call PROCESSING_RULE, $(FILE))))

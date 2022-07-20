.PHONY: all-in-one 
all-in-one: beatify.gv
	dot -Tjpeg -Gdpi=72 beatify.gv -o factorio.jpeg

.INTERMEDIATE: beatify.gv 
beatify.gv:
	gvpr -c 'N[outdegree==0]{color="red",penwidth=4}' -o $@ factorio.gv

.PHONY: groups
groups: tmp_factorio.gv
	dot -Tjpeg -Gdpi=72 tmp_factorio.gv -o factorio_groups.jpeg

.INTERMEDIATE: tmp_factorio.gv
tmp_factorio.gv:
	go run main.go > $@


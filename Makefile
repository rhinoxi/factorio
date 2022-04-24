.PHONY: jpeg
jpeg: beatify.gv
	dot -Tjpeg -Gdpi=72 beatify.gv -o factorio.jpeg

.INTERMEDIATE: beatify.gv 
beatify.gv:
	gvpr -c 'N[outdegree==0]{color="red",penwidth=4}' -o $@ factorio.gv


package main

import (
	"fmt"
	"io/ioutil"

	"github.com/awalterschulze/gographviz"
)

func convertAttrs(attrs gographviz.Attrs) map[string]string {
	out := make(map[string]string)
	for k, v := range attrs {
		out[string(k)] = v
	}
	return out
}

func groupNodeName(groupName, nodeName string) string {
	return groupName + "_" + nodeName
}

func _group(graph *gographviz.Graph, groupName string, leafNodeName string) ([]string, error) {
	var groupNodes []string
	n, ok := graph.Nodes.Lookup[leafNodeName]
	if !ok {
		return nil, nil
	}
	groupNodes = append(groupNodes, leafNodeName)

	graph.AddNode(groupName, groupNodeName(groupName, leafNodeName), convertAttrs(n.Attrs))

	for srcNodeName, edges := range graph.Edges.DstToSrcs[leafNodeName] {
		nodes, err := _group(graph, groupName, srcNodeName)
		if err != nil {
			return nil, err
		}

		for _, edge := range edges {
			graph.AddEdge(
				groupNodeName(groupName, srcNodeName),
				groupNodeName(groupName, leafNodeName),
				edge.Dir,
				convertAttrs(edge.Attrs),
			)
		}

		groupNodes = append(groupNodes, nodes...)
	}

	return groupNodes, nil
}

func group(graph *gographviz.Graph, groupName string, leafNodeNames []string) ([]string, error) {
	graph.AddSubGraph("", groupName, map[string]string{"color": "red"})
	var groupNodes []string
	for _, leafNodeName := range leafNodeNames {
		nodes, err := _group(graph, groupName, leafNodeName)
		if err != nil {
			return nil, err
		}
		groupNodes = append(groupNodes, nodes...)
	}

	return groupNodes, nil
}

func unique(l []string) []string {
	m := make(map[string]struct{})
	for _, s := range l {
		m[s] = struct{}{}
	}
	var out []string
	for s := range m {
		out = append(out, s)
	}

	return out
}

var groups = map[string][]string{
	"red":      {"AutomationSciencePack_RED"},
	"green":    {"LogisticSciencePack_GREEN"},
	"blue":     {"ChemicalSciencePack_BLUE"},
	"yellow":   {"UtilitySciencePack_YELLOW"},
	"gray":     {"MilitarySciencePack_GRAY"},
	"purple":   {"ProductionSciencePack_PURPLE"},
	"inserter": {"LongHandedInserter", "StackFilterInserter", "FilterInserter"},
	"belt":     {"ExpressTransportBelt", "ExpressUndergroundBelt", "ExpressSplitter"},
	"rocket":   {"Satellite", "RocketSilo", "RocketPart"},
}

func main() {
	path := "factorio.gv"

	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	graphAst, err := gographviz.Parse(b)
	if err != nil {
		panic(err)
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}

	var nodes []string
	for groupName, leafNodeNames := range groups {
		ns, err := group(graph, groupName, leafNodeNames)
		if err != nil {
			panic(err)
		}
		nodes = append(nodes, ns...)
	}

	for _, node := range nodes {
		graph.RemoveNode("", node)
	}

	graph.SetStrict(true)

	fmt.Println(graph.String())
}

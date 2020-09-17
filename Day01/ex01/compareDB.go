package main

import (
	"fmt"
	"os"
	"sort"
)

func usage() {
	fmt.Println("Usage: compareDB --old filename.xml | filename.json --new filename.xml | filename.json")
}

func main() {
	if len(os.Args) != 5 || os.Args[1] != "--old" || os.Args[3] != "--new" {
		usage()
		return
	}
	recipesOld, err := readUniversal(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		return
	}
	recipesNew, err := readUniversal(os.Args[4])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		return
	}
	sort.Slice(recipesOld.Cake, func(i, j int) bool { return recipesOld.Cake[i].Name < recipesNew.Cake[j].Name })
	sort.Slice(recipesNew.Cake, func(i, j int) bool { return recipesNew.Cake[i].Name < recipesNew.Cake[j].Name })
	for i, j, n, m := 0, 0, len(recipesOld.Cake), len(recipesNew.Cake); i < n || j < m; {
		if i == n || recipesOld.Cake[i].Name > recipesNew.Cake[j].Name {
			fmt.Println("ADDED cake \"" + recipesNew.Cake[j].Name + "\"")
			j++
		} else if j == m || recipesOld.Cake[i].Name < recipesNew.Cake[j].Name {
			fmt.Println("REMOVED cake \"" + recipesOld.Cake[i].Name + "\"")
			i++
		} else {
			if recipesOld.Cake[i].Stovetime != recipesNew.Cake[j].Stovetime {
				fmt.Println("CHANGED cooking time for cake \"" + recipesOld.Cake[i].Name + "\" - \"" + recipesNew.Cake[j].Stovetime + "\" instead of \"" + recipesOld.Cake[i].Stovetime + "\"")
			}
			sort.Slice(recipesOld.Cake[i].Item, func(a, b int) bool { return recipesOld.Cake[i].Item[a].Itemname < recipesOld.Cake[i].Item[b].Itemname })
			sort.Slice(recipesNew.Cake[j].Item, func(a, b int) bool { return recipesNew.Cake[j].Item[a].Itemname < recipesNew.Cake[j].Item[b].Itemname })
			for a, b, la, lb := 0, 0, len(recipesOld.Cake[i].Item), len(recipesNew.Cake[j].Item); a < la || b < lb; {
				if a == la || recipesOld.Cake[i].Item[a].Itemname > recipesNew.Cake[j].Item[b].Itemname {
					fmt.Println("ADDED ingredient \"" + recipesNew.Cake[j].Item[b].Itemname + "\" for cake \"" + recipesNew.Cake[j].Name + "\"")
					b++
				} else if b == lb || recipesOld.Cake[i].Item[a].Itemname < recipesNew.Cake[j].Item[b].Itemname {
					fmt.Println("REMOVED ingredient \"" + recipesOld.Cake[i].Item[a].Itemname + "\" for cake \"" + recipesOld.Cake[i].Name + "\"")
					a++
				} else {
					if recipesOld.Cake[i].Item[a].Itemunit == "" && recipesNew.Cake[j].Item[b].Itemunit != "" {
						fmt.Println("ADDED unit \"" + recipesNew.Cake[j].Item[b].Itemunit + "\" for ingredient \"" + recipesNew.Cake[j].Item[b].Itemname + "\" for cake \"" + recipesNew.Cake[j].Name + "\"")
					} else if recipesOld.Cake[i].Item[a].Itemunit != "" && recipesNew.Cake[j].Item[b].Itemunit == "" {
						fmt.Println("REMOVED unit \"" + recipesOld.Cake[i].Item[a].Itemunit + "\" for ingredient \"" + recipesOld.Cake[i].Item[a].Itemname + "\" for cake \"" + recipesOld.Cake[i].Name + "\"")
					} else if recipesOld.Cake[i].Item[a].Itemunit != "" && recipesNew.Cake[j].Item[b].Itemunit != "" && recipesOld.Cake[i].Item[a].Itemunit != recipesNew.Cake[j].Item[b].Itemunit {
						fmt.Println("CHANGED unit for ingredient \"" + recipesOld.Cake[i].Item[a].Itemname + "\" for cake \"" + recipesOld.Cake[i].Name + "\" - \"" + recipesNew.Cake[j].Item[b].Itemunit + "\" instead of \"" + recipesOld.Cake[i].Item[a].Itemunit + "\"")
					}
					if recipesOld.Cake[i].Item[a].Itemcount != recipesNew.Cake[j].Item[b].Itemcount {
						fmt.Println("CHANGED unit count for ingredient \"" + recipesOld.Cake[i].Item[a].Itemname + "\" for cake \"" + recipesOld.Cake[i].Name + "\" - \"" + recipesNew.Cake[j].Item[b].Itemcount + "\" instead of \"" + recipesOld.Cake[i].Item[a].Itemcount + "\"")
					}
					a++
					b++
				}
			}
			i++
			j++
		}
	}
}

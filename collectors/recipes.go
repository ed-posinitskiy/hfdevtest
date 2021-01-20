package collectors

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type RecipesCollector struct {
	recipesCounter map[string]int
	recipesNames   []string
	searchPattern  *regexp.Regexp
}

func (c *RecipesCollector) Process(r *Record) {
	if _, ok := c.recipesCounter[r.Recipe]; !ok {
		c.recipesCounter[r.Recipe] = 0
		c.recipesNames = append(c.recipesNames, r.Recipe)
	}
	c.recipesCounter[r.Recipe] += 1
}

func (c *RecipesCollector) Report() map[string]interface{} {
	sort.Strings(c.recipesNames)

	return map[string]interface{}{
		"unique_recipe_count": len(c.recipesCounter),
		"count_per_recipe":    c.sortedCounters(),
		"match_by_name":       c.searchRecipes(),
	}
}

func (c *RecipesCollector) sortedCounters() []map[string]interface{} {
	var counters []map[string]interface{}

	for _, recipe := range c.recipesNames {
		entry := map[string]interface{}{
			"recipe": recipe,
			"count":  c.recipesCounter[recipe],
		}

		counters = append(counters, entry)
	}

	return counters
}

func (c *RecipesCollector) searchRecipes() []string {
	var found []string

	if c.searchPattern == nil {
		return found
	}

	for _, recipe := range c.recipesNames {
		if match := c.searchPattern.MatchString(recipe); match {
			found = append(found, recipe)
		}
	}

	return found
}

func NewRecipesCollector(searchRecipes []string) *RecipesCollector {
	var searchPattern *regexp.Regexp

	if len(searchRecipes) > 0 {
		searchPattern = regexp.MustCompile(fmt.Sprintf("(%s)", strings.Join(searchRecipes, "|")))
	}

	return &RecipesCollector{
		recipesCounter: make(map[string]int),
		searchPattern:  searchPattern,
		recipesNames:   []string{},
	}
}

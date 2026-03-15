package random

import (
	"fmt"
	"sort"
	"strings"

	"github.com/m-oons/lexi/internal/random"
	"github.com/spf13/cobra"
)

type options struct {
	count     int
	words     string
	separator string
	minLength int
	maxLength int
	regex     string
	format    string
	order     string
	seed      int64
}

func NewCmd() *cobra.Command {
	var opts options

	cmd := &cobra.Command{
		Use:   "random",
		Short: "Generate random words from a list",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get word list
			var words []string
			if opts.words == "" {
				words = random.LoadDefaultWords(opts.separator)
			} else {
				loadedWords, err := random.LoadWords(opts.words, opts.separator)
				if err != nil {
					return err
				}
				words = loadedWords
			}

			// Filter words
			filteredWords, err := random.FilterWords(words, opts.minLength, opts.maxLength, opts.regex)
			if err != nil {
				return err
			}

			// Generate words
			count := max(opts.count, 1)
			results := make([]string, 0, count)
			rng := random.NewGenerator(opts.seed)
			placeholders := strings.Count(opts.format, "{}")
			placeholdersLower := strings.Count(opts.format, "{l}")
			placeholdersUpper := strings.Count(opts.format, "{u}")

			for range count {
				result := opts.format
				for range placeholders {
					word := filteredWords[rng.Intn(len(filteredWords))]
					result = strings.Replace(result, "{}", word, 1)
				}
				for range placeholdersLower {
					word := filteredWords[rng.Intn(len(filteredWords))]
					result = strings.Replace(result, "{l}", strings.ToLower(word), 1)
				}
				for range placeholdersUpper {
					word := filteredWords[rng.Intn(len(filteredWords))]
					result = strings.Replace(result, "{u}", strings.ToUpper(word), 1)
				}
				results = append(results, result)
			}

			// Order results
			if opts.order != "" {
				switch strings.ToLower(opts.order) {
				case "asc", "ascending":
					sort.Strings(results)
				case "desc", "descending":
					sort.Sort(sort.Reverse(sort.StringSlice(results)))
				default:
					return fmt.Errorf("invalid order %q: must be one of [asc, ascending, desc, descending]", opts.order)
				}
			}

			// Output results
			for _, word := range results {
				fmt.Fprintln(cmd.OutOrStdout(), word)
			}

			return nil
		},
	}

	cmd.PersistentFlags().IntVarP(&opts.count, "count", "n", 1, "maximum number of words to generate")
	cmd.PersistentFlags().StringVarP(&opts.words, "words", "w", "", "words to use for randomization")
	cmd.PersistentFlags().StringVarP(&opts.separator, "separator", "s", "", "separator to use between words (defaults to newline for files, comma for inline words)")
	cmd.PersistentFlags().IntVar(&opts.minLength, "min", 0, "minimum length of words to generate")
	cmd.PersistentFlags().IntVar(&opts.maxLength, "max", 0, "maximum length of words to generate")
	cmd.PersistentFlags().StringVarP(&opts.regex, "regex", "r", "", "regular expression to filter words")
	cmd.PersistentFlags().StringVarP(&opts.format, "format", "f", "{}", "output format for each line; use {}/{l}/{u} placeholders for random words")
	cmd.PersistentFlags().StringVarP(&opts.order, "order", "o", "", "order of output words")
	cmd.PersistentFlags().Int64Var(&opts.seed, "seed", -1, "seed for random number generator")

	return cmd
}

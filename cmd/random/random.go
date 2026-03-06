package random

import (
	"fmt"

	"github.com/m-oons/lexi/internal/random"
	"github.com/spf13/cobra"
)

type options struct {
	count     int
	words     string
	separator string
	minLength int
	maxLength int
	prefix    string
	suffix    string
	regex     string
	uppercase bool
	unique    bool
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
			words = filteredWords

			// Generate words
			rng := random.NewGenerator(opts.seed)

			if opts.unique {
				remaining := words
				seenCap := max(opts.count, 0)
				seen := make(map[string]struct{}, seenCap)

				for len(seen) < opts.count && len(remaining) > 0 {
					index := rng.Intn(len(remaining))
					word := remaining[index]

					last := len(remaining) - 1
					remaining[index] = remaining[last]
					remaining = remaining[:last]

					transformed := random.TransformWord(word, opts.prefix, opts.suffix, opts.uppercase)
					if _, ok := seen[transformed]; ok {
						continue
					}

					seen[transformed] = struct{}{}
					fmt.Println(transformed)
				}

				return nil
			}

			for i := 0; i < opts.count; i++ {
				if len(words) == 0 {
					return nil
				}

				index := rng.Intn(len(words))
				word := words[index]
				transformed := random.TransformWord(word, opts.prefix, opts.suffix, opts.uppercase)
				fmt.Println(transformed)
			}

			return nil
		},
	}

	cmd.PersistentFlags().IntVarP(&opts.count, "count", "n", 1, "maximum number of words to generate")
	cmd.PersistentFlags().StringVarP(&opts.words, "words", "w", "", "words to use for randomization")
	cmd.PersistentFlags().StringVar(&opts.separator, "separator", "", "separator to use between words (defaults to newline for files, comma for inline words)")
	cmd.PersistentFlags().IntVar(&opts.minLength, "min", 0, "minimum length of words to generate")
	cmd.PersistentFlags().IntVar(&opts.maxLength, "max", 0, "maximum length of words to generate")
	cmd.PersistentFlags().StringVarP(&opts.prefix, "prefix", "p", "", "prefix to add to each word")
	cmd.PersistentFlags().StringVarP(&opts.suffix, "suffix", "s", "", "suffix to add to each word")
	cmd.PersistentFlags().StringVarP(&opts.regex, "regex", "r", "", "regular expression to filter words")
	cmd.PersistentFlags().BoolVar(&opts.uppercase, "uppercase", false, "convert words to uppercase")
	cmd.PersistentFlags().BoolVarP(&opts.unique, "unique", "u", false, "don't generate duplicate words")
	cmd.PersistentFlags().Int64Var(&opts.seed, "seed", -1, "seed for random number generator")

	return cmd
}

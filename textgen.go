package textgen

import (
    "bufio"
    "io"
    "math/rand"
    "strings"
    "time"
)

func init() {
    rand.Seed(time.Now().UTC().UnixNano())
}

/*
Prefix is used to create keys in the markov chain map
*/
type Prefix []string

/*
PushBack appends word to the prefix pushing all
current back one index which removes the first word
*/
func (p Prefix) PushBack(word string) string {
    rval := p[0]
    size := len(p)
    for i := 1; i < size; i += 1 {
        p[i-1] = p[i]
    }
    p[size-1] = word
    return rval;
}

/*
NewPrefix creates a new prefix with the given size
*/
func NewPrefix(size uint) Prefix {
    return make(Prefix, size)
}

/*
String formats the prefix for use as a map key
*/
func (p Prefix) String() string {
    return strings.Join(p, " ")
}

/*
Generator is the main object for generating random text
using markov chains
*/
type Generator struct {
    chains map[string][]string
    chainlen uint
}

/*
NewGenerator creates a new text generator with the
specified markov chain length
*/
func NewGenerator(chainlen uint) *Generator {
    return &Generator{make(map[string][]string), chainlen}
}

/*
Load reads text to create markov chains which are used
to generate new random text
*/
func (g *Generator) Load(r *bufio.Reader) error {
    prefix := NewPrefix(g.chainlen)
    chunk, err := r.ReadString(' ')
    for err == nil {
        for _, word := range tokenize(chunk, "\n") {
            if word != "" {
                /* 
                this makes a newly craated chain match
                the start of new paragraphs
                */
                if word == "\n" {
                    word = ""
                }
                key := prefix.String()
                words := g.chains[key]
                words = append(words, word)
                g.chains[key] = words
                prefix.PushBack(word)
            }
        }
        chunk, err = r.ReadString(' ')
    }
    if err == io.EOF {
        err = nil
    }
    return err
}

/*
tokenize splits the chunk into parts using the token separator.
any tokens in the string will be their own item in the returned slice
*/
func tokenize(chunk string, token string) []string {
    var rval []string
    index := strings.Index(chunk, token)
    for index > -1 {
        part := strings.TrimSpace(chunk[:index])
        rval = append(rval, part)
        rval = append(rval, token)
        chunk = chunk[index+1:]
        index = strings.Index(chunk, token)
    }
    return append(rval, strings.TrimSpace(chunk))
}

/*
Generate writes new random text which will not be longer than maxlen.
the generated text could be shorter if the markov chain hits a dead end
*/
func (g *Generator) Generate(w *bufio.Writer, maxlen uint) error {
    var err error
    prefix := NewPrefix(g.chainlen)
    var count uint
    for ; err == nil && count < maxlen; count += 1 {
        key := prefix.String()
        words, exists := g.chains[key]
        if !exists {
            break
        }
        randindex := rand.Intn(len(words))
        nextword := words[randindex]
        _, err = w.WriteString(nextword + " ")
        prefix.PushBack(nextword)
    }
    /* TODO try to end with a period? */
    if err == nil {
        err = w.Flush()
    }
    return err
}


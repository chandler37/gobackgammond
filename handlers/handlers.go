// HTTP handlers. NB: You must call math/rand.Seed() first.
//
// TODO(chandler37): Support playing just one side of the game, and support
// choosing which AI to play. Support match play. Support doubling.
package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"

	"github.com/chandler37/gobackgammon/ai"
	"github.com/chandler37/gobackgammon/brd"
	"github.com/chandler37/gobackgammon/json"
	"github.com/chandler37/gobackgammon/svg"

	mysvg "github.com/chandler37/gobackgammond/svg"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
<html>

<head>
  <title>Backgammon</title>
</head>
<body>

<a href="/">This</a> is <a href="https://github.com/chandler37/gobackgammond">github.com/chandler37/gobackgammond</a>.<br>

<a href="/game">Click here to start a game.</a><br>

<a href="/game?s=qlYqKlWyUjIxIhBShoawsEK4JNwEYa2hiVItIAAA__8&t=">Click here for a game that is all but over (so you can see what sweet, sweet victory looks like).</a><br>

</body></html>`))
}

func SvgHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := token(r)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, fmt.Sprintf("error getting token: %s", err))
		return
	}
	decompressedToken, err := json.Decompress(tok)
	if err != nil {
		decompressedToken = tok
	}
	b, err := json.Deserialize(decompressedToken)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, fmt.Sprintf("error deserializing game state: %s", err))
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	svg.Board(240, b, mysvg.New(w))
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := token(r)
	var b *brd.Board
	if err == noGameFoundError {
		tok = ""
		b = brd.New(false)
	} else if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, fmt.Sprintf("error getting token: %s", err))
		return
	} else {
		decompressedToken, err := json.Decompress(tok)
		if err != nil {
			decompressedToken = tok
		}
		b, err = json.Deserialize(decompressedToken)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprint(w, fmt.Sprintf("error deserializing game state: %s", err))
			return
		}
	}
	if shouldTakeTurn(r) {
		victor, stakes, score := b.TakeTurn(nil, nil)
		if victor != brd.NoChecker {
			victorString := "Red"
			if victor != brd.Red {
				victorString = "White"
			}
			stakesString := fmt.Sprintf("%d points", stakes)
			if stakes == 1 {
				stakesString = "1 point"
			}
			scoreString := fmt.Sprintf("White:%d, Red:%d", score.WhiteScore, score.RedScore)
			if score.Goal > 0 {
				scoreString += fmt.Sprintf(" playing to %d", score.Goal)
			}
			vs := victoryState{
				Serialization: tok,
				Stakes:        stakesString,
				Victor:        victorString,
				Score:         scoreString,
			}
			var buf bytes.Buffer // so we don't render half a template
			err := victoryTemplate.Execute(&buf, vs)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, fmt.Sprintf("error executing template: %v", err))
				return
			}
			io.Copy(w, &buf)
			return
		}
	}
	boards, err := smartBoards(b.LegalContinuations())
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, fmt.Sprintf("error making continuations: %s", err))
		return
	}
	bs, err := json.Serialize(b)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, fmt.Sprintf("error serializing game state: %s", err))
		return
	}
	theState := state{
		Token:              tok,
		Board:              smartBoard{b, json.Compress(bs), ""},
		LegalContinuations: boards}
	var buf bytes.Buffer // so we don't render half a template
	err = gameTemplate.Execute(&buf, theState)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, fmt.Sprintf("error executing template: %v", err))
		return
	}
	io.Copy(w, &buf)
}

type state struct {
	Token              string
	Board              smartBoard
	LegalContinuations []smartBoard
}

type victoryState struct {
	Stakes        string
	Serialization string
	Victor        string
	Score         string
}

var victoryTemplate = newVictoryTemplate()
var gameTemplate = newGameTemplate()

func newVictoryTemplate() *template.Template {
	return template.Must(
		template.New("Victory").Parse(
			`<!DOCTYPE html>
<html>

<head>
  <title>Backgammon Game Results</title>
</head>
<body>

<a href="/">This</a> is <a href="https://github.com/chandler37/gobackgammond">github.com/chandler37/gobackgammond</a>.<br>

<a href="/game">Click here for a new game.</a><br>

The final board is <img src='/game.svg?s={{.Serialization}}'><br>

The final score is {{.Score}}.<br>

<br>

Congratulations on winning {{.Stakes}}, {{.Victor}}!<br>

</body></html>`))
}

func newGameTemplate() *template.Template {
	return template.Must(
		template.New("Game").Parse(
			`<!DOCTYPE html>
<html>

<head>
  <title>Backgammon Game</title>
</head>
<body>
<a href="/">This</a> is <a href="https://github.com/chandler37/gobackgammond">github.com/chandler37/gobackgammond</a>.<br>

<a href="/game">Click here for a new game.</a><br>

The current board is <img src='/game.svg?s={{.Board.Serialization}}'><br>

Which of the following is your play?<br>
<ul>
    {{range .LegalContinuations}}
        <li>{{.Hint}}<br><a href='/game?s={{.Serialization}}&t='><img src='/game.svg?s={{.Serialization}}'></a></li>
    {{end}}
</ul>
</body></html>`))
}

type smartBoard struct {
	Board         *brd.Board
	Serialization string
	Hint          string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func smartBoards(bb []*brd.Board) ([]smartBoard, error) {
	aiChoices := ai.MakePlayerConservative(0, nil)(bb)
	r := make([]smartBoard, 0, len(bb))
	for i, ab := range aiChoices {
		ser, err := json.Serialize(ab.Board)
		if err != nil {
			return nil, fmt.Errorf("serialization error: %v", err)
		}
		hint := ""
		sum := ""
		if ab.Analysis != nil {
			sum = " because it " + ab.Analysis.Summary()
		}
		if i == 0 {
			hint = fmt.Sprintf("PlayerConservative's choice%s: ", sum)
		} else if ab.Analysis != nil {
			hint = fmt.Sprintf("PlayerConservative's analysis: %s", ab.Analysis.Summary())
		}
		r = append(r, smartBoard{ab.Board, json.Compress(ser), hint})
	}
	return r, nil
}

var noGameFoundError error = fmt.Errorf("no game found")

func token(r *http.Request) (string, error) {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", fmt.Errorf("bad URL: %v", err) // TODO(chandler37): Use pkg.errors.Wrap
	}
	s := m["s"]
	if len(s) == 0 {
		return "", noGameFoundError
	}
	if len(s) > 1 {
		return "", fmt.Errorf("too many games found")
	}
	if len(s[0]) == 0 {
		return "", fmt.Errorf("empty game found")
	}
	return s[0], nil
}

func shouldTakeTurn(r *http.Request) bool {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return false
	}
	s := m["t"]
	if len(s) == 0 {
		return false
	}
	return true
}

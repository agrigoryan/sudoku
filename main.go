package main

import (
	"fmt"
	"image/jpeg"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/agrigoryan/gocsp/csp"
	"github.com/agrigoryan/sudoku/sudoku"
)

type solverListener struct {
	boards     []sudoku.Board
	varIndices []int
}

func (sl *solverListener) addBoard(a *csp.Assignment, varIdx int) {
	resultBoard := sudoku.NewBoard()
	for i := 0; i < a.NumVariables(); i++ {
		v, ok := a.AssignedValue(i)
		if !ok {
			resultBoard[i] = 0
			continue
		}
		resultBoard[i] = byte(v)
	}
	sl.boards = append(sl.boards, resultBoard)
	sl.varIndices = append(sl.varIndices, varIdx)
}

func (sl *solverListener) ValueAssigned(a *csp.Assignment, varIdx int) {
	sl.addBoard(a, varIdx)
}

func (sl *solverListener) ValueUnassigned(a *csp.Assignment, varIdx int) {
	// sl.addBoard(a, varIdx)
}

func main() {
	server := http.NewServeMux()
	server.HandleFunc("/sudoku", func(writer http.ResponseWriter, request *http.Request) {
		board := sudoku.NewBoardFromString(request.FormValue("board"))
		if !board.IsValid() {
			writer.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(writer, "invalid board")
			return
		}

		problem := sudoku.New(board)
		solver := csp.NewBacktrackingSolver(csp.MRVVariableSelector, csp.FirstDomainValueSelector, sudoku.InferenceFunc)
		listener := &solverListener{}
		solver.Listener = listener
		result := solver.Solve(problem)
		if result == nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(writer, "can not solve the puzzle")
			return
		}

		mw := multipart.NewWriter(writer)
		defer mw.Close()

		writer.Header().Set("Content-Type", fmt.Sprintf(
			"multipart/x-mixed-replace;boundary=%s",
			mw.Boundary(),
		))
		writer.WriteHeader(200)

		ticker := time.NewTicker(time.Second / 30)
		defer ticker.Stop()
		frame := 0

		for {
			select {
			case <-ticker.C:
				w, _ := mw.CreatePart(textproto.MIMEHeader{
					"Content-Type": []string{"image/jpeg"},
				})
				jpeg.Encode(w, sudoku.CreateBoardImage(listener.boards[frame], listener.varIndices[frame]), &jpeg.Options{Quality: jpeg.DefaultQuality})
				frame++
				if frame == len(listener.boards) {
					return
				}
			case <-request.Context().Done():
				return
			}
		}
	})

	err := http.ListenAndServe("0.0.0.0:4040", server)
	if err != nil {
		slog.Error("error starting http server", "error", err)
	}
}

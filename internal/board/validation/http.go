// Copyright 2023 The chessBox Crew
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation

import (
	"net/http"

	"github.com/sadmadrus/chessBox/internal/board"
)

// Validator — http.HandlerFunc для валидации позиции на доске.
func Validator(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" && r.Method != "HEAD" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	b, err := board.FromUsFEN(r.URL.Query().Get("board"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if IsLegal(*b) {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusForbidden)
}
